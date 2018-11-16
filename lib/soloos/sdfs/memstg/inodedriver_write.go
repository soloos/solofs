package memstg

import (
	"soloos/sdfs/types"
)

func (p *INodeDriver) PWrite(uINode types.INodeUintptr, data []byte, offset int64) error {
	var err error
	pINode := uINode.Ptr()
	pINode.WriteRWMutex.RLock()

	// write in memblock
	memBlockIndex := int(offset / int64(pINode.MemBlockSize))
	memBlockWriteOffset := int(offset - int64(memBlockIndex)*int64(pINode.MemBlockSize))
	memBlockWriteEnd := memBlockWriteOffset + len(data) + 1
	uMemBlock, _ := p.memBlockDriver.MustGetBlockWithReadAcquire(uINode, memBlockIndex)
	uMemBlock.Ptr().PWrite(data, memBlockWriteOffset)

	// write in netblock
	netBlockIndex := int(offset / int64(pINode.NetBlockSize))
	uNetBlock, _ := p.netBlockDriver.MustGetBlock(uINode, netBlockIndex)
	err = p.netBlockDriver.PWrite(uINode, uNetBlock, uMemBlock, memBlockIndex, memBlockWriteOffset, memBlockWriteEnd)
	if err != nil {
		goto WRITE_DATA_DONE
	}

WRITE_DATA_DONE:
	uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
	pINode.WriteRWMutex.RUnlock()
	return err
}

func (p *INodeDriver) FlushMemBlock(uINode types.INodeUintptr, uMemBlock types.MemBlockUintptr) error {
	var err error
	uINode.Ptr().WriteRWMutex.Lock()
	err = p.netBlockDriver.Flush(uMemBlock)
	if err != nil {
		goto FLUSH_DATA_DONE
	}

FLUSH_DATA_DONE:
	uINode.Ptr().WriteRWMutex.Unlock()
	return err
}
