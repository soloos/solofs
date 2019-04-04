package types

import (
	"soloos/sdbone/offheap"
)

type MockMemBlockPool struct {
	offheapDriver *offheap.OffheapDriver
	ichunkSize    int
	hkvTable      *offheap.HKVTableWithBytes12
}

func (p *MockMemBlockPool) Init(offheapDriver *offheap.OffheapDriver, ichunkSize int) error {
	var err error
	p.offheapDriver = offheapDriver
	p.ichunkSize = ichunkSize
	p.hkvTable, err = p.offheapDriver.CreateHKVTableWithBytes12("mock",
		int(MemBlockStructSize+uintptr(p.ichunkSize)), -1, DefaultKVTableSharedCount,
		p.ChunkPoolInvokePrepareNewChunk, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *MockMemBlockPool) ChunkPoolInvokePrepareNewChunk(uObject uintptr) {
	uMemBlock := MemBlockUintptr(uObject)
	uMemBlock.Ptr().Reset()
	uMemBlock.Ptr().Bytes.Data = uObject + MemBlockStructSize
	uMemBlock.Ptr().Bytes.Len = p.ichunkSize
	uMemBlock.Ptr().Bytes.Cap = uMemBlock.Ptr().Bytes.Len
}

func (p *MockMemBlockPool) AllocMemBlock() MemBlockUintptr {
	var memBlockID PtrBindIndex
	uObject, _ := p.hkvTable.MustGetObjectWithReadAcquire(memBlockID)
	uChunk := offheap.ChunkUintptr(uObject)
	uMemBlock := (MemBlockUintptr)(uChunk.Ptr().Data)
	return uMemBlock
}
