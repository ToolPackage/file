package service

type PartitionCache interface {
	GetChunk(id PartitionId, onFail func() FileChunk) FileChunk
}

type LRUPartitionCache struct {
	partitionMap map[PartitionId]FileChunk
}

func (c *LRUPartitionCache) GetChunk(id PartitionId, onFail func() FileChunk) FileChunk {
	chunk, ok := c.partitionMap[id]
	if ok {
		return chunk
	}

	chunk = onFail()
	c.partitionMap[id] = chunk
	return chunk
}
