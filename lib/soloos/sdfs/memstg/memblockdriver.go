package memstg

import (
	"soloos/sdfs/types"
	"soloos/util/offheap"
)

type MemBlockDriver struct {
	offheapDriver *offheap.OffheapDriver
	pools         map[int]*MemBlockPool
}

func (p *MemBlockDriver) Init(options MemBlockDriverOptions,
	offheapDriver *offheap.OffheapDriver) error {
	var err error

	p.offheapDriver = offheapDriver

	var memblockPool *MemBlockPool
	p.pools = make(map[int]*MemBlockPool)
	for _, memblockPoolOptions := range options.MemBlockPoolOptionsList {
		memblockPool = new(MemBlockPool)
		err = memblockPool.Init(memblockPoolOptions, p)
		if err != nil {
			return err
		}

		p.pools[memblockPool.options.ChunkSize] = memblockPool
	}

	return nil
}

// MustGetBlockWithReadAcquire get or init a memblock's offheap
func (p *MemBlockDriver) MustGetBlockWithReadAcquire(uINode types.INodeUintptr,
	memBlockIndex int) (types.MemBlockUintptr, bool) {
	var memBlockID types.PtrBindIndex
	types.EncodePtrBindIndex(&memBlockID, uintptr(uINode), memBlockIndex)
	return p.pools[uINode.Ptr().MemBlockSize].MustGetBlockWithReadAcquire(memBlockID)
}
