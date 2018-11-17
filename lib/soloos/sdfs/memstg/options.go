package memstg

type MemBlockDriverOptions struct {
	MemBlockPoolOptionsList []MemBlockPoolOptions
}

type MemBlockPoolOptions struct {
	ChunkSize   int
	ChunksLimit int32
}
