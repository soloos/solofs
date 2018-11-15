package memstg

import (
	"soloos/util/offheap"
	"soloos/sdfs/types"
)

type MemBlockDriver struct {
	offheapDriver *offheap.OffheapDriver
	pools         map[int]*MemBlockPool
}

func (p *MemBlockDriver) Init(options MemBlockDriverOptions,
	offheapDriver *offheap.OffheapDriver) error {
	var err error

	p.offheapDriver = offheapDriver

	var inodeBlockPool *MemBlockPool
	p.pools = make(map[int]*MemBlockPool)
	for _, inodeBlockPoolOptions := range options.MemBlockPoolOptionsList {
		inodeBlockPool = new(MemBlockPool)
		err = inodeBlockPool.Init(inodeBlockPoolOptions, p)
		if err != nil {
			return err
		}

		p.pools[inodeBlockPool.options.ChunkPoolOptions.ChunkSize] = inodeBlockPool
	}

	return nil
}

// MustGetBlockWithReadAcquire get or init a inodeblock's offheap
func (p *MemBlockDriver) MustGetBlockWithReadAcquire(uINode types.INodeUintptr,
	memBlockIndex int) (types.MemBlockUintptr, bool) {
	var memBlockID types.PtrBindIndex
	types.EncodePtrBindIndex(&memBlockID, uintptr(uINode), memBlockIndex)
	return p.pools[uINode.Ptr().MemBlockSize].MustGetBlockWithReadAcquire(memBlockID)
}
