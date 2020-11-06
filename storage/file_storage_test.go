package storage

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFileStorage(t *testing.T) {
	fs := NewFileStorage()
	defer fs.Destroy()

	var dataSize = MaxFileChunkDataSize * 3
	data := make([]byte, dataSize)
	data[0] = 'A'
	data[MaxFileChunkDataSize] = 'B'
	data[MaxFileChunkDataSize*2] = 'C'
	// write
	file, err := fs.SaveFile("testFile", "html/txt", bytes.NewReader(data))
	assert.Nil(t, err)
	id := file.Id
	assert.True(t, file.Size == uint32(len(data)))
	// read
	file, ok := fs.GetFile(id)
	assert.True(t, ok)
	assert.Equal(t, file.Id, id)
	assert.Equal(t, file.Name, "testFile")
	assert.True(t, file.Size == uint32(dataSize))
	assert.Equal(t, file.ContentType, "html/txt")
	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(file.OpenStream())
	assert.Nil(t, err)
	assert.True(t, n == int64(dataSize))
	assert.Equal(t, buf.Bytes(), data)
	// delete
	ok = fs.DeleteFile(id)
	assert.True(t, ok)
}
