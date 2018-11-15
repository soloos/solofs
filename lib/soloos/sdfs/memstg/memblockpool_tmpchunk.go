package memstg

import "soloos/sdfs/types"

// ChunkPoolInvokeReleaseChunk call by offheap.BlockPool
func (p *memBlockPoolChunk) ChunkPoolInvokeReleaseTmpChunk() {
	uMemBlock := p.takeTmpBlockForRelease()
	if uMemBlock == 0 {
		return
	}
	p.releaseBlock(uMemBlock)
}

func (p *memBlockPoolChunk) takeTmpBlockForRelease() types.MemBlockUintptr {
	iRet := p.workingTmpChunkPool.IteratorAndPop(func(x interface{}) (bool, interface{}) {
		uMemBlock := x.(types.MemBlockUintptr)
		pMemBlock := uMemBlock.Ptr()
		if !pMemBlock.IsInited() && pMemBlock.Chunk.Ptr().Accessor > 0 {
			return false, 0
		}
		return true, uMemBlock
	})

	if iRet == nil {
		iRet = p.workingTmpChunkPool.IteratorAndPop(func(x interface{}) (bool, interface{}) {
			uMemBlock := x.(types.MemBlockUintptr)
			pMemBlock := uMemBlock.Ptr()
			if !pMemBlock.IsInited() {
				return false, 0
			}
			return true, uMemBlock
		})
	}

	if iRet == nil {
		return 0
	}

	return iRet.(types.MemBlockUintptr)
}
