package memstg

import "soloos/util/offheap"

type INodePoolOptions struct {
	RawChunksLimit int32
}

type MemBlockDriverOptions struct {
	MemBlockPoolOptionsList []MemBlockPoolOptions
}

type MemBlockPoolOptions struct {
	ChunksLimit      int32
	ChunkPoolOptions offheap.ChunkPoolOptions
}
