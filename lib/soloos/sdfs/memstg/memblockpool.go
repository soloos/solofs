package memstg

import (
	"soloos/sdfs/types"
)

type MemBlockPool struct {
	options MemBlockPoolOptions
	driver  *MemBlockDriver
	chunk   memBlockPoolChunk
}

func (p *MemBlockPool) Init(options MemBlockPoolOptions, driver *MemBlockDriver) error {
	var err error

	p.options = options
	p.driver = driver

	err = p.chunk.Init(p, p.options.ChunkSize, p.options.ChunksLimit)
	if err != nil {
		return err
	}

	return nil
}

// MustGetBlockWithReadAcquire get or init a inodeblock
func (p *MemBlockPool) MustGetBlockWithReadAcquire(memBlockID types.PtrBindIndex) (types.MemBlockUintptr, bool) {
	return p.chunk.MustGetBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockPool) AllocTmpBlockWithWriteAcquire() types.MemBlockUintptr {
	return p.chunk.AllocTmpBlockWithWriteAcquire()
}

func (p *MemBlockPool) ReleaseTmpBlock(uMemBlock types.MemBlockUintptr) {
	p.chunk.ReleaseTmpBlock(uMemBlock)
}
