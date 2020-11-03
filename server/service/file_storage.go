package service

import (
	"errors"
	"github.com/ToolPackage/fse/utils"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
)

type FileStorage struct {
	storagePath string
	files       map[string]*File
	dataFiles   []*SequentialFile
	cache       PartitionCache
}

type File struct {
	fileName    string // 128
	fileSize    int64  // 8
	contentType string // 32
	createdAt   int64  // 8
	partitions  Partitions
}

//PartitionId = sequential file id + file chunk id
type PartitionId int32
type Partitions []PartitionId

func NewFileStorage() (fs *FileStorage, err error) {
	storagePath := getStoragePath()

	// scan storage path and open all sequential files
	dataFilePath := path.Join(storagePath, "datafiles")
	fileNames, err := filepath.Glob(dataFilePath)
	if err != nil {
		return
	}

	dataFiles := make([]*SequentialFile, len(fileNames))
	for _, fileName := range fileNames {
		id, err := strconv.ParseInt(fileName, 10, 16)
		if err != nil {
			return
		}

		dataFile, err := NewSequentialFile(path.Join(dataFilePath, fileName),
			MaxFileChunkDataSize, MaxFileChunkNum)
		dataFiles[id] = dataFile
	}

	// read file metadata

	fs = &FileStorage{
		storagePath: storagePath,
		dataFiles:   dataFiles,
	}

	return
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

func (fs *FileStorage) OpenStream(file *File) io.Reader {
	return newFileDataReader(fs, file.partitions)
}

func (fs *FileStorage) SaveFileData(input io.Reader) (Partitions, error) {
	return nil, nil
}

func (fs *FileStorage) GetChunk(id PartitionId) (*FileChunk, error) {
	return nil, nil
}

type FileDataReader struct {
	fs               *FileStorage
	partitions       Partitions
	nextPartitionIdx int
	currentChunk     *FileChunk
	chunkReadOffset  int
}

func newFileDataReader(fs *FileStorage, partitions Partitions) *FileDataReader {
	return &FileDataReader{
		fs:               fs,
		partitions:       partitions,
		nextPartitionIdx: -1,
		currentChunk:     nil,
		chunkReadOffset:  0,
	}
}

func (r *FileDataReader) Read(p []byte) (n int, err error) {
	chunk, err := r.getAvailableChunk()
	if err != nil {
		return
	}

	availableBytes := len(chunk.content) - r.chunkReadOffset

	n = utils.Min(availableBytes, len(p))
	result := copy(p, chunk.content[r.chunkReadOffset:r.chunkReadOffset+n])
	if n != result {
		n = result
		err = errors.New("copy chunk data error")
	}
	r.chunkReadOffset += n
	return
}

// return nil when all chunks are consumed or chunk couldn't be load by file storage
func (r *FileDataReader) getAvailableChunk() (*FileChunk, error) {
	var err error
	if r.currentChunk == nil || r.chunkReadOffset >= len(r.currentChunk.content) {
		// get next chunk
		r.nextPartitionIdx++
		if r.nextPartitionIdx >= len(r.partitions) {
			err = EndOfPartitionStreamError
		} else {
			r.currentChunk, err = r.fs.GetChunk(r.partitions[r.nextPartitionIdx])
		}
	}
	return r.currentChunk, err
}
