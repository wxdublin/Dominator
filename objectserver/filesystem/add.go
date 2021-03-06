package filesystem

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/Symantec/Dominator/lib/hash"
	"github.com/Symantec/Dominator/lib/objectcache"
	"os"
	"path"
)

const buflen = 65536

func (objSrv *ObjectServer) addObjects(datas [][]byte,
	expectedHashes []*hash.Hash) ([]hash.Hash, error) {
	hashes := make([]hash.Hash, len(datas))
	numAdded := 0
	for index, data := range datas {
		var err error
		var add bool
		hashes[index], add, err = objSrv.addObject(data, expectedHashes[index])
		if err != nil {
			objSrv.logger.Printf("AddObjects(): error: %s", err.Error())
			return nil, err
		}
		if add {
			numAdded++
		}
	}
	objSrv.logger.Printf("AddObjects(): %d of %d are new objects",
		numAdded, len(datas))
	return hashes, nil
}

func (objSrv *ObjectServer) addObject(data []byte, expectedHash *hash.Hash) (
	hash.Hash, bool, error) {
	var hash hash.Hash
	if len(data) < 1 {
		return hash, false, errors.New("zero length object cannot be added")
	}
	hasher := sha512.New()
	if hasher.Size() != len(hash) {
		return hash, false, errors.New("Incompatible hash size")
	}
	_, err := hasher.Write(data)
	if err != nil {
		return hash, false, err
	}
	copy(hash[:], hasher.Sum(nil))
	if expectedHash != nil {
		if hash != *expectedHash {
			return hash, false, errors.New(fmt.Sprintf(
				"Hash mismatch. Computed=%x, expected=%x", hash, *expectedHash))
		}
	}
	filename := path.Join(objSrv.baseDir, objectcache.HashToFilename(hash))
	// Check for existing object and collision.
	fi, err := os.Lstat(filename)
	if err == nil {
		if !fi.Mode().IsRegular() {
			return hash, false, errors.New("Existing non-file: " + filename)
		}
		err := collisionCheck(data, filename, fi.Size())
		if err != nil {
			return hash, false, errors.New("Collision detected: " + err.Error())
		}
		// No collision and no error: it's the same object. Go home early.
		return hash, false, nil
	}
	err = os.MkdirAll(path.Dir(filename), 0755)
	if err != nil {
		return hash, false, err
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		return hash, false, err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return hash, false, err
	}
	objSrv.sizesMap[hash] = uint64(len(data))
	return hash, true, nil
}

func collisionCheck(data []byte, filename string, size int64) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if int64(len(data)) != size {
		return errors.New(fmt.Sprintf(
			"length mismatch. Data=%d, existing object=%d",
			len(data), size))
	}
	reader := bufio.NewReader(file)
	buffer := make([]byte, 0, buflen)
	for len(data) > 0 {
		numToRead := len(data)
		if numToRead > cap(buffer) {
			numToRead = cap(buffer)
		}
		buf := buffer[:numToRead]
		nread, err := reader.Read(buf)
		if err != nil {
			return err
		}
		if bytes.Compare(data[:nread], buf[:nread]) != 0 {
			return errors.New("content mismatch")
		}
		data = data[nread:]
	}
	return nil
}
