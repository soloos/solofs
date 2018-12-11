package memstg

import (
	"soloos/sdfs/types"
)

func (p *NetINodeDriver) PWrite(uNetINode types.NetINodeUintptr, data []byte, offset int64) error {
	var (
		isSuccess           bool
		memBlockIndex       int
		memBlockBytesOffset int
		memBlockBytesEnd    int
		uMemBlock           types.MemBlockUintptr
		netBlockIndex       int
		uNetBlock           types.NetBlockUintptr
		err                 error
	)
	pNetINode := uNetINode.Ptr()
	pNetINode.MetaDataMutex.RLock()

	// write in memblock
	memBlockIndex = int(offset / int64(pNetINode.MemBlockCap))
	memBlockBytesOffset = int(offset - int64(memBlockIndex)*int64(pNetINode.MemBlockCap))
	memBlockBytesEnd = memBlockBytesOffset + len(data)
	uMemBlock, _ = p.memBlockDriver.MustGetBlockWithReadAcquire(uNetINode, memBlockIndex)
	isSuccess = uMemBlock.Ptr().PWrite(data, memBlockBytesOffset)
	if isSuccess == false {
		// TODO memblock load data
		panic("write error")
	}

	// write in netblock
	netBlockIndex = int(offset / int64(pNetINode.NetBlockCap))
	uNetBlock, err = p.netBlockDriver.MustGetBlock(uNetINode, netBlockIndex)
	if err != nil {
		goto WRITE_DATA_DONE
	}

	err = p.netBlockDriver.PWrite(uNetINode, uNetBlock, uMemBlock, memBlockIndex, memBlockBytesOffset, memBlockBytesEnd)
	if err != nil {
		goto WRITE_DATA_DONE
	}

WRITE_DATA_DONE:
	uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
	pNetINode.MetaDataMutex.RUnlock()
	return err
}

func (p *NetINodeDriver) FlushMemBlock(uNetINode types.NetINodeUintptr,
	uMemBlock types.MemBlockUintptr) error {
	var err error
	uNetINode.Ptr().MetaDataMutex.Lock()
	err = p.netBlockDriver.FlushMemBlock(uMemBlock)
	if err != nil {
		goto FLUSH_DATA_DONE
	}

FLUSH_DATA_DONE:
	uNetINode.Ptr().MetaDataMutex.Unlock()
	return err
}
