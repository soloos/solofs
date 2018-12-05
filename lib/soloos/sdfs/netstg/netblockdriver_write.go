package netstg

import (
	"soloos/sdfs/types"
)

func (p *NetBlockDriver) PWrite(uINode types.INodeUintptr,
	uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int,
	offset, end int) error {
	var err error

	err = p.netBlockDriverUploader.PWrite(uNetBlock, uMemBlock, memBlockIndex, offset, end)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetBlockDriver) FlushMemBlock(uMemBlock types.MemBlockUintptr) error {
	return p.netBlockDriverUploader.FlushMemBlock(uMemBlock)
}
