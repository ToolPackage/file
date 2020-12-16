package storage

import (
	"errors"
	"github.com/ToolPackage/fse/utils"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var S = NewFileStorage()

func init() {
	log.Println("file storage init")
	runtime.SetFinalizer(S, func(fs *FileStorage) {
		// TODO: won't be invoked
		fs.Destroy()
		log.Println("file storage stopped gracefully")
	})
}

// TODO: 并发
type FileStorage struct {
	storagePath string
	files       map[string]*File
	dataFiles   []*SequentialFile
	cache       PartitionCache
}

const maxPartitionNum = 0xffff - 1 // 65535, 2Bytes
const storageDirName = ".fse"
const dataFilesDirName = "datafiles"
const storageMetadataFileName = "metadata.esf"

func NewFileStorage() *FileStorage {
	fs := &FileStorage{storagePath: getStoragePath(), cache: new(LRUPartitionCache)}
	fs.initStoragePath()
	fs.loadStorageMetadata()
	fs.loadDataFiles()
	return fs
}

func (fs *FileStorage) initStoragePath() {
	_, err := os.Stat(fs.storagePath)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(fs.storagePath, 0644)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}

func (fs *FileStorage) loadStorageMetadata() {
	// read file metadata
	metadataFile, err := NewEntrySequenceFile(
		filepath.Join(fs.storagePath, storageMetadataFileName), ReadMode)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := metadataFile.Close(); err != nil {
			log.Println(err)
		}

		if err := recover(); err != io.EOF {
			log.Fatal("failed to load storage metadata", err)
		}
	}()
	var readEntry = func() []byte {
		entry, err := metadataFile.ReadEntry()
		if err != nil {
			panic(err)
		}
		return entry
	}

	var files = make(map[string]*File)
	fs.files = files
	for {
		file := &File{}
		file.fs = fs
		file.Id = string(readEntry())
		file.Name = string(readEntry())
		file.Size = utils.ConvertByteToUint32(readEntry(), 0)
		file.ContentType = string(readEntry())
		file.CreatedAt = utils.ConvertByteToInt64(readEntry(), 0)
		partitionNum := utils.ConvertByteToUint16(readEntry(), 0)
		partitions := make([]PartitionId, partitionNum)
		for i := uint16(0); i < partitionNum; i++ {
			partitions[i] = PartitionId(utils.ConvertByteToUint32(readEntry(), 0))
		}
		files[file.Id] = file
	}
}

func (fs *FileStorage) saveStorageMetadata() {
	metadataFile, err := NewEntrySequenceFile(
		filepath.Join(fs.storagePath, storageMetadataFileName), WriteMode)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := metadataFile.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	var writeEntry = func(data []byte) {
		if err := metadataFile.WriteEntry(data); err != nil {
			log.Fatal("failed to save storage metadata", err)
		}
	}

	var buf []byte
	for _, file := range fs.files {
		writeEntry([]byte(file.Id))
		writeEntry([]byte(file.Name))

		buf = make([]byte, 4)
		utils.ConvertUint32ToByte(file.Size, buf, 0)
		writeEntry(buf)

		writeEntry([]byte(file.ContentType))

		buf = make([]byte, 8)
		utils.ConvertInt64ToByte(file.CreatedAt, buf, 0)
		writeEntry(buf)

		buf = make([]byte, 2)
		utils.ConvertUint16ToByte(uint16(len(file.Partitions)), buf, 0)
		writeEntry(buf)

		buf = make([]byte, 4)
		for _, id := range file.Partitions {
			utils.ConvertUint32ToByte(uint32(id), buf, 0)
			writeEntry(buf)
		}
	}
}

func (fs *FileStorage) loadDataFiles() {
	dataFilePath := filepath.Join(fs.storagePath, dataFilesDirName)
	// create directory if not exists
	_, err := os.Stat(dataFilePath)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dataFilePath, 0644)
		if errDir != nil {
			log.Fatal(err)
		}
	}
	// scan storage path and open all sequential files
	files, err := ioutil.ReadDir(dataFilePath)
	if err != nil {
		log.Fatal(err)
	}

	dataFiles := make([]*SequentialFile, len(files))
	for _, fileInfo := range files {
		id, err := strconv.ParseInt(fileInfo.Name(), 10, 16)
		if err != nil {
			log.Fatal(err)
		}

		dataFile, err := NewSequentialFile(filepath.Join(dataFilePath, fileInfo.Name()), 0, 0)
		dataFiles[id] = dataFile
	}

	fs.dataFiles = dataFiles
}

func getStoragePath() string {
	return filepath.Join(getUserHomeDir(), storageDirName)
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

func (fs *FileStorage) GetFile(id string) (f *File, ok bool) {
	f, ok = fs.files[id]
	return
}

func (fs *FileStorage) GetAllFiles() []*File {
	list := make([]*File, len(fs.files))
	i := 0
	for _, file := range fs.files {
		list[i] = file
		i++
	}
	return list
}

//
// DuplicateFileNameError
// write chunk error
// PartitionNumLimitError
// read input error
func (fs *FileStorage) SaveFile(fileName string, contentType string, reader io.Reader) (*File, error) {
	if _, ok := fs.files[fileName]; ok {
		return nil, DuplicateFileNameError
	}

	uid, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}

	file := &File{
		fs:          fs,
		Id:          strings.Replace(uid.String(), "-", "", -1),
		Name:        fileName,
		Size:        0,
		ContentType: contentType,
		CreatedAt:   time.Now().UnixNano(),
		Partitions:  make(Partitions, 0),
	}

	chunkBuf := make([]byte, MaxFileChunkDataSize)
	dataFile, fileId := fs.getAvailableDataFile()
	var fileSize uint32
	for {
		// read input
		n, err := reader.Read(chunkBuf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if len(file.Partitions) >= maxPartitionNum {
			return nil, PartitionNumLimitError
		}
		fileSize += uint32(n)
		// append input to data file
		chunkId, err := dataFile.AppendChunk(chunkBuf)
		if err != nil {
			return nil, err
		}
		// maintain partition info
		partitionId := createPartitionId(fileId, chunkId)
		file.Partitions = append(file.Partitions, partitionId)
	}
	file.Size = fileSize

	fs.files[file.Id] = file
	return file, nil
}

func (fs *FileStorage) getAvailableDataFile() (*SequentialFile, uint16) {
	sz := len(fs.dataFiles)
	if sz == 0 {
		return fs.createDataFile(), 0
	}

	file := fs.dataFiles[sz-1]
	if !file.IsWritable() {
		return fs.createDataFile(), uint16(sz)
	}

	return file, uint16(sz - 1)
}

func (fs *FileStorage) createDataFile() *SequentialFile {
	fileId := len(fs.dataFiles)
	file, err := NewSequentialFile(filepath.Join(fs.storagePath, dataFilesDirName, strconv.Itoa(fileId)),
		MaxFileChunkDataSize, MaxFileChunkNum)
	if err != nil {
		log.Fatal("failed to create new data file", err)
	}
	fs.dataFiles = append(fs.dataFiles, file)
	return file
}

func (fs *FileStorage) DeleteFile(id string) bool {
	file, ok := fs.files[id]
	if ok {
		delete(fs.files, id)
		for _, partitionId := range file.Partitions {
			// mark all Partitions deleted
			if err := fs.deleteChunk(partitionId); err != nil {
				log.Println("failed to delete file chunk, partition id = ", partitionId, err)
				ok = false
			}
		}
	}

	return ok
}

func (fs *FileStorage) getChunk(id PartitionId) (*FileChunk, error) {
	fileId, chunkId := id.split()
	if int(fileId) >= len(fs.dataFiles) {
		return nil, InvalidPartitionIdError
	}
	file := fs.dataFiles[fileId]
	return file.ReadChunk(chunkId)
}

func (fs *FileStorage) deleteChunk(id PartitionId) error {
	fileId, chunkId := id.split()
	if int(fileId) >= len(fs.dataFiles) {
		return InvalidPartitionIdError
	}
	file := fs.dataFiles[fileId]
	return file.DeleteChunk(chunkId)
}

func (fs *FileStorage) Destroy() {
	fs.saveStorageMetadata()
	fs.files = nil
	for _, file := range fs.dataFiles {
		if err := file.Close(); err != nil {
			log.Println("failed to close sequential file handle, path = ", file.path, ", err = ", err)
		}
	}
	fs.cache.Destroy()
}

type File struct {
	fs          *FileStorage
	Id          string // 32
	Name        string // 128
	Size        uint32 // 8
	ContentType string // 32
	CreatedAt   int64  // 8
	Partitions  Partitions
}

// PartitionId = sequential file id + file chunk id
type PartitionId uint32
type Partitions []PartitionId

func (f *File) OpenStream() io.Reader {
	return newFileDataReader(f.fs, f)
}

type FileDataReader struct {
	fs               *FileStorage
	file             *File
	nextPartitionIdx int
	currentChunk     *FileChunk
	chunkReadOffset  int
}

func newFileDataReader(fs *FileStorage, file *File) *FileDataReader {
	return &FileDataReader{
		fs:               fs,
		file:             file,
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

	nRead := utils.Min(availableBytes, len(p))
	n = copy(p, chunk.content[r.chunkReadOffset:r.chunkReadOffset+nRead])
	if n != nRead {
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
		if r.nextPartitionIdx >= len(r.file.Partitions) {
			err = io.EOF
		} else {
			r.currentChunk, err = r.fs.getChunk(r.file.Partitions[r.nextPartitionIdx])
			r.chunkReadOffset = 0
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
