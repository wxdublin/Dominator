package filesystem

import (
	"path"
)

func (fs *FileSystem) rebuildInodePointers() {
	fs.Directory.rebuildInodePointers(fs)
}

func (directory *Directory) rebuildInodePointers(fs *FileSystem) {
	for _, entry := range directory.RegularFileList {
		entry.rebuildInodePointers(fs)
	}
	for _, entry := range directory.SymlinkList {
		entry.rebuildInodePointers(fs)
	}
	for _, entry := range directory.FileList {
		entry.rebuildInodePointers(fs)
	}
	for _, entry := range directory.DirectoryList {
		entry.rebuildInodePointers(fs)
	}
}

func (file *RegularFile) rebuildInodePointers(fs *FileSystem) {
	file.inode = fs.RegularInodeTable[file.InodeNumber]
}

func (symlink *Symlink) rebuildInodePointers(fs *FileSystem) {
	symlink.inode = fs.SymlinkInodeTable[symlink.InodeNumber]
}

func (file *File) rebuildInodePointers(fs *FileSystem) {
	file.inode = fs.InodeTable[file.InodeNumber]
}

func (fs *FileSystem) buildFilenamesTable() {
	fs.FilenamesTable = make(FilenamesTable)
	fs.Directory.buildFilenamesTable(fs, "")
}

func (directory *Directory) buildFilenamesTable(fs *FileSystem,
	parentName string) {
	myPathName := path.Join(parentName, directory.Name)
	for _, entry := range directory.RegularFileList {
		fs.addFilenameToTable(entry.InodeNumber, myPathName, entry.Name)
	}
	for _, entry := range directory.SymlinkList {
		fs.addFilenameToTable(entry.InodeNumber, myPathName, entry.Name)
	}
	for _, entry := range directory.FileList {
		fs.addFilenameToTable(entry.InodeNumber, myPathName, entry.Name)
	}
	for _, entry := range directory.DirectoryList {
		entry.buildFilenamesTable(fs, myPathName)
	}
}

func (fs *FileSystem) addFilenameToTable(inode uint64,
	parentName, entryName string) {
	filenames := fs.FilenamesTable[inode]
	filenames = append(filenames, path.Join(parentName, entryName))
	fs.FilenamesTable[inode] = filenames
}

func (fs *FileSystem) computeTotalDataBytes() {
	fs.TotalDataBytes = 0
	for _, inode := range fs.RegularInodeTable {
		fs.TotalDataBytes += uint64(inode.Size)
	}
}

func (directory *Directory) buildEntryMap() {
	directory.EntriesByName = make(map[string]interface{})
	for _, entry := range directory.RegularFileList {
		directory.EntriesByName[entry.Name] = entry
	}
	for _, entry := range directory.SymlinkList {
		directory.EntriesByName[entry.Name] = entry
	}
	for _, entry := range directory.FileList {
		directory.EntriesByName[entry.Name] = entry
	}
	for _, entry := range directory.DirectoryList {
		directory.EntriesByName[entry.Name] = entry
		entry.buildEntryMap()
	}
}
