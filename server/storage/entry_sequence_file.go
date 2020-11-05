package storage

import (
	"github.com/ToolPackage/fse/utils"
	"os"
)

type EntrySequenceFile struct {
	path         string
	mode         int
	file         *os.File
	entrySizeBuf []byte
}

const entrySizeSize = 2     // unit: bytes
const maxEntrySize = 0xffff // 16bits
const (
	ReadMode = iota
	WriteMode
)

func NewEntrySequenceFile(path string, mode int) (*EntrySequenceFile, error) {
	var file *os.File

	m := os.O_WRONLY
	if mode == ReadMode {
		m = os.O_RDONLY
	}

	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// create file
		if file, err = os.OpenFile(path, os.O_CREATE|m, 0644); err != nil {
			return nil, err
		}
	} else {
		// open file
		if file, err = os.OpenFile(path, m, 0644); err != nil {
			return nil, err
		}
	}

	return &EntrySequenceFile{
		path:         path,
		mode:         mode,
		file:         file,
		entrySizeBuf: make([]byte, entrySizeSize),
	}, nil
}

func (f *EntrySequenceFile) WriteEntry(entry []byte) error {
	if f.mode == ReadMode {
		return InvalidOperationError
	}
	entrySize := len(entry)
	if entrySize > maxEntrySize {
		return EntryTooLargeError
	}

	utils.ConvertUint16ToByte(uint16(entrySize), f.entrySizeBuf, 0)
	n, err := f.file.Write(f.entrySizeBuf)
	if err != nil {
		return err
	}
	if n != entrySizeSize {
		return InvalidRetValue
	}

	n, err = f.file.Write(entry)
	if err != nil {
		return err
	}
	if n != entrySize {
		return InvalidRetValue
	}

	return nil
}

// InvalidOperationError
// io.EOF
// others
// InvalidRetValue
func (f *EntrySequenceFile) ReadEntry() ([]byte, error) {
	if f.mode == WriteMode {
		return nil, InvalidOperationError
	}

	n, err := f.file.Read(f.entrySizeBuf)
	if err != nil {
		// include EOF
		return nil, err
	}
	if n != entrySizeSize {
		return nil, InvalidRetValue
	}

	entrySize := utils.ConvertByteToUint16(f.entrySizeBuf, 0)
	buf := make([]byte, entrySize)
	n, err = f.file.Read(buf)
	if err != nil {
		// include EOF
		return nil, err
	}
	if n != int(entrySize) {
		return nil, InvalidRetValue
	}

	return buf, nil
}

func (f *EntrySequenceFile) Close() error {
	return f.file.Close()
}
