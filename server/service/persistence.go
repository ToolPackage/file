package service

import (
	constants "github.com/ToolPackage/fse/server/common"
	"path"
)

// PersistenceService responds for persist file chunks
type PersistenceService struct {
	rootPath  string
	indexFile *SequentialFile
	dataFile  *SequentialFile
}

func NewPersistenceService(rootPath string) (p *PersistenceService, err error) {
	p = &PersistenceService{
		rootPath: rootPath,
	}

	p.indexFile, err = NewSequentialFile(path.Join(rootPath, constants.IndexFileName), 0, 0)
	p.dataFile, err = NewSequentialFile(path.Join(rootPath, constants.DataFileName), 0, 0)

	if err != nil {
		p = nil
	}
	return
}

func (p *PersistenceService) SaveFile(file *ChunkedFile) {
	for _, chunk := range file.chunks {
		p.dataFile.Append(chunk.content)
	}
}
