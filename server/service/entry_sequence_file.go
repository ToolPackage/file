package service

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

const entrySizeSize = 2         // unit: bytes
const maxEntrySize = 0xffffffff // 16bits
const (
	ReadMode = iota
	WriteMode
)

func NewEntrySequenceFile(path string, mode int) (f *EntrySequenceFile, err error) {
	var file *os.File

	if _, err = os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		if mode == ReadMode {
			return
		}
		// create and open file
		if file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return
		}
	} else {
		mode := os.O_RDONLY
		if mode == WriteMode {
			mode = os.O_WRONLY
		}

		// open file
		if file, err = os.OpenFile(path, mode, 0644); err != nil {
			return
		}
	}

	f = &EntrySequenceFile{
		path:         path,
		mode:         mode,
		file:         file,
		entrySizeBuf: make([]byte, entrySizeSize),
	}
	return
}

func (f *EntrySequenceFile) WriteEntry(entry []byte) error {
	if f.mode == ReadMode {
		return InvalidOperationError
	}

	entrySize := len(entry)
	if entrySize > maxEntrySize {
		return EntryTooLargeError
	}

	utils.ConvertInt16ToByte(uint16(entrySize), f.entrySizeBuf, 0)
	n, err := f.file.Write(f.entrySizeBuf)
	if n != entrySizeSize || err != nil {
		return err
	}

	n, err = f.file.Write(entry)
	return err
}

func (f *EntrySequenceFile) ReadEntry() ([]byte, error) {
	if f.mode == WriteMode {
		return nil, InvalidOperationError
	}

	n, err := f.file.Read(f.entrySizeBuf)
	if n != entrySizeSize || err != nil {
		return nil, err
	}

	entrySize := utils.ConvertByteToInt16(f.entrySizeBuf, 0)
	buf := make([]byte, entrySize)
	n, err = f.file.Read(buf)
	if n != entrySizeSize || err != nil {
		return nil, err
	}

	return buf, nil
}

func (f *EntrySequenceFile) Close() error {
	return f.file.Close()
}
