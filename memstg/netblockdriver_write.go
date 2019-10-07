package memstg

import (
	"soloos/common/solofstypes"
)

func (p *NetBlockDriver) PWrite(uNetINode solofstypes.NetINodeUintptr,
	uNetBlock solofstypes.NetBlockUintptr, netBlockIndex int32,
	uMemBlock solofstypes.MemBlockUintptr, memBlockIndex int32,
	offset, end int) error {
	var (
		err error
	)

	err = p.NetBlockUploader.PWrite(uNetINode,
		uNetBlock, netBlockIndex,
		uMemBlock, memBlockIndex,
		offset, end)
	if err != nil {
		return err
	}

	return nil
}
