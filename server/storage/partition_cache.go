package storage

type PartitionCache interface {
	GetChunk(id PartitionId, onFail func() FileChunk) FileChunk
	Destroy()
}

type LRUPartitionCache struct {
	partitionMap map[PartitionId]FileChunk
}

func (c *LRUPartitionCache) GetChunk(id PartitionId, onFail func() FileChunk) FileChunk {
	// TODO:
	//chunk, ok := c.partitionMap[id]
	//if ok {
	//	return chunk
	//}
	//
	//chunk = onFail()
	//c.partitionMap[id] = chunk
	//return chunk
	return onFail()
}

func (c *LRUPartitionCache) Destroy() {

}
