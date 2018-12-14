package netstg

import (
	"soloos/sdfs/types"
)

func (p *NetBlockDriver) PWrite(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int,
	offset, end int) error {
	var err error

	err = p.netBlockDriverUploader.PWrite(uNetINode, uNetBlock, uMemBlock, memBlockIndex, offset, end)
	if err != nil {
		return err
	}

	return nil
}
