package storage

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

const sequentialFilePath = "../tmp/sequential_file.tmp"

func TestNewSequentialFile(t *testing.T) {
	var (
		chunkNum = 100
		stubPos  = 512
	)
	f, err := NewSequentialFile(sequentialFilePath, MaxFileChunkDataSize, MaxFileChunkNum)
	assert.Nil(t, err)

	// prepare test Partitions
	data := make([]byte, MaxFileChunkDataSize)
	data[stubPos] = 'A'

	var chunkId uint16
	for i := 0; i < chunkNum; i++ {
		chunkId, err = f.AppendChunk(data)
		assert.Nil(t, err)
		assert.True(t, chunkId == uint16(i))
	}

	err = f.Close()
	assert.Nil(t, err)

	f, err = NewSequentialFile(sequentialFilePath, 0, 0)
	assert.Nil(t, err)

	var chunk *FileChunk
	for i := 0; i < chunkNum; i++ {
		chunk, err = f.ReadChunk(uint16(i))
		assert.Nil(t, err)
		assert.True(t, chunk.content[stubPos] == 'A')
		assert.True(t, chunk.Validate())
	}

	err = f.Close()
	assert.Nil(t, err)
	err = f.Delete()
	assert.Nil(t, err)
}

func TestSequentialFile_AppendChunk(t *testing.T) {
	f, err := NewSequentialFile(sequentialFilePath, 20, 5)
	assert.Nil(t, err)

	data := make([]byte, 13)
	data[10] = 'A'

	chunkId, err := f.AppendChunk(data)
	assert.Nil(t, err)
	assert.True(t, chunkId == 0)

	err = f.Close()
	assert.Nil(t, err)

	f, err = NewSequentialFile(sequentialFilePath, 0, 0)
	assert.Nil(t, err)

	chunk, err := f.ReadChunk(0)
	assert.Nil(t, err)

	assert.True(t, chunk.chunkId == 0)
	assert.True(t, bytes.Equal(chunk.content, data))
	assert.True(t, chunk.Validate())

	err = f.Close()
	assert.Nil(t, err)
	err = f.Delete()
	assert.Nil(t, err)
}
