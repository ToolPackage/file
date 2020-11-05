package storage

import (
	"bytes"
	"github.com/go-playground/assert/v2"
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
	assert.Equal(t, err, nil)
	assert.Equal(t, file.Size, uint32(len(data)))
	// read
	file, ok := fs.GetFile("testFile")
	assert.Equal(t, ok, true)
	assert.Equal(t, file.Name, "testFile")
	assert.Equal(t, file.Size, uint32(dataSize))
	assert.Equal(t, file.ContentType, "html/txt")
	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(file.OpenStream())
	assert.Equal(t, err, nil)
	assert.Equal(t, n, int64(dataSize))
	assert.Equal(t, buf.Bytes(), data)
	// delete
	err = fs.DeleteFile("testFile")
	assert.Equal(t, err, nil)
}
