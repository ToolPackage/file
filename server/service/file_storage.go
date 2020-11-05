package service

import (
	"errors"
	log "github.com/Luncert/slog"
	"github.com/ToolPackage/fse/utils"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

// TODO: 并发
type FileStorage struct {
	storagePath string
	files       map[string]*File
	dataFiles   []*SequentialFile
	cache       PartitionCache
}

type File struct {
	fileName    string // 128
	fileSize    uint32 // 8
	contentType string // 32
	createdAt   int64  // 8
	partitions  Partitions
}

//PartitionId = sequential file id + file chunk id
type PartitionId uint32
type Partitions []PartitionId

const maxPartitionNum = 0xffff - 1 // 65535, 2Bytes
const dataFilesDirName = "datafiles"
const storageMetadataFileName = "metadata.esf"

func NewFileStorage() *FileStorage {
	storagePath := getStoragePath()

	// scan storage path and open all sequential files
	dataFilePath := path.Join(storagePath, dataFilesDirName)
	fileNames, err := filepath.Glob(dataFilePath)
	if err != nil {
		panic(err)
	}

	dataFiles := make([]*SequentialFile, len(fileNames))
	for _, fileName := range fileNames {
		id, err := strconv.ParseInt(fileName, 10, 16)
		if err != nil {
			panic(err)
		}

		dataFile, err := NewSequentialFile(path.Join(dataFilePath, fileName), 0, 0)
		dataFiles[id] = dataFile
	}

	files := loadStorageMetadata(storagePath)

	return &FileStorage{
		storagePath: storagePath,
		files:       files,
		dataFiles:   dataFiles,
	}
}

func loadStorageMetadata(storagePath string) map[string]*File {
	// TODO: bug
	defer func() {
		if err := recover(); err != nil && err == io.EOF {
			err = nil
		}
	}()

	// read file metadata
	metadataFile := NewEntrySequenceFile(path.Join(storagePath, storageMetadataFileName), ReadMode)
	defer metadataFile.Close()

	var files = make(map[string]*File)
	for true {
		file := &File{}
		file.fileName = string(metadataFile.ReadEntry())
		file.fileSize = utils.ConvertByteToUint32(metadataFile.ReadEntry(), 0)
		file.contentType = string(metadataFile.ReadEntry())
		file.createdAt = utils.ConvertByteToInt64(metadataFile.ReadEntry(), 0)
		partitionNum := utils.ConvertByteToUint16(metadataFile.ReadEntry(), 0)
		partitions := make([]PartitionId, partitionNum)
		for i := uint16(0); i < partitionNum; i++ {
			partitions[i] = PartitionId(utils.ConvertByteToUint32(metadataFile.ReadEntry(), 0))
		}
		files[file.fileName] = file
	}

	return files
}

func (fs *FileStorage) saveStorageMetadata() {
	metadataFile := NewEntrySequenceFile(path.Join(fs.storagePath, storageMetadataFileName), WriteMode)
	defer metadataFile.Close()

	var buf []byte
	for _, file := range fs.files {
		metadataFile.WriteEntry([]byte(file.fileName))

		buf = make([]byte, 4)
		utils.ConvertUint32ToByte(file.fileSize, buf, 0)
		metadataFile.WriteEntry(buf)

		metadataFile.WriteEntry([]byte(file.contentType))

		buf = make([]byte, 8)
		utils.ConvertInt64ToByte(file.createdAt, buf, 0)
		metadataFile.WriteEntry(buf)

		buf = make([]byte, 2)
		utils.ConvertUint16ToByte(uint16(len(file.partitions)), buf, 0)
		metadataFile.WriteEntry(buf)

		buf = make([]byte, 4)
		for _, id := range file.partitions {
			utils.ConvertUint32ToByte(uint32(id), buf, 0)
			metadataFile.WriteEntry(buf)
		}
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

func (fs *FileStorage) OpenStream(file *File) io.Reader {
	return newFileDataReader(fs, file.partitions)
}

func (fs *FileStorage) SaveFile(fileName string, contentType string, fileSize uint32, reader io.Reader) (*File, error) {
	file := &File{
		fileName:    fileName,
		fileSize:    fileSize,
		contentType: contentType,
		createdAt:   time.Now().UnixNano(),
		partitions:  make(Partitions, 0),
	}

	chunkBuf := make([]byte, MaxFileChunkDataSize)
	dataFile := fs.dataFiles[len(fs.dataFiles)-1]
	for true {
		// read input
		if _, err := reader.Read(chunkBuf); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if len(file.partitions) >= maxPartitionNum {
			return nil, PartitionNumLimitError
		}
		// append input to data file
		chunkId, err := dataFile.AppendChunk(chunkBuf)
		if err == DataOutOfFileError {
			if dataFile, err = fs.createDataFile(); err != nil {
				return nil, err
			}
			if chunkId, err = dataFile.AppendChunk(chunkBuf); err != nil {
				return nil, err
			}
		}
		// maintain partition info
		partitionId := createPartitionId(uint16(len(fs.dataFiles)), chunkId)
		file.partitions = append(file.partitions, partitionId)
	}
	return file, nil
}

func (fs *FileStorage) createDataFile() (*SequentialFile, error) {
	fileId := len(fs.dataFiles)
	file, err := NewSequentialFile(path.Join(fs.storagePath, dataFilesDirName, strconv.Itoa(fileId)),
		MaxFileChunkDataSize, MaxFileChunkNum)
	if err != nil {
		return nil, err
	}
	fs.dataFiles = append(fs.dataFiles, file)
	return file, nil
}

func (fs *FileStorage) GetChunk(id PartitionId) (*FileChunk, error) {
	fileId, chunkId := id.split()
	if int(fileId) >= len(fs.dataFiles) {
		return nil, InvalidPartitionIdError
	}
	file := fs.dataFiles[fileId]
	if chunkId >= file.chunkNum {
		return nil, InvalidPartitionIdError
	}

	return file.ReadChunk(chunkId)
}

func (fs *FileStorage) Destroy() {
	fs.saveStorageMetadata()
	fs.files = nil
	for _, file := range fs.dataFiles {
		if err := file.Close(); err != nil {
			log.Info("failed to close sequential file handle, path = ", file.path, ", err = ", err)
		}
	}
	fs.cache.Destroy()
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

func createPartitionId(fileId uint16, chunkId uint16) PartitionId {
	return PartitionId(uint32(fileId)<<16 + uint32(chunkId))
}

func (id PartitionId) split() (uint16, uint16) {
	return uint16((id >> 16) & 0xffff), uint16(id & 0xffff)
}
