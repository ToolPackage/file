package storage

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const entrySequenceFilePath = "../../tmp/entry_sequential_file.tmp"

func TestNewEntrySequenceFile(t *testing.T) {
	f := NewEntrySequenceFile(entrySequenceFilePath, WriteMode)

	testData := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(maxEntrySize)
		entry := randomString(n)
		testData[i] = entry
		f.WriteEntry([]byte(entry))
	}

	f.Close()

	f = NewEntrySequenceFile(entrySequenceFilePath, ReadMode)

	for i := 0; i < 100; i++ {
		entry := f.ReadEntry()
		if string(entry) != testData[i] {
			panic(fmt.Sprintf("Expected: %s, got: %s", testData[i], string(entry)))
		}
	}

	if entry := f.ReadEntry(); entry != nil {
		panic(fmt.Sprintf("Expect to read nil"))
	}

	f.Close()
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
