package memstg

import "soloos/sdfs/types"

func (p *INodeDriver) WriteAt(uINode types.INodeUintptr, data []byte, offset int64) error {
	pINode := uINode.Ptr()

	// write in memblock
	memBlockIndex := int(offset / int64(pINode.MemBlockSize))
	memBlockWriteOffset := int(offset - int64(memBlockIndex)*int64(pINode.MemBlockSize))
	uMemBlock, _ := p.memBlockDriver.MustGetBlockWithReadAcquire(uINode, memBlockIndex)
	memBlockData := *uMemBlock.Ptr().BytesSlice()
	copy(memBlockData[memBlockWriteOffset:], data)

	// write in netblock
	// netBlockIndex := int(offset / int64(pINode.NetBlockSize))
	// uNetBlock, _ := p.netBlockDriver.MustGetBlock(uINode, netBlockIndex)

	uMemBlock.Ptr().Chunk.Ptr().ReadRelease()

	return nil
}
