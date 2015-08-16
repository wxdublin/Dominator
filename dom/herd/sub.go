package herd

import (
	"fmt"
	"github.com/Symantec/Dominator/lib/constants"
	"github.com/Symantec/Dominator/lib/hash"
	subproto "github.com/Symantec/Dominator/proto/sub"
	"net/rpc"
	"strings"
)

func (herd *Herd) pollNextSub() bool {
	if herd.nextSubToPoll >= uint(len(herd.subsByIndex)) {
		herd.nextSubToPoll = 0
		return true
	}
	sub := herd.subsByIndex[herd.nextSubToPoll]
	herd.nextSubToPoll++
	sub.poll(herd)
	return false
}

func (sub *Sub) poll(herd *Herd) {
	if sub.connection == nil {
		hostname := strings.SplitN(sub.hostname, "*", 2)[0]
		var err error
		sub.connection, err = rpc.DialHTTP("tcp",
			fmt.Sprintf("%s:%d", hostname, constants.SubPortNumber))
		if err != nil {
			fmt.Printf("Error dialing\t%s\n", err)
			return
		}
	}
	var request subproto.PollRequest
	request.HaveGeneration = sub.generationCount
	var reply subproto.PollResponse
	err := sub.connection.Call("Subd.Poll", request, &reply)
	if err != nil {
		fmt.Printf("Error calling\t%s\n", err)
		return
	}
	fs := reply.FileSystem
	if fs != nil {
		fs.RebuildPointers()
		sub.fileSystem = fs
		sub.generationCount = reply.GenerationCount
		fmt.Printf("Polled: %s, GenerationCount=%d\n",
			sub.hostname, reply.GenerationCount)
		sub.fetchMissingObjects(herd, sub.requiredImage)
	}
	// TODO(rgooch): Compare required image with polled image and send RPCs
	//               to fetch files and later update sub.
}

// Returns true if all required objects are available.
func (sub *Sub) fetchMissingObjects(herd *Herd, imageName string) bool {
	if imageName == "" {
		return false
	}
	image := herd.getImage(imageName)
	if image == nil {
		return false
	}
	missingObjects := make(map[hash.Hash]bool)
	for _, inode := range image.FileSystem.RegularInodeTable {
		if inode.Size > 0 {
			missingObjects[inode.Hash] = true
		}
	}
	for _, hash := range sub.fileSystem.ObjectCache {
		delete(missingObjects, hash)
	}
	for _, inode := range sub.fileSystem.RegularInodeTable {
		if inode.Size > 0 {
			delete(missingObjects, inode.Hash)
		}
	}
	if len(missingObjects) < 1 {
		return true
	}
	// TODO(rgooch): Remove debugging output.
	fmt.Printf("Objects needing to be fetched: %d\n", len(missingObjects))
	var request subproto.FetchRequest
	var reply subproto.FetchResponse
	request.ServerAddress = herd.ImageServerAddress
	for hash, _ := range missingObjects {
		request.Objects = append(request.Objects, hash)
	}
	err := sub.connection.Call("Subd.Fetch", request, &reply)
	if err != nil {
		fmt.Printf("Error calling\t%s\n", err)
		return false
	}
	return false
}
