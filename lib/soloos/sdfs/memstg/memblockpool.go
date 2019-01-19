package memstg

import (
	"soloos/sdfs/types"
)

type MemBlockPool struct {
	options MemBlockPoolOptions
	driver  *MemBlockDriver
	chunk   memBlockPoolChunk
}

func (p *MemBlockPool) Init(
	options MemBlockPoolOptions,
	driver *MemBlockDriver,
) error {
	var err error

	p.options = options
	p.driver = driver

	err = p.chunk.Init(p, p.options.ChunkSize, p.options.ChunksLimit)
	if err != nil {
		return err
	}

	return nil
}

// MustGetMemBlockWithReadAcquire get or init a netINodeblock
func (p *MemBlockPool) MustGetMemBlockWithReadAcquire(memBlockID types.PtrBindIndex) (types.MemBlockUintptr, bool) {
	return p.chunk.MustGetMemBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockPool) TryGetMemBlockWithReadAcquire(memBlockID types.PtrBindIndex) types.MemBlockUintptr {
	return p.chunk.TryGetMemBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockPool) ReleaseMemBlockWithReadRelease(uMemBlock types.MemBlockUintptr) {
	uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
}

func (p *MemBlockPool) AllocTmpMemBlockWithWriteAcquire() types.MemBlockUintptr {
	return p.chunk.AllocTmpMemBlockWithWriteAcquire()
}

func (p *MemBlockPool) ReleaseTmpMemBlock(uMemBlock types.MemBlockUintptr) {
	p.chunk.ReleaseTmpMemBlock(uMemBlock)
}
