package service

type FileChunk struct {
	chunkId int
	md5     string
	content []byte
}

type ChunkedFile struct {
	filename    string
	contentType string
	createdAt   int64
	chunks      []FileChunk
}
