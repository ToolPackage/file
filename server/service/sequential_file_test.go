package service

import (
	"fmt"
	"os"
	"testing"
)

const TestFilePath = "../../tmp/sequential_file.tmp"

func BenchmarkSequentialFile_Append(b *testing.B) {
	f, err := NewSequentialFile(TestFilePath, MaxFileChunkDataSize, MaxFileChunkNum)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// prepare test partitions
	data := make([]byte, MaxFileChunkDataSize)
	data[512] = 'A'

	b.StartTimer()
	var chunkId uint16
	for i := 0; i < b.N; i++ {
		chunkId, err = f.AppendChunk(data)
		if int(chunkId) != i {
			panic(fmt.Sprintf("Expected: %d, got: %d", i, chunkId))
		}
	}
	b.StopTimer()
}

func setup() {

}

func shutdown() {
	// delete test file
	if err := os.Remove(TestFilePath); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
