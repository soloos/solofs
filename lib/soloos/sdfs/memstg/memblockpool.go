package memstg

import (
	"math"
	"soloos/sdbone/offheap"
	"soloos/sdfs/types"
)

type MemBlockPool struct {
	options MemBlockPoolOptions
	driver  *MemBlockDriver

	ichunkSize       int
	chunkSize        uintptr
	tmpMemBlockTable *offheap.HKVTableWithBytes12
	memBlockTable    *offheap.HKVTableWithBytes12
}

func (p *MemBlockPool) Init(
	options MemBlockPoolOptions,
	driver *MemBlockDriver,
) error {
	var err error

	p.options = options
	p.driver = driver

	chunkSize := p.options.ChunkSize
	chunksLimit := p.options.ChunksLimit

	chunkPoolChunksLimit := int32(math.Ceil(float64(chunksLimit) * 0.9))
	p.ichunkSize = chunkSize
	p.chunkSize = uintptr(chunkSize)

	p.memBlockTable, err =
		p.driver.offheapDriver.CreateHKVTableWithBytes12("MemBlock",
			int(types.MemBlockStructSize+p.chunkSize), chunkPoolChunksLimit, types.DefaultKVTableSharedCount,
			p.hkvTableInvokePrepareNewBlock,
			p.hkvTableInvokeBeforeReleaseBlock,
		)
	if err != nil {
		return err
	}

	tmpChunkPoolChunksLimit := chunksLimit - chunkPoolChunksLimit
	if tmpChunkPoolChunksLimit == 0 {
		tmpChunkPoolChunksLimit = 1
	}
	p.tmpMemBlockTable, err =
		p.driver.offheapDriver.CreateHKVTableWithBytes12("TmpMemBlock",
			int(types.MemBlockStructSize+p.chunkSize), tmpChunkPoolChunksLimit, types.DefaultKVTableSharedCount,
			p.hkvTableInvokePrepareNewBlock,
			p.hkvTableInvokeBeforeReleaseTmpBlock,
		)
	if err != nil {
		return err
	}

	return nil
}

func (p *MemBlockPool) hkvTableInvokePrepareNewBlock(uMemBlock uintptr) {
	pMemBlock := types.MemBlockUintptr(uMemBlock).Ptr()
	pMemBlock.Reset()
	pMemBlock.Bytes.Data = uMemBlock + types.MemBlockStructSize
	pMemBlock.Bytes.Len = p.ichunkSize
	pMemBlock.Bytes.Cap = pMemBlock.Bytes.Len
}
