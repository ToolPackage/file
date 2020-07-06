package service

import (
	"fmt"
	constants "github.com/ToolPackage/fse/server/common"
	"os"
	"testing"
)

const TestFilePath = "../../tmp/sequential_file.tmp"

func BenchmarkSequentialFile_Append(b *testing.B) {
	f, err := NewSequentialFile(TestFilePath, constants.MaxSequentialFileSize, 0)
	if err != nil {
		panic(err)
	}

	// prepare test data
	data := make([]byte, constants.FileChunkSize)
	data[512] = 'A'

	b.StartTimer()
	var n int
	for i := 0; i < b.N; i++ {
		n, err = f.Append(data)
		if n != len(data) {
			panic(fmt.Sprintf("Expected: %d, got: %d", len(data), n))
		}
	}
	b.StopTimer()

	f.Close()
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
