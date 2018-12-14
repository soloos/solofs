package memstg

import (
	"soloos/sdfs/types"
)

func (p *NetINodeDriver) preadMemBlock(uNetINode types.NetINodeUintptr, memBlockIndex int, data []byte, offset int) error {
	var (
		// isSuccess bool
		err error
	)
	pNetINode := uNetINode.Ptr()

	end := offset + len(data)

	// check memblock
	uMemBlock, _ := p.memBlockDriver.MustGetBlockWithReadAcquire(uNetINode, memBlockIndex)
	// TODO maybe rebase is not needed
	if uMemBlock.Ptr().Contains(offset, end) == false {
		var (
			netBlockIndex int
			uNetBlock     types.NetBlockUintptr
		)
		netBlockIndex = memBlockIndex * pNetINode.MemBlockCap / pNetINode.NetBlockCap
		uNetBlock, err = p.netBlockDriver.MustGetBlock(uNetINode, netBlockIndex)
		if err != nil {
			goto READ_DATA_DONE
		}

		err = p.unsafeMemBlockRebaseNetBlock(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex)
		if err != nil {
			goto READ_DATA_DONE
		}
	}

	// read memblock
	uMemBlock.Ptr().PRead(data, offset)

READ_DATA_DONE:
	uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
	return err
}

func (p *NetINodeDriver) PRead(uNetINode types.NetINodeUintptr, data []byte, offset int64) error {
	var (
		memBlockIndex       int
		memBlockBytesOffset int
		dataOffset          int
		dataEnd             int
		err                 error
	)
	pNetINode := uNetINode.Ptr()
	pNetINode.MetaDataMutex.RLock()

	// read data from first memblock
	// memBlockIndex = int(math.Ceil(float64(offset) / float64(pNetINode.MemBlockCap)))
	memBlockIndex = int(offset / int64(pNetINode.MemBlockCap))
	memBlockBytesOffset = int(offset - (int64(memBlockIndex) * int64(pNetINode.MemBlockCap)))
	dataOffset = 0
	dataEnd = dataOffset + (pNetINode.MemBlockCap - memBlockBytesOffset)
	if dataEnd > len(data) {
		dataEnd = len(data)
	}
	err = p.preadMemBlock(uNetINode, memBlockIndex, data[dataOffset:dataEnd], memBlockBytesOffset)
	if err != nil {
		goto READ_DATA_DONE
	}

	// read data from other memblock
	dataOffset += dataEnd
	for ; dataOffset < len(data); dataOffset += pNetINode.MemBlockCap {
		dataEnd = dataOffset + pNetINode.MemBlockCap
		if dataEnd > len(data) {
			dataEnd = len(data)
		}
		memBlockIndex = int((offset + int64(dataOffset)) / int64(pNetINode.MemBlockCap))
		err = p.preadMemBlock(uNetINode, memBlockIndex, data[dataOffset:dataEnd], 0)
		if err != nil {
			goto READ_DATA_DONE
		}
	}

READ_DATA_DONE:
	pNetINode.MetaDataMutex.RUnlock()
	return err
}
