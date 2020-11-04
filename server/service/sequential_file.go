package service

import (
	"crypto/md5"
	"errors"
	"github.com/ToolPackage/fse/utils"
	"os"
)

const MaxFileChunkDataSize = 64*1024 - 1                                 // (64kb = 65536) overflow uint16
const MaxSequentialFileDataSize = 512 * 1024 * 1024                      // 512MB
const MaxFileChunkNum = MaxSequentialFileDataSize / MaxFileChunkDataSize // 8192
const SequentialFileMetadataSize = 32                                    // chunkSize + chunkNum, unit: bytes
const SequentialFileChunkMetadataSize = 16                               // md5, unit: bytes

type SequentialFile struct {
	path      string
	file      *os.File
	chunkSize uint16
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
func NewSequentialFile(path string, chunkSize uint16, chunkNum uint16) (s *SequentialFile, err error) {
	realChunkSize := uint32(chunkSize) + uint32(SequentialFileChunkMetadataSize)
	totalDataSize := realChunkSize * uint32(chunkNum)
	totalSize := int64(totalDataSize) + int64(SequentialFileMetadataSize)

	if totalDataSize > MaxSequentialFileDataSize || totalDataSize == 0 {
		err = errors.New("invalid chunkSize or chunkNum")
		return
	}

	var file *os.File

	if _, err = os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		// create and open file
		if file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644); err != nil {
			return
		}
		// re-alloc file space to totalSize
		if err = os.Truncate(path, totalSize); err != nil {
			return
		}
		// write metadata
		if err = writeMetadata(file, chunkSize, chunkNum); err != nil {
			return
		}
	} else {
		// file exists
		if file, err = os.OpenFile(path, os.O_RDWR, 0644); err != nil {
			return
		}
		// read metadata
		chunkSize, chunkNum, err = readMetadata(file)
		if err != nil {
			return
		}
	}

	s = &SequentialFile{
		path:      path,
		file:      file,
		chunkSize: chunkSize,
		chunkNum:  chunkNum,
	}
	return
}

func writeMetadata(f *os.File, chunkSize uint16, chunkNum uint16) (err error) {
	var metadataByteNum int64 = SequentialFileMetadataSize / 8
	// seek cursor relating to end fo the file
	_, err = f.Seek(-metadataByteNum, 2)
	if err != nil {
		return
	}

	buf := make([]byte, metadataByteNum)
	utils.ConvertUint16ToByte(chunkSize, buf, 0)
	utils.ConvertUint16ToByte(chunkNum, buf, 1)
	return
}

func readMetadata(f *os.File) (chunkSize uint16, chunkNum uint16, err error) {
	var metadataByteNum int64 = SequentialFileMetadataSize / 8
	// seek cursor relating to end fo the file
	_, err = f.Seek(-metadataByteNum, 2)
	if err != nil {
		return
	}

	buf := make([]byte, metadataByteNum)
	n, err := f.Read(buf)
	if n != len(buf) || err != nil {
		return
	}

	chunkSize = utils.ConvertByteToUint16(buf, 0)
	chunkNum = utils.ConvertByteToUint16(buf, 2)
	return
}

// read block at specified position in the file
func (s *SequentialFile) ReadChunk(chunkId uint16) (chunk *FileChunk, err error) {
	offset := int64(chunkId * s.chunkSize)
	_, err = s.file.Seek(offset, 0) // seek from the start
	if err != nil {
		return nil, err
	}

	var n int
	buf := make([]byte, s.chunkSize)

	// read chunk metadata
	var md5Bytes [16]byte
	n, err = s.file.Read(buf)
	if n != len(buf) || err != nil {
		return
	}
	for i := 0; i < 16; i++ {
		md5Bytes[i] = buf[i]
	}

	// read chunk data
	n, err = s.file.Read(buf)
	if n != len(buf) || err != nil {
		return
	}

	chunk = &FileChunk{
		chunkId: chunkId,
		md5:     md5Bytes,
		content: buf[SequentialFileChunkMetadataSize:],
	}
	return
}

// Append data to new chunk
func (s *SequentialFile) AppendChunk(data []byte) (chunkId uint16, err error) {
	if len(data) > int(s.chunkSize) {
		err = DataOutOfChunkError
		return
	}

	chunkId = s.chunkNum + 1

	// seek file
	offset := int64(chunkId * s.chunkSize)
	// 0 indicates referring the start of the file
	if _, err = s.file.Seek(offset, 0); err != nil {
		return
	}

	var n int
	md5Bytes := md5.Sum(data)
	n, err = s.file.Write(md5Bytes[:])
	if err != nil || n != 16 {
		return
	}
	n, err = s.file.Write(data)
	if err != nil || n != len(data) {
		return
	}

	s.chunkNum++
	return
}

func (s *SequentialFile) Close() error {
	if err := s.file.Close(); err != nil {
		return err
	}
	s.path = ""
	s.chunkNum = 0
	s.chunkSize = 0
	s.file = nil

	return nil
}
