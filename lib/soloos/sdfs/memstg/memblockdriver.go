package memstg

import (
	"soloos/sdfs/types"
	"soloos/util/offheap"
)

type MemBlockDriver struct {
	offheapDriver *offheap.OffheapDriver
	pools         map[int]*MemBlockPool
}

func (p *MemBlockDriver) Init(
	offheapDriver *offheap.OffheapDriver,
	options MemBlockDriverOptions,
) error {
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

// MustGetMemBlockWithReadAcquire get or init a memblock's offheap
func (p *MemBlockDriver) MustGetMemBlockWithReadAcquire(uNetINode types.NetINodeUintptr,
	memBlockIndex int) (types.MemBlockUintptr, bool) {
	var memBlockID types.PtrBindIndex
	types.EncodePtrBindIndex(&memBlockID, uintptr(uNetINode), memBlockIndex)
	return p.pools[uNetINode.Ptr().MemBlockCap].MustGetMemBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockDriver) TryGetBlockWithReadAcquire(uNetINode types.NetINodeUintptr,
	memBlockIndex int) types.MemBlockUintptr {
	var memBlockID types.PtrBindIndex
	types.EncodePtrBindIndex(&memBlockID, uintptr(uNetINode), memBlockIndex)
	return p.pools[uNetINode.Ptr().MemBlockCap].TryGetBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockDriver) AllocTmpBlockWithWriteAcquire(uNetINode types.NetINodeUintptr) types.MemBlockUintptr {
	return p.pools[uNetINode.Ptr().MemBlockCap].AllocTmpBlockWithWriteAcquire()
}

func (p *MemBlockDriver) ReleaseTmpBlock(uNetINode types.NetINodeUintptr, uMemBlock types.MemBlockUintptr) {
	p.pools[uNetINode.Ptr().MemBlockCap].ReleaseTmpBlock(uMemBlock)
}
