package netstg

import (
	"soloos/common/sdfsapitypes"
)

func (p *NetBlockDriver) PWrite(uNetINode sdfsapitypes.NetINodeUintptr,
	uNetBlock sdfsapitypes.NetBlockUintptr, netBlockIndex int32,
	uMemBlock sdfsapitypes.MemBlockUintptr, memBlockIndex int32,
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
