package service

import (
	"fmt"
	log "github.com/Luncert/slog"
	constants "github.com/ToolPackage/fse/server/common"
	"os"
)

type SequentialFile struct {
	path         string
	file         *os.File
	totalSize    int64 // the pre-alloc size of the file, normally it's MaxSequentialFileSize
	appendOffset int64 // the offset of the append cursor
}

func NewSequentialFile(path string, totalSize int64, appendOffset int64) (*SequentialFile, error) {
	// the totalSize and appendOffset are stored in Mongo
	s := &SequentialFile{path: path, totalSize: totalSize, appendOffset: appendOffset}

	var err error
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// create and open file
			s.file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
			// re-alloc file space to MaxDataFileSize
			err = os.Truncate(path, constants.MaxSequentialFileSize)
			return s, err
		}
	}

	// file is already created
	s.file, err = os.OpenFile(path, os.O_RDWR, 0644)
	return s, err
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
	if s.totalSize < int64(len(data))+s.appendOffset+1 {
		return 0, fmt.Errorf("no enough space to write data")
	}

	// seek file
	_, err = s.file.Seek(s.appendOffset, 0) // 0 indicates referring the start of the file
	if err != nil {
		return 0, err
	}

	n, err = s.file.Write(data)
	return
}

func (s *SequentialFile) Close() {
	if err := s.file.Close(); err != nil {
		log.Error(err)
	}
	s.path = ""
	s.totalSize = 0
	s.appendOffset = 0
	s.file = nil
}
