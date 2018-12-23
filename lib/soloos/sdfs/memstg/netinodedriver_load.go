package memstg

import (
	"soloos/sdfs/types"
	"soloos/util/offheap"
)

func (p *NetINodeDriver) unsafeMemBlockRebaseNetBlock(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr,
	netBlockIndex int,
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

	uTmpMemBlock = p.memBlockDriver.AllocTmpBlockWithWriteAcquire(uNetINode)
	err = p.netBlockDriver.PReadMemBlock(uNetINode, uNetBlock, netBlockIndex, uTmpMemBlock, memBlockIndex,
		int64(memBlockIndex)*int64(uNetINode.Ptr().MemBlockCap),
		uNetINode.Ptr().MemBlockCap,
	)
	if err != nil {
		goto READ_DONE
	}

	pMemBlock.AvailMask.MergeElementRWMutex.Lock()
	for i := 0; i < pMemBlock.AvailMask.MaskArrayLen; i++ {
		chunkMaskEntry = pMemBlock.AvailMask.MaskArray[i]
		uTmpMemBlock.Ptr().PWriteWithMem(
			(*pMemBlock.BytesSlice())[chunkMaskEntry.Offset:chunkMaskEntry.End],
			chunkMaskEntry.Offset)
	}
	copy(*pMemBlock.BytesSlice(), *uTmpMemBlock.Ptr().BytesSlice())
	pMemBlock.AvailMask.Set(0, pMemBlock.Bytes.Len)
	pMemBlock.AvailMask.MergeElementRWMutex.Unlock()

READ_DONE:
	uTmpMemBlock.Ptr().Chunk.Ptr().WriteRelease()
	p.memBlockDriver.ReleaseTmpBlock(uNetINode, uTmpMemBlock)
	pMemBlock.RebaseNetBlockMutex.Unlock()

	return err
}
