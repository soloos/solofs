package memstg

import (
	"soloos/sdfs/types"
	"soloos/util/offheap"
)

func (p *INodeDriver) unsafeMemBlockRebaseNetBlock(uINode types.INodeUintptr,
	uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int) error {
	var (
		chunkMaskEntry offheap.ChunkMaskEntry
		pMemBlock      *types.MemBlock
		uTmpMemBlock   types.MemBlockUintptr
		err            error
	)

	pMemBlock = uMemBlock.Ptr()
	pMemBlock.RebaseNetBlockMutex.Lock()
	if pMemBlock.Contains(0, pMemBlock.Bytes.Len) {
		pMemBlock.RebaseNetBlockMutex.Unlock()
		return nil
	}

	uTmpMemBlock = p.memBlockDriver.AllocTmpBlockWithWriteAcquire(uINode)
	err = p.netBlockDriver.PRead(uINode, uNetBlock, uTmpMemBlock, memBlockIndex,
		memBlockIndex*uINode.Ptr().MemBlockCap,
		uINode.Ptr().MemBlockCap,
	)
	if err != nil {
		goto READ_DONE
	}

	pMemBlock.AvailMask.MergeElementRWMutex.Lock()
	for i := 0; i < pMemBlock.AvailMask.MaskArrayLen; i++ {
		chunkMaskEntry = pMemBlock.AvailMask.MaskArray[i]
		uTmpMemBlock.Ptr().PWrite(
			(*pMemBlock.BytesSlice())[chunkMaskEntry.Offset:chunkMaskEntry.End],
			chunkMaskEntry.Offset)
	}
	copy(*pMemBlock.BytesSlice(), *uTmpMemBlock.Ptr().BytesSlice())
	pMemBlock.AvailMask.Set(0, pMemBlock.Bytes.Len)
	pMemBlock.AvailMask.MergeElementRWMutex.Unlock()

READ_DONE:
	uTmpMemBlock.Ptr().Chunk.Ptr().WriteRelease()
	p.memBlockDriver.ReleaseTmpBlock(uINode, uTmpMemBlock)
	pMemBlock.RebaseNetBlockMutex.Unlock()

	return err
}
