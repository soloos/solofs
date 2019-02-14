package memstg

import (
	"soloos/sdfs/types"
)

// ChunkPoolInvokeReleaseChunk call by offheap.BlockPool
func (p *memBlockPoolChunk) ChunkPoolInvokeReleaseTmpChunk() {
	uMemBlock := p.takeTmpBlockForRelease()
	if uMemBlock == 0 {
		return
	}
	p.ReleaseTmpMemBlock(uMemBlock)
}

func (p *memBlockPoolChunk) allocTmpChunkFromTmpChunkPool() types.MemBlockUintptr {
	uChunk := p.tmpChunkPool.AllocChunk()
	uMemBlock := (types.MemBlockUintptr)(uChunk.Ptr().Data)
	return uMemBlock
}

func (p *memBlockPoolChunk) releaseTmpChunkToTmpChunkPool(uMemBlock types.MemBlockUintptr) {
	uMemBlock.Ptr().Reset()
	p.tmpChunkPool.ReleaseChunk(uMemBlock.Ptr().Chunk)
}

func (p *memBlockPoolChunk) takeTmpBlockForRelease() types.MemBlockUintptr {
	iRet := p.workingTmpChunkPool.IteratorAndPop(func(x uintptr) (bool, uintptr) {
		uMemBlock := types.MemBlockUintptr(x)
		pMemBlock := uMemBlock.Ptr()
		if pMemBlock.IsInited() == false && pMemBlock.Chunk.Ptr().Accessor > 0 {
			return false, 0
		}
		return true, uintptr(uMemBlock)
	})

	if iRet == 0 {
		iRet = p.workingTmpChunkPool.IteratorAndPop(func(x uintptr) (bool, uintptr) {
			uMemBlock := types.MemBlockUintptr(x)
			pMemBlock := uMemBlock.Ptr()
			if pMemBlock.IsInited() == false {
				return false, 0
			}
			return true, uintptr(uMemBlock)
		})
	}

	if iRet == 0 {
		return 0
	}

	return types.MemBlockUintptr(iRet)
}

func (p *memBlockPoolChunk) beforeReleaseTmpMemBlock(pMemBlock *types.MemBlock) {
	if pMemBlock.IsInited() == false {
		return
	}
	pMemBlock.SetReleasable()
}

func (p *memBlockPoolChunk) ReleaseTmpMemBlock(uMemBlock types.MemBlockUintptr) {
	pMemBlock := uMemBlock.Ptr()

	p.beforeReleaseTmpMemBlock(pMemBlock)
	pMemBlock.Chunk.Ptr().WriteAcquire()
	if pMemBlock.EnsureRelease() {
		p.workingTmpChunkPool.IteratorAndPop(func(x uintptr) (bool, uintptr) {
			uLocalMemBlock := types.MemBlockUintptr(x)
			if uLocalMemBlock == uMemBlock {
				return true, uintptr(uLocalMemBlock)
			}
			return false, 0
		})

		p.releaseTmpChunkToTmpChunkPool(uMemBlock)
		pMemBlock.Chunk.Ptr().WriteRelease()
	} else {
		pMemBlock.Chunk.Ptr().WriteRelease()
	}
}

func (p *memBlockPoolChunk) MustGetTmpMemBlockWithReadAcquire(blockID types.PtrBindIndex) types.MemBlockUintptr {
	uMemBlock := p.allocTmpChunkFromTmpChunkPool()
	uMemBlock.Ptr().Chunk.Ptr().WriteAcquire()
	uMemBlock.Ptr().CompleteInit()
	p.workingTmpChunkPool.Put(uintptr(uMemBlock))
	return uMemBlock
}
