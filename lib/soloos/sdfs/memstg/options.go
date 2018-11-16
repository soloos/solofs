package memstg

type INodePoolOptions struct {
	RawChunksLimit int32
}

type MemBlockDriverOptions struct {
	MemBlockPoolOptionsList []MemBlockPoolOptions
}

type MemBlockPoolOptions struct {
	ChunkSize   int
	ChunksLimit int32
}
