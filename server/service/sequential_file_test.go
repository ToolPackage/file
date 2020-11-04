package service

import (
	"fmt"
	"os"
	"testing"
)

const sequentialFilePath = "../../tmp/sequential_file.tmp"

func TestSequentialFile_Append(t *testing.T) {
	var (
		chunkNum = 100
		stubPos  = 512
	)
	f, err := NewSequentialFile(sequentialFilePath, MaxFileChunkDataSize, MaxFileChunkNum)
	if err != nil {
		panic(err)
	}

	// prepare test partitions
	data := make([]byte, MaxFileChunkDataSize)
	data[stubPos] = 'A'

	var chunkId uint16
	for i := 0; i < chunkNum; i++ {
		chunkId, err = f.AppendChunk(data)
		if err != nil {
			panic(err)
		}
		if chunkId != uint16(i) {
			panic(fmt.Sprintf("Expected: %d, got: %d", i, chunkId))
		}
	}

	f.Close()

	f, err = NewSequentialFile(sequentialFilePath, 0, 0)
	if err != nil {
		panic(err)
	}

	var chunk *FileChunk
	for i := 0; i < chunkNum; i++ {
		chunk, err = f.ReadChunk(uint16(i))
		if err != nil {
			panic(err)
		}
		if chunk.content[stubPos] != 'A' {
			panic("stub character check failed")
		}
		if !chunk.Validate() {
			panic("md5 check failed")
		}
	}

	f.Close()
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
