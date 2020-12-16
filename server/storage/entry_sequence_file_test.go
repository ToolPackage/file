package storage

import (
	"github.com/go-playground/assert/v2"
	"io"
	"math/rand"
	"testing"
	"time"
)

const entrySequenceFilePath = "../tmp/entry_sequential_file.tmp"

func TestNewEntrySequenceFile(t *testing.T) {
	f, err := NewEntrySequenceFile(entrySequenceFilePath, WriteMode)
	assert.Equal(t, err, nil)

	testData := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(maxEntrySize)
		entry := randomString(n)
		testData[i] = entry
		err = f.WriteEntry([]byte(entry))
		assert.Equal(t, err, nil)
	}

	err = f.Close()
	assert.Equal(t, err, nil)

	f, err = NewEntrySequenceFile(entrySequenceFilePath, ReadMode)
	assert.Equal(t, err, nil)

	for i := 0; i < 100; i++ {
		entry, err := f.ReadEntry()
		assert.Equal(t, err, nil)
		assert.Equal(t, entry, []byte(testData[i]))
	}

	entry, err := f.ReadEntry()
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, entry, nil)

	err = f.Close()
	assert.Equal(t, err, nil)
	err = f.Delete()
	assert.Equal(t, err, nil)
}

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(length int) string {
	return randomStringWithCharset(length, charset)
}

func randomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
