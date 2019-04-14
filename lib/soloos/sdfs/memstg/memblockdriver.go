package memstg

import (
	"soloos/sdbone/offheap"
	"soloos/sdfs/types"
)

type MemBlockDriver struct {
	offheapDriver *offheap.OffheapDriver
	tables        map[int]*MemBlockTable
}

func (p *MemBlockDriver) Init(
	offheapDriver *offheap.OffheapDriver,
	memBlockDriverOptions MemBlockDriverOptions,
) error {
	p.offheapDriver = offheapDriver
	p.tables = make(map[int]*MemBlockTable)
	for _, memBlockTableOptions := range memBlockDriverOptions.MemBlockTableOptionsList {
		p.PrepareMemBlockTable(memBlockTableOptions)
	}
	return nil
}

func (p *MemBlockDriver) PrepareMemBlockTable(memBlockTableOptions MemBlockTableOptions) error {
	if _, exists := p.tables[memBlockTableOptions.ObjectSize]; exists {
		return nil
	}

	var (
		memBlockTable *MemBlockTable
		err           error
	)

	memBlockTable = new(MemBlockTable)
	err = memBlockTable.Init(memBlockTableOptions, p)
	if err != nil {
		return err
	}

	p.tables[memBlockTable.options.ObjectSize] = memBlockTable
	return nil
}

// MustGetMemBlockWithReadAcquire get or init a memblock's offheap
func (p *MemBlockDriver) MustGetMemBlockWithReadAcquire(uNetINode types.NetINodeUintptr,
	memBlockIndex int32) (types.MemBlockUintptr, bool) {
	var memBlockID types.PtrBindIndex
	types.EncodePtrBindIndex(&memBlockID, uintptr(uNetINode), memBlockIndex)
	return p.tables[uNetINode.Ptr().MemBlockCap].MustGetMemBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockDriver) TryGetMemBlockWithReadAcquire(uNetINode types.NetINodeUintptr,
	memBlockIndex int32) types.MemBlockUintptr {
	var memBlockID types.PtrBindIndex
	types.EncodePtrBindIndex(&memBlockID, uintptr(uNetINode), memBlockIndex)
	return p.tables[uNetINode.Ptr().MemBlockCap].TryGetMemBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockDriver) ReleaseMemBlockWithReadRelease(uMemBlock types.MemBlockUintptr) {
	if uMemBlock != 0 {
		p.tables[uMemBlock.Ptr().Bytes.Cap].ReleaseMemBlockWithReadRelease(uMemBlock)
	}
}

func (p *MemBlockDriver) MustGetTmpMemBlockWithReadAcquire(uNetINode types.NetINodeUintptr, memBlockID types.PtrBindIndex) types.MemBlockUintptr {
	return p.tables[uNetINode.Ptr().MemBlockCap].MustGetTmpMemBlockWithReadAcquire(memBlockID)
}

func (p *MemBlockDriver) ReleaseTmpMemBlock(uMemBlock types.MemBlockUintptr) {
	p.tables[uMemBlock.Ptr().Bytes.Cap].ReleaseTmpMemBlock(uMemBlock)
}
