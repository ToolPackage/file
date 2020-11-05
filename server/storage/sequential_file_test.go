package storage

import (
	"github.com/go-playground/assert/v2"
	"os"
	"testing"
)

const sequentialFilePath = "../../tmp/sequential_file.tmp"

// TODO: md5 bug, reopen

func TestSequentialFile_Append(t *testing.T) {
	var (
		chunkNum = 100
		stubPos  = 512
	)
	f, err := NewSequentialFile(sequentialFilePath, MaxFileChunkDataSize, MaxFileChunkNum)
	assert.Equal(t, err, nil)

	// prepare test partitions
	data := make([]byte, MaxFileChunkDataSize)
	data[stubPos] = 'A'

	var chunkId uint16
	for i := 0; i < chunkNum; i++ {
		chunkId, err = f.AppendChunk(data)
		assert.Equal(t, err, nil)
		assert.Equal(t, chunkId, uint16(i))
	}

	err = f.Close()
	assert.Equal(t, err, nil)

	f, err = NewSequentialFile(sequentialFilePath, 0, 0)
	assert.Equal(t, err, nil)

	var chunk *FileChunk
	for i := 0; i < chunkNum; i++ {
		chunk, err = f.ReadChunk(uint16(i))
		assert.Equal(t, err, nil)
		assert.Equal(t, chunk.content[stubPos] == 'A', true)
		assert.Equal(t, chunk.Validate(), true)
	}

	err = f.Close()
	assert.Equal(t, err, nil)
}

func setup() {

}

func shutdown() {
	deleteTmpFile(sequentialFilePath)
	deleteTmpFile(entrySequenceFilePath)
}

func deleteTmpFile(path string) {
	// delete test file
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return
		}
	}
	if err := os.Remove(path); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
