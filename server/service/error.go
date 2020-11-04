package service

import "errors"

var (
	EndOfPartitionStreamError = errors.New("all partitions have been consumed")
	DataOutOfChunkError       = errors.New("not enough space to write chunk data")
	InvalidOperationError     = errors.New("invalid operation")
	EntryTooLargeError        = errors.New("entry is too large")
)
