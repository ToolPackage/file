package service

import (
	constants "github.com/ToolPackage/fse/server/common"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var (
	PersistenceServiceIns = NewPersistenceService(getStoragePath())
)

// PersistenceService responds for persist file chunks
type PersistenceService struct {
	rootPath  string
	indexFile *SequentialFile
	dataFile  *SequentialFile
}

type FileInfo struct {
	FileId   string
	FileName string
	FileSize int64
}

func NewPersistenceService(rootPath string) (p *PersistenceService) {
	var err error
	var indexFile *SequentialFile
	if indexFile, err = NewSequentialFile(path.Join(rootPath, constants.IndexFileName), 0, 0); err != nil {
		panic(err)
	}

	var dataFile *SequentialFile
	if dataFile, err = NewSequentialFile(path.Join(rootPath, constants.DataFileName), 0, 0); err != nil {
		panic(err)
	}

	p = &PersistenceService{
		rootPath:  rootPath,
		indexFile: indexFile,
		dataFile:  dataFile,
	}
	return
}

func (p *PersistenceService) SaveFile(fileName string, contentType string, data io.Reader) *FileInfo {
	//file := &ChunkedFile{fileName: fileName,
	//	contentType:contentType,
	//	createdAt:time.Now().Unix(),
	//	chunks: []FileChunk{},
	//}
	//
	//var (
	//	chunk *FileChunk
	//	ret int
	//	err error
	//)
	//
	//for true {
	//	chunk = &FileChunk{}
	//	ret, err = data.Read(chunk.content)
	//	chunk.chunkId =
	//}
	//data.Read()
	return nil
}

func (p *PersistenceService) SaveChunkedFile(file *ChunkedFile) {
	for _, chunk := range file.chunks {
		p.dataFile.Append(chunk.content)
	}
}

func getStoragePath() string {
	return filepath.Join(getUserHomeDir(), ".fse")
}

func getUserHomeDir() string {
	home := os.Getenv("HOME")
	if home != "" {
		return home
	}

	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
		if home != "" {
			return home
		}

		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home != "" {
			return home
		}
	}

	panic("could not detect home directory")
}
