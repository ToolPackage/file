package service

import "errors"

var (
	EndOfPartitionStreamError = errors.New("all partitions have been consumed")
	DataOutOfChunkError       = errors.New("not enough space to write chunk data")
	DataOutOfFileError        = errors.New("not enough space to write data to file")
	InvalidOperationError     = errors.New("invalid operation")
	EntryTooLargeError        = errors.New("entry is too large")
	InvalidPartitionIdError   = errors.New("invalid partition id")
	InvalidChunkIdError       = errors.New("invalid chunk id")
	InvalidRetValue           = errors.New("invalid ret value")
	PartitionNumLimitError    = errors.New("partition num limit error")
)
