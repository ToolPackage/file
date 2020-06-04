package service

import (
	"os"
	"path"
)

const IndexFileName = "index.fse"
const DataFileName = "data.fse"
const FileBlockSize = 64 * 1024 // 64kb

type StorageService struct {
	rootPath  string
	indexFile *SequentialFile
	dataFile  *SequentialFile
}

func New(rootPath string) (s *StorageService, err error) {
	s = &StorageService{
		rootPath: rootPath,
	}

	s.indexFile, err = NewSequentialFile(path.Join(rootPath, IndexFileName))
	s.dataFile, err = NewSequentialFile(path.Join(rootPath, DataFileName))

	if err != nil {
		s = nil
	}
	return
}

type SequentialFile struct {
	path string
	file *os.File
}

func NewSequentialFile(path string) (s *SequentialFile, err error) {
	s = &SequentialFile{path: path}

	s.file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		s = nil
	}
	return
}

func (s *SequentialFile) GetFileSize() int64 {
	fileInfo, _ := os.Stat(s.path)
	return fileInfo.Size()
}

// read block at specified position in the file
func (s *SequentialFile) Read(offset int64, size int) (data []byte, n int, err error) {
	_, err = s.file.Seek(offset, 0) // seek from the start
	if err != nil {
		return nil, 0, err
	}

	data = make([]byte, size)
	n, err = s.file.Read(data)
	return
}

func (s *SequentialFile) Append(data []byte) (n int, err error) {
	_, err = s.file.Seek(0, 2)
	if err != nil {
		return 0, err
	}

	n, err = s.file.Write(data)
	return
}
