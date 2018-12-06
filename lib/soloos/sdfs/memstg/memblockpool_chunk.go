package memstg

import (
	"math"
	"soloos/sdfs/types"
	"soloos/util/offheap"
	"sync"
)

// memBlockPoolChunk
// user -> MustGetBlockWithReadAcquire -> allocChunkFromChunkPool ->
//      offheap.BlockPool.AllocBlock ->
//		BlockPoolAssistant.ChunkPoolInvokeReleaseChunk ->
//      takeBlockForRelease -> beforeReleaseBlock -> releaseChunkToChunkPool ->
//		BlockPoolAssistant.ReleaseBlock -> user
// user -> MustGetBlockWithReadAcquire -> offheap.BlockPool.AllocBlock ->
// 		BlockPoolAssistant.ChunkPoolInvokePrepareNewChunk ->
type memBlockPoolChunk struct {
	memBlockPool *MemBlockPool

	ichunkSize int
	chunkSize  uintptr

	tmpChunkPool        offheap.ChunkPool
	workingTmpChunkPool types.HotPool

	chunkPool        offheap.ChunkPool
	workingChunkPool types.HotPool
	memBlocksRWMutex sync.RWMutex
	memBlocks        map[types.PtrBindIndex]types.MemBlockUintptr
}

func (p *memBlockPoolChunk) Init(memBlockPool *MemBlockPool, chunkSize int, chunksLimit int32) error {
	var err error

	p.memBlockPool = memBlockPool

	chunkPoolChunksLimit := int32(math.Ceil(float64(chunksLimit) * 0.9))
	p.chunkSize = uintptr(chunkSize)
	p.ichunkSize = chunkSize
	err = p.memBlockPool.driver.offheapDriver.InitChunkPool(&p.chunkPool,
		int(p.chunkSize+types.MemBlockStructSize),
		chunkPoolChunksLimit,
		p.ChunkPoolInvokePrepareNewChunk,
		p.ChunkPoolInvokeReleaseChunk)
	if err != nil {
		return err
	}
	p.workingChunkPool.Init()

	tmpChunkPoolChunksLimit := chunksLimit - chunkPoolChunksLimit
	if tmpChunkPoolChunksLimit == 0 {
		tmpChunkPoolChunksLimit = 1
	}
	err = p.memBlockPool.driver.offheapDriver.InitChunkPool(&p.tmpChunkPool,
		int(p.chunkSize+types.MemBlockStructSize),
		tmpChunkPoolChunksLimit,
		p.ChunkPoolInvokePrepareNewChunk,
		p.ChunkPoolInvokeReleaseTmpChunk)
	if err != nil {
		return err
	}

	p.workingTmpChunkPool.Init()

	p.memBlocks = make(map[types.PtrBindIndex]types.MemBlockUintptr)

	return nil
}

// ChunkPoolInvokePrepareNewChunk call by offheap.BlockPool
func (p *memBlockPoolChunk) ChunkPoolInvokePrepareNewChunk(uChunk offheap.ChunkUintptr) {
	uMemBlock := types.MemBlockUintptr(uChunk.Ptr().Data)
	uMemBlock.Ptr().Reset()
	uMemBlock.Ptr().Chunk = uChunk
	uMemBlock.Ptr().Bytes.Data = uChunk.Ptr().Data + types.MemBlockStructSize
	uMemBlock.Ptr().Bytes.Len = p.ichunkSize
	uMemBlock.Ptr().Bytes.Cap = uMemBlock.Ptr().Bytes.Len
}

// ChunkPoolInvokeReleaseChunk call by offheap.BlockPool
func (p *memBlockPoolChunk) ChunkPoolInvokeReleaseChunk() {
	uMemBlock := p.takeBlockForRelease()
	if uMemBlock == 0 {
		return
	}
	p.releaseBlock(uMemBlock)
}

func (p *memBlockPoolChunk) allocChunkFromChunkPool() types.MemBlockUintptr {
	uChunk := p.chunkPool.AllocChunk()
	uMemBlock := (types.MemBlockUintptr)(uChunk.Ptr().Data)
	return uMemBlock
}

func (p *memBlockPoolChunk) releaseChunkToChunkPool(uMemBlock types.MemBlockUintptr) {
	uMemBlock.Ptr().Reset()
	p.chunkPool.ReleaseChunk(uMemBlock.Ptr().Chunk)
}

func (p *memBlockPoolChunk) takeBlockForRelease() types.MemBlockUintptr {
	iRet := p.workingChunkPool.IteratorAndPop(func(x interface{}) (bool, interface{}) {
		uMemBlock := x.(types.MemBlockUintptr)
		pMemBlock := uMemBlock.Ptr()
		if pMemBlock.IsInited() == false && pMemBlock.Chunk.Ptr().Accessor > 0 {
			return false, 0
		}
		return true, uMemBlock
	})

	if iRet != nil {
		return iRet.(types.MemBlockUintptr)
	}

	// Get Block From workingChunkPool
	iRet = p.workingChunkPool.IteratorAndPop(func(x interface{}) (bool, interface{}) {
		uMemBlock := x.(types.MemBlockUintptr)
		pMemBlock := uMemBlock.Ptr()
		if pMemBlock.IsInited() == false {
			return false, 0
		}
		return true, uMemBlock
	})

	if iRet == nil {
		return 0
	}

	return iRet.(types.MemBlockUintptr)
}

func (p *memBlockPoolChunk) beforeReleaseBlock(pMemBlock *types.MemBlock) {
	if pMemBlock.IsInited() == false {
		return
	}
	pMemBlock.UploadJob.UploadSig.Wait()
	pMemBlock.SetReleasable()
}

func (p *memBlockPoolChunk) releaseBlock(uMemBlock types.MemBlockUintptr) {
	pMemBlock := uMemBlock.Ptr()

	pMemBlock.Chunk.Ptr().WriteAcquire()
	p.beforeReleaseBlock(pMemBlock)
	if pMemBlock.EnsureRelease() {
		p.memBlocksRWMutex.Lock()
		delete(p.memBlocks, pMemBlock.ID)
		p.memBlocksRWMutex.Unlock()
		pMemBlock.Chunk.Ptr().WriteRelease()
		p.releaseChunkToChunkPool(uMemBlock)
	} else {
		pMemBlock.Chunk.Ptr().WriteRelease()
	}
}

func (p *memBlockPoolChunk) checkBlock(blockID types.PtrBindIndex, uMemBlock types.MemBlockUintptr) bool {
	if uMemBlock.Ptr().ID != blockID {
		return false
	}
	return true
}

func (p *memBlockPoolChunk) MustGetBlockWithReadAcquire(blockID types.PtrBindIndex) (types.MemBlockUintptr, bool) {
	var (
		uMemBlock types.MemBlockUintptr
		loaded    bool = false
	)

	p.memBlocksRWMutex.RLock()
	uMemBlock, _ = p.memBlocks[blockID]
	if uMemBlock != 0 {
		loaded = true
		uMemBlock.Ptr().Chunk.Ptr().ReadAcquire()
	}
	p.memBlocksRWMutex.RUnlock()

	if uMemBlock != 0 && p.checkBlock(blockID, uMemBlock) == false {
		uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
		uMemBlock = 0
	}

	if uMemBlock == 0 {
		var (
			pMemBlock    *types.MemBlock
			uNewMemBlock = p.allocChunkFromChunkPool()
		)

		p.memBlocksRWMutex.Lock()

		uMemBlock, _ = p.memBlocks[blockID]
		if uMemBlock == 0 {
			loaded = false
			uMemBlock = uNewMemBlock
			pMemBlock = uMemBlock.Ptr()
			pMemBlock.ID = blockID
			p.memBlocks[blockID] = uMemBlock
		} else {
			loaded = true
			pMemBlock = uMemBlock.Ptr()
		}

		pMemBlock.Chunk.Ptr().ReadAcquire()

		p.memBlocksRWMutex.Unlock()

		if loaded {
			p.releaseChunkToChunkPool(uNewMemBlock)
		} else {
			pMemBlock.CompleteInit()
			p.workingChunkPool.Put(uMemBlock)
		}
	}

	return uMemBlock, loaded
}
