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

	chunkPool        offheap.ChunkPool
	workingChunkPool types.HotPool

	tmpChunkPool        offheap.ChunkPool
	workingTmpChunkPool types.HotPool

	memBlocksRWMutex sync.RWMutex
	memBlocks        map[types.PtrBindIndex]types.MemBlockUintptr
}

func (p *memBlockPoolChunk) Init(memBlockPool *MemBlockPool) error {
	var err error

	p.memBlockPool = memBlockPool

	chunksLimit := p.memBlockPool.options.ChunkPoolOptions.ChunksLimit

	chunkPoolOptions := p.memBlockPool.options.ChunkPoolOptions
	p.chunkSize = uintptr(chunkPoolOptions.ChunkSize)
	p.ichunkSize = int(p.chunkSize)
	chunkPoolOptions.ChunksLimit = int32(math.Ceil(float64(chunksLimit) * 0.9))
	chunkPoolOptions.SetChunkPoolAssistant(p.ChunkPoolInvokePrepareNewChunk, p.ChunkPoolInvokeReleaseChunk)
	chunkPoolOptions.ChunkSize = int(uintptr(chunkPoolOptions.ChunkSize) +
		types.MemBlockStructSize)
	err = p.memBlockPool.driver.offheapDriver.InitChunkPool(chunkPoolOptions, &p.chunkPool)
	if err != nil {
		return err
	}
	p.workingChunkPool.Init()

	tmpChunkPoolOptions := chunkPoolOptions
	tmpChunkPoolOptions.ChunksLimit = chunksLimit - chunkPoolOptions.ChunksLimit
	if tmpChunkPoolOptions.ChunksLimit == 0 {
		tmpChunkPoolOptions.ChunksLimit = 1
	}
	tmpChunkPoolOptions.SetChunkPoolAssistant(p.ChunkPoolInvokePrepareNewChunk, p.ChunkPoolInvokeReleaseTmpChunk)
	err = p.memBlockPool.driver.offheapDriver.InitChunkPool(tmpChunkPoolOptions, &p.tmpChunkPool)
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
	p.chunkPool.ReleaseChunk(offheap.ChunkUintptr(
		uMemBlock.Ptr().Bytes.Data -
			types.MemBlockStructSize -
			offheap.ChunkStructDataOffset))
}

func (p *memBlockPoolChunk) takeBlockForRelease() types.MemBlockUintptr {
	iRet := p.workingChunkPool.IteratorAndPop(func(x interface{}) (bool, interface{}) {
		uMemBlock := x.(types.MemBlockUintptr)
		pMemBlock := uMemBlock.Ptr()
		if !pMemBlock.IsInited() && pMemBlock.Chunk.Ptr().Accessor > 0 {
			return false, 0
		}
		return true, uMemBlock
	})

	if iRet != nil {
		return iRet.(types.MemBlockUintptr)
	}

	return (types.MemBlockUintptr)(p.tmpChunkPool.AllocChunk().Ptr().Data)

	// Get Block From workingChunkPool
	// iRet = p.workingChunkPool.IteratorAndPop(func(x interface{}) (bool, interface{}) {
	// uMemBlock := x.(types.MemBlockUintptr)
	// pMemBlock := uMemBlock.Ptr()
	// if !pMemBlock.IsInited() {
	// return false, 0
	// }
	// return true, uMemBlock
	// })

	// if iRet == nil {
	// return 0
	// }

	// return iRet.(types.MemBlockUintptr)
}

func (p *memBlockPoolChunk) beforeReleaseBlock(pMemBlock *types.MemBlock) {
	if pMemBlock.IsInited() == false {
		return
	}
	pMemBlock.SetReleasable()
}

func (p *memBlockPoolChunk) releaseBlock(uMemBlock types.MemBlockUintptr) {
	pMemBlock := uMemBlock.Ptr()

	p.beforeReleaseBlock(pMemBlock)
	pMemBlock.Chunk.Ptr().WriteAcquire()
	if pMemBlock.EnsureRelease() {
		p.memBlocksRWMutex.Lock()
		delete(p.memBlocks, pMemBlock.MemID)
		p.memBlocksRWMutex.Unlock()
		pMemBlock.Chunk.Ptr().WriteRelease()
		p.releaseChunkToChunkPool(uMemBlock)
	} else {
		pMemBlock.Chunk.Ptr().WriteRelease()
	}
}

func (p *memBlockPoolChunk) checkBlock(blockID types.PtrBindIndex, uMemBlock types.MemBlockUintptr) bool {
	if uMemBlock.Ptr().MemID != blockID {
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

	if uMemBlock != 0 && !p.checkBlock(blockID, uMemBlock) {
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
			pMemBlock.MemID = blockID
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
