package storage

import (
	"bytes"
	"github.com/go-playground/assert/v2"
	"testing"
)

const sequentialFilePath = "../../tmp/sequential_file.tmp"

func TestNewSequentialFile(t *testing.T) {
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
	err = f.Delete()
	assert.Equal(t, err, nil)
}

func TestSequentialFile_AppendChunk(t *testing.T) {
	f, err := NewSequentialFile(sequentialFilePath, 20, 5)
	assert.Equal(t, err, nil)

	data := make([]byte, 13)
	data[10] = 'A'

	chunkId, err := f.AppendChunk(data)
	assert.Equal(t, err, nil)
	assert.Equal(t, chunkId == 0, true)

	err = f.Close()
	assert.Equal(t, err, nil)

	f, err = NewSequentialFile(sequentialFilePath, 0, 0)
	assert.Equal(t, err, nil)

	chunk, err := f.ReadChunk(0)
	assert.Equal(t, err, nil)

	assert.Equal(t, chunk.chunkId == 0, true)
	assert.Equal(t, bytes.Equal(chunk.content, data), true)
	assert.Equal(t, chunk.Validate(), true)

	err = f.Close()
	assert.Equal(t, err, nil)
	err = f.Delete()
	assert.Equal(t, err, nil)
}
