package memstg

import (
	"fmt"
	"math"
	"soloos/sdbone/offheap"
	"soloos/sdfs/types"
)

type memBlockPoolHKVTable struct {
	memBlockPool     *MemBlockPool
	chunkSize        uintptr
	tmpMemBlockTable *offheap.HKVTable
	memBlockTable    *offheap.HKVTable
}

func (p *memBlockPoolHKVTable) Init(memBlockPool *MemBlockPool, chunkSize int, chunksLimit int32) error {
	var err error

	p.memBlockPool = memBlockPool

	chunkPoolChunksLimit := int32(math.Ceil(float64(chunksLimit) * 0.9))
	p.chunkSize = uintptr(chunkSize)

	hkvTableKeyType := fmt.Sprintf("Bytes%d", types.PtrBindIndexSize)

	p.memBlockTable, err =
		p.memBlockPool.driver.offheapDriver.CreateHKVTable("MemBlock",
			int(types.MemBlockStructSize+p.chunkSize),
			chunkPoolChunksLimit,
			hkvTableKeyType, 32,
			p.beforeReleaseBlock)
	if err != nil {
		return err
	}

	tmpChunkPoolChunksLimit := chunksLimit - chunkPoolChunksLimit
	if tmpChunkPoolChunksLimit == 0 {
		tmpChunkPoolChunksLimit = 1
	}
	p.tmpMemBlockTable, err =
		p.memBlockPool.driver.offheapDriver.CreateHKVTable("TmpMemBlock",
			int(types.MemBlockStructSize+p.chunkSize),
			tmpChunkPoolChunksLimit,
			hkvTableKeyType, 32,
			p.beforeReleaseTmpBlock)
	if err != nil {
		return err
	}

	return nil
}

func (p *memBlockPoolHKVTable) beforeReleaseBlock(uMemBlock uintptr) {
	pMemBlock := types.MemBlockUintptr(uMemBlock).Ptr()
	if pMemBlock.IsInited() == false {
		return
	}
	pMemBlock.UploadJob.SyncDataSig.Wait()
	pMemBlock.SetReleasable()
}

func (p *memBlockPoolHKVTable) MustGetMemBlockWithReadAcquire(blockID types.PtrBindIndex) (types.MemBlockUintptr, bool) {
	var (
		uObject types.MemBlockUintptr
		u       uintptr
		loaded  bool
	)
	u, loaded = p.memBlockTable.MustGetObjectByBytes12WithReadAcquire(blockID)
	uObject = types.MemBlockUintptr(u)
	return uObject, loaded
}

func (p *memBlockPoolHKVTable) TryGetMemBlockWithReadAcquire(blockID types.PtrBindIndex) types.MemBlockUintptr {
	var uObject types.MemBlockUintptr
	uObject = types.MemBlockUintptr(p.memBlockTable.TryGetObjectByBytes12WithReadAcquire(blockID))
	return uObject
}
