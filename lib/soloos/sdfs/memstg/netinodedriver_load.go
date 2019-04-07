package memstg

import (
	"soloos/sdbone/offheap"
	"soloos/sdfs/types"
)

func (p *NetINodeDriver) unsafeMemBlockRebaseNetBlock(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr,
	netBlockIndex int32,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int32) error {
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

	uTmpMemBlock = p.memBlockDriver.MustGetTmpMemBlockWithReadAcquire(uNetINode, pMemBlock.ID)
	_, err = p.netBlockDriver.PReadMemBlock(uNetINode, uNetBlock, netBlockIndex, uTmpMemBlock, memBlockIndex,
		uint64(memBlockIndex)*uint64(uNetINode.Ptr().MemBlockCap),
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
	uTmpMemBlock.Ptr().ReadRelease()
	p.memBlockDriver.ReleaseTmpMemBlock(uTmpMemBlock)
	pMemBlock.RebaseNetBlockMutex.Unlock()

	return err
}
