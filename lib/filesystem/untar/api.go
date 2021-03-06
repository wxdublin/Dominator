package untar

import (
	"archive/tar"
	"github.com/Symantec/Dominator/lib/filesystem"
	"github.com/Symantec/Dominator/lib/filter"
	"github.com/Symantec/Dominator/lib/hash"
)

type DataHandler interface {
	HandleData(data []byte) (hash.Hash, error)
}

func Decode(tarReader *tar.Reader, dataHandler DataHandler,
	filter *filter.Filter) (*filesystem.FileSystem, error) {
	return decode(tarReader, dataHandler, filter)
}
