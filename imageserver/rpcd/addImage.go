package rpcd

import (
	"errors"
	"fmt"
	"github.com/Symantec/Dominator/lib/hash"
	"github.com/Symantec/Dominator/proto/imageserver"
)

func (t *rpcType) AddImage(request imageserver.AddImageRequest,
	reply *imageserver.AddImageResponse) error {
	if imageDataBase.CheckImage(request.ImageName) {
		return errors.New("image already exists")
	}
	if request.Image == nil {
		return errors.New("nil image")
	}
	if request.Image.FileSystem == nil {
		return errors.New("nil file-system")
	}
	// Verify all objects are available.
	hashes := make([]hash.Hash, 0,
		len(request.Image.FileSystem.RegularInodeTable))
	for _, inode := range request.Image.FileSystem.RegularInodeTable {
		if inode.Size > 0 {
			hashes = append(hashes, inode.Hash)
		}
	}
	objectSizes, err := imageDataBase.ObjectServer().CheckObjects(hashes)
	if err != nil {
		return err
	}
	for index, size := range objectSizes {
		if size < 1 {
			return errors.New(fmt.Sprintf("object: %x is not available",
				hashes[index]))
		}
	}
	request.Image.FileSystem.RebuildInodePointers()
	logger.Printf("AddImage(%s)\n", request.ImageName)
	return imageDataBase.AddImage(request.Image, request.ImageName)
}
