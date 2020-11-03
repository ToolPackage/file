package service

import (
	"github.com/ToolPackage/fse/server/common"
)

type File struct {
	indexFile *SequentialFile
	dataFile  *SequentialFile
}

const IndexFileSuffix = ".fseidx"
const DataFileSuffix = ".fsedata"

func NewFile(filePath string) (f *File, err error) {
	var indexFile *SequentialFile
	if indexFile, err = NewSequentialFile(filePath+IndexFileSuffix,
		common.MaxSequentialFileSize/common.FileChunkSize, 0); err != nil {
		return
	}

	var dataFile *SequentialFile
	if dataFile, err = NewSequentialFile(filePath+DataFileSuffix,
		common.MaxSequentialFileSize, 0); err != nil {
		return
	}

	f = &File{
		indexFile: indexFile,
		dataFile:  dataFile,
	}
	return
}
