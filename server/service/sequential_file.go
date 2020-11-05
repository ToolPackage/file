package service

import (
	"crypto/md5"
	"errors"
	"github.com/ToolPackage/fse/utils"
	"os"
)

// all bytes
const MaxFileChunkDataSize = 64*1024 - 1   // (64kb = 65536) overflow uint16
const SequentialFileChunkMetadataSize = 16 // md5
const MaxFileChunkNum = 8192

const SequentialFileMetadataSize = 6        // chunkSize + chunkCap + chunkNum
const MaxSequentialFileDataSize = 536993792 // almost 512MB, = (MaxFileChunkDataSize + SequentialFileChunkMetadataSize) * MaxFileChunkNum

type SequentialFile struct {
	path      string
	file      *os.File
	chunkSize uint16
	chunkCap  uint16
	chunkNum  uint16
}

type FileChunk struct {
	chunkId uint16
	md5     [16]byte
	content []byte
}

// Open or create sequential file, if target file exists,
// chunkSize and chunkNum will be read from file's metadata
// instead of using function arguments.
func NewSequentialFile(path string, chunkSize uint16, chunkCap uint16) (s *SequentialFile, err error) {

	var file *os.File
	var chunkNum uint16 = 0

	if _, err = os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		// only validate arguments when file doesn't exist
		realChunkSize := uint32(chunkSize) + uint32(SequentialFileChunkMetadataSize)
		totalDataSize := realChunkSize * uint32(chunkCap)
		if totalDataSize > MaxSequentialFileDataSize || totalDataSize == 0 {
			err = errors.New("invalid chunkSize or chunkNum")
			return
		}
		// create and open file
		if file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644); err != nil {
			return
		}
		// re-alloc file space to totalSize, totalSize is the file size in bytes
		totalSize := int64(totalDataSize) + int64(SequentialFileMetadataSize)
		if err = os.Truncate(path, totalSize); err != nil {
			return
		}
		// write metadata
		if err = writeMetadata(file, chunkSize, chunkCap, chunkNum); err != nil {
			return
		}
	} else {
		// file exists
		if file, err = os.OpenFile(path, os.O_RDWR, 0644); err != nil {
			return
		}
		// read metadata
		chunkSize, chunkCap, chunkNum, err = readMetadata(file)
		if err != nil {
			return
		}
	}

	s = &SequentialFile{
		path:      path,
		file:      file,
		chunkSize: chunkSize,
		chunkCap:  chunkCap,
		chunkNum:  chunkNum,
	}
	return
}

func writeMetadata(f *os.File, chunkSize uint16, chunkCap uint16, chunkNum uint16) (err error) {
	// seek cursor relating to end fo the file
	_, err = f.Seek(-SequentialFileMetadataSize, 2)
	if err != nil {
		return
	}

	buf := make([]byte, SequentialFileMetadataSize)
	utils.ConvertUint16ToByte(chunkSize, buf, 0)
	utils.ConvertUint16ToByte(chunkCap, buf, 2)
	utils.ConvertUint16ToByte(chunkNum, buf, 4)
	n, err := f.Write(buf)
	if n != SequentialFileMetadataSize {
		err = InvalidRetValue
	}
	return
}

func readMetadata(f *os.File) (chunkSize uint16, chunkCap uint16, chunkNum uint16, err error) {
	// seek cursor relating to end fo the file
	_, err = f.Seek(-SequentialFileMetadataSize, 2)
	if err != nil {
		return
	}

	buf := make([]byte, SequentialFileMetadataSize)
	n, err := f.Read(buf)
	if n != len(buf) || err != nil {
		return
	}

	chunkSize = utils.ConvertByteToUint16(buf, 0)
	chunkCap = utils.ConvertByteToUint16(buf, 2)
	chunkNum = utils.ConvertByteToUint16(buf, 4)
	return
}

// read block at specified position in the file
func (s *SequentialFile) ReadChunk(chunkId uint16) (chunk *FileChunk, err error) {
	if chunkId >= s.chunkNum {
		err = InvalidChunkIdError
		return
	}

	offset := int64(chunkId) * (int64(s.chunkSize) + int64(SequentialFileChunkMetadataSize))
	_, err = s.file.Seek(offset, 0) // seek from the start
	if err != nil {
		return nil, err
	}

	var n int

	// read chunk metadata
	var md5Bytes [16]byte
	n, err = s.file.Read(md5Bytes[:])
	if err != nil {
		return
	}
	if n != 16 {
		err = InvalidRetValue
		return
	}

	// read chunk data
	buf := make([]byte, s.chunkSize)
	n, err = s.file.Read(buf)
	if err != nil {
		return
	}
	if n != int(s.chunkSize) {
		err = InvalidRetValue
		return
	}

	chunk = &FileChunk{
		chunkId: chunkId,
		md5:     md5Bytes,
		content: buf,
	}
	return
}

// Append data to new chunk
// DataOutOfFileError,
// DataOutOfChunkError,
// seek failure,
// write failure,
// InvalidRetValue
func (s *SequentialFile) AppendChunk(data []byte) (chunkId uint16, err error) {
	if s.chunkNum >= s.chunkCap {
		err = DataOutOfFileError
		return
	}
	if len(data) > int(s.chunkSize) {
		err = DataOutOfChunkError
		return
	}

	// fit data to chunkSize to calculate md5

	chunkId = s.chunkNum

	// seek file
	offset := int64(chunkId) * (int64(s.chunkSize) + int64(SequentialFileChunkMetadataSize))
	// 0 indicates referring the start of the file
	if _, err = s.file.Seek(offset, 0); err != nil {
		return
	}

	var n int
	// write metadata
	md5Bytes := md5.Sum(data)
	n, err = s.file.Write(md5Bytes[:])
	if err != nil {
		return
	}
	if n != 16 {
		err = InvalidRetValue
		return
	}
	// write chunk data
	n, err = s.file.Write(data)
	if err != nil {
		return
	}
	if n != len(data) {
		err = InvalidRetValue
		return
	}

	s.chunkNum++
	return
}

func (s *SequentialFile) Close() (err error) {
	if err = writeMetadata(s.file, s.chunkSize, s.chunkCap, s.chunkNum); err != nil {
		return
	}

	if err = s.file.Close(); err != nil {
		return
	}
	s.path = ""
	s.chunkNum = 0
	s.chunkSize = 0
	s.file = nil

	return
}

func (c *FileChunk) Validate() bool {
	return c.md5 == md5.Sum(c.content)
}
