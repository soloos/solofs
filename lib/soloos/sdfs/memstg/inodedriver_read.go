package memstg

import "soloos/sdfs/types"

func (p *INodeDriver) PRead(uINode types.INodeUintptr, data []byte, offset int64) error {
	var (
		// isSuccess bool
		err error
	)
	pINode := uINode.Ptr()
	pINode.AccessRWMutex.RLock()

	// check memblock
	// memBlockIndex := int(offset / int64(pINode.MemBlockSize))
	// memBlockBytesOffset := int(offset - int64(memBlockIndex)*int64(pINode.MemBlockSize))
	// memBlockBytesEnd := memBlockBytesOffset + len(data) + 1
	// uMemBlock, _ := p.memBlockDriver.MustGetBlockWithReadAcquire(uINode, memBlockIndex)
	// if !uMemBlock.Ptr().Contains(memBlockBytesOffset, memBlockBytesEnd) {
	// }

	pINode.AccessRWMutex.RUnlock()
	return err
}
