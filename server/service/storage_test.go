package service

import (
	"os"
	"testing"
)

const TestFilePath = "../tmp/sequential_file.tmp"
const TestFileSize = 5 * 1024 * 1024 * 1024

func Before() {
	fileInfo, err := os.Stat(TestFilePath)
	if err != nil {
		panic(err)
	}
	if fileInfo.Size() != TestFileSize {
		// create file with specified size
		file, err = os.OpenFile(TestFilePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}

		// TODO: mmap
	}
}

func BenchmarkSequentialFile_Read(b *testing.B) {
	f, err := NewSequentialFile(TestFilePath)
	if err != nil {
		panic(err)
	}

}
