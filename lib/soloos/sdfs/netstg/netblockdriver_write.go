package netstg

import (
	"soloos/sdfs/types"
)

func (p *NetBlockDriver) PWrite(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr, netBlockIndex int32,
	uMemBlock types.MemBlockUintptr, memBlockIndex int32,
	offset, end int) error {
	var (
		err error
	)

	err = p.netBlockDriverUploader.PWrite(uNetINode,
		uNetBlock, netBlockIndex,
		uMemBlock, memBlockIndex,
		offset, end)
	if err != nil {
		return err
	}

	return nil
}
