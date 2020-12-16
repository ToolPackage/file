package storage

import (
	"crypto/md5"
	"errors"
	"github.com/ToolPackage/fse/common/utils"
	"os"
)

// all bytes
const MaxFileChunkDataSize = 64*1024 - 1   // (64kb = 65536) overflow uint16
const SequentialFileChunkMetadataSize = 18 // md5 + chunkSize
const MaxFileChunkNum = 8192

const SequentialFileMetadataSize = 8        // chunkSize + chunkCap + chunkNum
const MaxSequentialFileDataSize = 537010176 // almost 512MB, = (MaxFileChunkDataSize + SequentialFileChunkMetadataSize) * MaxFileChunkNum

type SequentialFile struct {
	path            string
	file            *os.File
	chunkSize       uint16
	chunkCap        uint16
	deletedChunkNum uint16
	chunkNum        uint16
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
	s = &SequentialFile{path: path}

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
		if s.file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644); err != nil {
			return
		}
		// re-alloc file space to totalSize, totalSize is the file size in bytes
		totalSize := int64(totalDataSize) + int64(SequentialFileMetadataSize)
		if err = os.Truncate(path, totalSize); err != nil {
			return
		}
		// init metadata
		s.chunkSize = chunkSize
		s.chunkCap = chunkCap
		s.chunkNum = 0
		s.deletedChunkNum = 0
		// write metadata
		if err = s.writeMetadata(); err != nil {
			return
		}
	} else {
		// file exists
		if s.file, err = os.OpenFile(path, os.O_RDWR, 0644); err != nil {
			return
		}
		// read metadata
		if err = s.readMetadata(); err != nil {
			return
		}
	}

	return
}

func (s *SequentialFile) writeMetadata() (err error) {
	// seek cursor to the last 8 bytes of the file
	_, err = s.file.Seek(-SequentialFileMetadataSize, 2)
	if err != nil {
		return
	}

	buf := make([]byte, SequentialFileMetadataSize)
	utils.ConvertUint16ToByte(s.chunkSize, buf, 0)
	utils.ConvertUint16ToByte(s.chunkCap, buf, 2)
	utils.ConvertUint16ToByte(s.deletedChunkNum, buf, 4)
	utils.ConvertUint16ToByte(s.chunkNum, buf, 6)
	n, err := s.file.Write(buf)
	if err != nil {
		return
	}
	if n != SequentialFileMetadataSize {
		err = InvalidRetValue
	}
	return
}

func (s *SequentialFile) readMetadata() (err error) {
	// seek cursor relating to end fo the file
	_, err = s.file.Seek(-SequentialFileMetadataSize, 2)
	if err != nil {
		return
	}

	buf := make([]byte, SequentialFileMetadataSize)
	n, err := s.file.Read(buf)
	if n != len(buf) || err != nil {
		return
	}

	s.chunkSize = utils.ConvertByteToUint16(buf, 0)
	s.chunkCap = utils.ConvertByteToUint16(buf, 2)
	s.deletedChunkNum = utils.ConvertByteToUint16(buf, 4)
	s.chunkNum = utils.ConvertByteToUint16(buf, 6)
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

	// read chunk md5
	var md5Bytes [16]byte
	n, err = s.file.Read(md5Bytes[:])
	if err != nil {
		return
	}
	if n != 16 {
		err = InvalidRetValue
		return
	}
	// read chunk size
	buf := make([]byte, 2)
	n, err = s.file.Read(buf)
	if err != nil {
		return
	}
	if n != 2 {
		err = InvalidRetValue
		return
	}
	chunkSize := utils.ConvertByteToUint16(buf, 0)
	// read chunk data
	buf = make([]byte, chunkSize)
	n, err = s.file.Read(buf)
	if err != nil {
		return
	}
	if n != int(chunkSize) {
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

func (s *SequentialFile) IsWritable() bool {
	return s.chunkNum < s.chunkCap
}

// Append data to new chunk
// DataOutOfFileError,
// DataOutOfChunkError,
// seek failure,
// write failure,
// InvalidRetValue
// TODO: bug 写文件空间不足
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
	// write md5
	md5Bytes := md5.Sum(data)
	n, err = s.file.Write(md5Bytes[:])
	if err != nil {
		return
	}
	if n != 16 {
		err = InvalidRetValue
		return
	}
	// write data size
	buf := make([]byte, 2)
	utils.ConvertUint16ToByte(uint16(len(data)), buf, 0)
	n, err = s.file.Write(buf)
	if err != nil {
		return
	}
	if n != 2 {
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

func (s *SequentialFile) DeleteChunk(chunkId uint16) error {
	if chunkId >= s.chunkNum {
		return InvalidChunkIdError
	}

	s.deletedChunkNum++
	if s.deletedChunkNum == s.chunkNum {
		// when all chunks are deleted,
		// we could append new chunk to the start of the file now
		s.deletedChunkNum = 0
		s.chunkNum = 0
	}
	return nil
}

func (c *FileChunk) Validate() bool {
	return c.md5 == md5.Sum(c.content)
}

func (s *SequentialFile) Close() (err error) {
	if err = s.writeMetadata(); err != nil {
		return
	}

	if err = s.file.Close(); err != nil {
		return
	}
	//s.path = ""
	s.chunkNum = 0
	s.chunkSize = 0
	s.file = nil

	return
}

func (s *SequentialFile) Delete() error {
	if _, err := os.Stat(s.path); err != nil {
		return err
	}
	return os.Remove(s.path)
}
