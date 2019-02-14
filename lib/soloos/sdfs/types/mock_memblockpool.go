package types

import (
	"soloos/sdbone/offheap"
)

type MockMemBlockPool struct {
	offheapDriver *offheap.OffheapDriver
	ichunkSize    int
	chunkPool     offheap.ChunkPool
}

func (p *MockMemBlockPool) Init(offheapDriver *offheap.OffheapDriver, ichunkSize int) error {
	var err error
	p.offheapDriver = offheapDriver
	p.ichunkSize = ichunkSize
	chunkSize := int(uintptr(p.ichunkSize) +
		MemBlockStructSize)
	err = p.offheapDriver.InitChunkPool(&p.chunkPool, chunkSize, -1, p.ChunkPoolInvokePrepareNewChunk, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *MockMemBlockPool) ChunkPoolInvokePrepareNewChunk(uChunk offheap.ChunkUintptr) {
	uMemBlock := MemBlockUintptr(uChunk.Ptr().Data)
	uMemBlock.Ptr().Reset()
	uMemBlock.Ptr().Chunk = uChunk
	uMemBlock.Ptr().Bytes.Data = uChunk.Ptr().Data + MemBlockStructSize
	uMemBlock.Ptr().Bytes.Len = p.ichunkSize
	uMemBlock.Ptr().Bytes.Cap = uMemBlock.Ptr().Bytes.Len
}

func (p *MockMemBlockPool) AllocMemBlock() MemBlockUintptr {
	uChunk := p.chunkPool.AllocChunk()
	uMemBlock := (MemBlockUintptr)(uChunk.Ptr().Data)
	return uMemBlock
}
