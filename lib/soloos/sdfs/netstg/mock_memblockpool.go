package netstg

import (
	"soloos/sdfs/types"
	"soloos/util/offheap"
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
		types.MemBlockStructSize)
	err = p.offheapDriver.InitChunkPool(&p.chunkPool, chunkSize, -1, p.ChunkPoolInvokePrepareNewChunk, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *MockMemBlockPool) ChunkPoolInvokePrepareNewChunk(uChunk offheap.ChunkUintptr) {
	uMemBlock := types.MemBlockUintptr(uChunk.Ptr().Data)
	uMemBlock.Ptr().Reset()
	uMemBlock.Ptr().Chunk = uChunk
	uMemBlock.Ptr().Bytes.Data = uChunk.Ptr().Data + types.MemBlockStructSize
	uMemBlock.Ptr().Bytes.Len = p.ichunkSize
	uMemBlock.Ptr().Bytes.Cap = uMemBlock.Ptr().Bytes.Len
}

func (p *MockMemBlockPool) AllocMemBlock() types.MemBlockUintptr {
	uChunk := p.chunkPool.AllocChunk()
	uMemBlock := (types.MemBlockUintptr)(uChunk.Ptr().Data)
	return uMemBlock
}
