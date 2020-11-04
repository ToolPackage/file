package service

import (
	"fmt"
	"io"
	"math/rand"
	"testing"
	"time"
)

const testFilePath = "../../tmp/entry_sequential_file.tmp"

func TestNewEntrySequenceFile(t *testing.T) {
	f, err := NewEntrySequenceFile(testFilePath, WriteMode)
	if err != nil {
		panic(err)
	}

	testData := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(maxEntrySize)
		entry := randomString(n)
		testData[i] = entry
		if err = f.WriteEntry([]byte(entry)); err != nil {
			panic(err)
		}
	}

	f.Close()

	f, err = NewEntrySequenceFile(testFilePath, ReadMode)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		entry, err := f.ReadEntry()
		if err != nil {
			panic(err)
		}

		if string(entry) != testData[i] {
			panic(fmt.Sprintf("Expected: %s, got: %s", testData[i], string(entry)))
		}
	}

	_, err = f.ReadEntry()
	if err != io.EOF {
		panic(fmt.Sprintf("Expected: EOF, got: %s", err))
	}

	f.Close()
}

var seededRand *rand.Rand = rand.New(
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
