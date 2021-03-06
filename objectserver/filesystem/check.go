package filesystem

import (
	"errors"
	"fmt"
	"github.com/Symantec/Dominator/lib/hash"
	"github.com/Symantec/Dominator/lib/objectcache"
	"os"
	"path"
)

func (objSrv *ObjectServer) checkObjects(hashes []hash.Hash) ([]uint64, error) {
	sizesList := make([]uint64, len(hashes))
	for index, hash := range hashes {
		var err error
		sizesList[index], err = objSrv.checkObject(hash)
		if err != nil {
			return nil, err
		}
	}
	return sizesList, nil
}

func (objSrv *ObjectServer) checkObject(hash hash.Hash) (uint64, error) {
	if size, ok := objSrv.sizesMap[hash]; ok {
		return size, nil
	}
	filename := path.Join(objSrv.baseDir, objectcache.HashToFilename(hash))
	fi, err := os.Lstat(filename)
	if err != nil {
		return 0, nil
	}
	if fi.Mode().IsRegular() {
		if fi.Size() < 1 {
			return 0, errors.New(fmt.Sprintf("zero length file: %s", filename))
		}
		size := uint64(fi.Size())
		objSrv.sizesMap[hash] = size
		return size, nil
	}
	return 0, nil
}
