package memstg

import (
	"soloos/sdfs/types"
)

func (p *INodeDriver) preadMemBlock(uINode types.INodeUintptr, memBlockIndex int, data []byte, offset int) error {
	var (
		// isSuccess bool
		err error
	)
	pINode := uINode.Ptr()

	end := offset + len(data)

	// check memblock
	uMemBlock, _ := p.memBlockDriver.MustGetBlockWithReadAcquire(uINode, memBlockIndex)
	// TODO maybe rebase is not needed
	if uMemBlock.Ptr().Contains(offset, end) == false {
		netBlockIndex := memBlockIndex * pINode.MemBlockCap / pINode.NetBlockCap
		uNetBlock := p.netBlockDriver.MustGetBlock(uINode, netBlockIndex)
		err = p.unsafeMemBlockRebaseNetBlock(uINode, uNetBlock, uMemBlock, memBlockIndex)
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

func (p *INodeDriver) PRead(uINode types.INodeUintptr, data []byte, offset int64) error {
	var (
		memBlockIndex       int
		memBlockBytesOffset int
		dataOffset          int
		dataEnd             int
		err                 error
	)
	pINode := uINode.Ptr()
	pINode.AccessRWMutex.RLock()

	// read data from first memblock
	// memBlockIndex = int(math.Ceil(float64(offset) / float64(pINode.MemBlockCap)))
	memBlockIndex = int(offset / int64(pINode.MemBlockCap))
	memBlockBytesOffset = int(offset - (int64(memBlockIndex) * int64(pINode.MemBlockCap)))
	dataOffset = 0
	dataEnd = dataOffset + (pINode.MemBlockCap - memBlockBytesOffset)
	if dataEnd > len(data) {
		dataEnd = len(data)
	}
	err = p.preadMemBlock(uINode, memBlockIndex, data[dataOffset:dataEnd], memBlockBytesOffset)
	if err != nil {
		goto READ_DATA_DONE
	}

	// read data from other memblock
	dataOffset += dataEnd
	for ; dataOffset < len(data); dataOffset += pINode.MemBlockCap {
		dataEnd = dataOffset + pINode.MemBlockCap
		if dataEnd > len(data) {
			dataEnd = len(data)
		}
		memBlockIndex = int((offset + int64(dataOffset)) / int64(pINode.MemBlockCap))
		err = p.preadMemBlock(uINode, memBlockIndex, data[dataOffset:dataEnd], 0)
		if err != nil {
			goto READ_DATA_DONE
		}
	}

READ_DATA_DONE:
	pINode.AccessRWMutex.RUnlock()
	return err
}
