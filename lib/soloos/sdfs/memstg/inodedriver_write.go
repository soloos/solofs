package memstg

import (
	"soloos/sdfs/types"
)

func (p *INodeDriver) PWrite(uINode types.INodeUintptr, data []byte, offset int64) error {
	var (
		isSuccess bool
		err       error
	)
	pINode := uINode.Ptr()
	pINode.AccessRWMutex.RLock()

	// write in memblock
	memBlockIndex := int(offset / int64(pINode.MemBlockSize))
	memBlockBytesOffset := int(offset - int64(memBlockIndex)*int64(pINode.MemBlockSize))
	memBlockBytesEnd := memBlockBytesOffset + len(data)
	uMemBlock, _ := p.memBlockDriver.MustGetBlockWithReadAcquire(uINode, memBlockIndex)
	isSuccess = uMemBlock.Ptr().PWrite(data, memBlockBytesOffset)
	if !isSuccess {
		// TODO memblock load data
		panic("write error")
	}

	// write in netblock
	netBlockIndex := int(offset / int64(pINode.NetBlockSize))
	uNetBlock, _ := p.netBlockDriver.MustGetBlock(uINode, netBlockIndex)
	err = p.netBlockDriver.PWrite(uINode, uNetBlock, uMemBlock, memBlockIndex, memBlockBytesOffset, memBlockBytesEnd)
	if err != nil {
		goto WRITE_DATA_DONE
	}

WRITE_DATA_DONE:
	uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
	pINode.AccessRWMutex.RUnlock()
	return err
}

func (p *INodeDriver) FlushMemBlock(uINode types.INodeUintptr, uMemBlock types.MemBlockUintptr) error {
	var err error
	uINode.Ptr().AccessRWMutex.Lock()
	err = p.netBlockDriver.Flush(uMemBlock)
	if err != nil {
		goto FLUSH_DATA_DONE
	}

FLUSH_DATA_DONE:
	uINode.Ptr().AccessRWMutex.Unlock()
	return err
}
