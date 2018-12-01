package memstg

import "soloos/sdfs/types"

func (p *INodeDriver) unsafeLoadMemBlock(uINode types.INodeUintptr,
	uMemBlock types.MemBlockUintptr,
	netBlockIndex int,
	memBlockIndex int) error {
	uMemBlock.Ptr().AvailMask.Set(0, uMemBlock.Ptr().Bytes.Len)
	// uNetBlock, _ := p.netBlockDriver.MustGetBlock(uINode, netBlockIndex)
	return nil
}
