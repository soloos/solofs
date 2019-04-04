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
			int(types.MemBlockStructSize+p.chunkSize), chunkPoolChunksLimit, 32,
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
			int(types.MemBlockStructSize+p.chunkSize), tmpChunkPoolChunksLimit, 32,
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

func (p *MemBlockPool) hkvTableInvokeBeforeReleaseBlock(uMemBlock uintptr) {
	pMemBlock := types.MemBlockUintptr(uMemBlock).Ptr()
	pMemBlock.UploadJob.SyncDataSig.Wait()
	pMemBlock.SetReleasable()
}

// MustGetMemBlockWithReadAcquire get or init a netINodeblock
func (p *MemBlockPool) MustGetMemBlockWithReadAcquire(memBlockID types.PtrBindIndex) (types.MemBlockUintptr, bool) {
	var (
		uObject types.MemBlockUintptr
		u       uintptr
		loaded  bool
	)
	u, loaded = p.memBlockTable.MustGetObjectWithReadAcquire(memBlockID)
	uObject = types.MemBlockUintptr(u)
	return uObject, loaded
}

func (p *MemBlockPool) TryGetMemBlockWithReadAcquire(memBlockID types.PtrBindIndex) types.MemBlockUintptr {
	var uObject types.MemBlockUintptr
	uObject = types.MemBlockUintptr(p.memBlockTable.TryGetObjectWithReadAcquire(memBlockID))
	return uObject
}

func (p *MemBlockPool) ReleaseMemBlockWithReadRelease(uMemBlock types.MemBlockUintptr) {
	uMemBlock.Ptr().ReadRelease()
}

func (p *MemBlockPool) hkvTableInvokeBeforeReleaseTmpBlock(uMemBlock uintptr) {
	pMemBlock := types.MemBlockUintptr(uMemBlock).Ptr()
	pMemBlock.SetReleasable()
}

func (p *MemBlockPool) MustGetTmpMemBlockWithReadAcquire(memBlockID types.PtrBindIndex) types.MemBlockUintptr {
	var (
		uObject types.MemBlockUintptr
		u       uintptr
	)
	u, _ = p.tmpMemBlockTable.MustGetObjectWithReadAcquire(memBlockID)
	uObject = types.MemBlockUintptr(u)
	return uObject

}

func (p *MemBlockPool) ReleaseTmpMemBlock(uMemBlock types.MemBlockUintptr) {
	p.tmpMemBlockTable.DeleteObject(uMemBlock.Ptr().ID)
}
