package storage

import (
	"github.com/ToolPackage/fse/utils"
	"io"
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

func NewEntrySequenceFile(path string, mode int) *EntrySequenceFile {
	var file *os.File

	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
		if mode == ReadMode {
			panic(err)
		}
		// create and open file
		if file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			panic(err)
		}
	} else {
		mode := os.O_RDONLY
		if mode == WriteMode {
			mode = os.O_WRONLY
		}

		// open file
		if file, err = os.OpenFile(path, mode, 0644); err != nil {
			panic(err)
		}
	}

	return &EntrySequenceFile{
		path:         path,
		mode:         mode,
		file:         file,
		entrySizeBuf: make([]byte, entrySizeSize),
	}
}

func (f *EntrySequenceFile) WriteEntry(entry []byte) {
	if f.mode == ReadMode {
		panic(InvalidOperationError)
	}
	entrySize := len(entry)
	if entrySize > maxEntrySize {
		panic(InvalidOperationError)
	}

	utils.ConvertUint16ToByte(uint16(entrySize), f.entrySizeBuf, 0)
	n, err := f.file.Write(f.entrySizeBuf)
	if n != entrySizeSize || err != nil {
		panic(err)
	}

	n, err = f.file.Write(entry)
	if n != entrySize || err != nil {
		panic(err)
	}
}

func (f *EntrySequenceFile) ReadEntry() []byte {
	if f.mode == WriteMode {
		panic(InvalidOperationError)
	}

	n, err := f.file.Read(f.entrySizeBuf)
	if n != entrySizeSize || err != nil {
		if err == io.EOF {
			return nil
		}
		panic(err)
	}

	entrySize := utils.ConvertByteToUint16(f.entrySizeBuf, 0)
	buf := make([]byte, entrySize)
	n, err = f.file.Read(buf)
	if n != int(entrySize) || err != nil {
		if err == io.EOF {
			return nil
		}
		panic(err)
	}

	return buf
}

func (f *EntrySequenceFile) Close() {
	if err := f.file.Close(); err != nil {
		panic(err)
	}
}
