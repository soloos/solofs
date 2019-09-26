package memstg

import (
	"soloos/common/solofsapitypes"
)

func (p *NetBlockDriver) PWrite(uNetINode solofsapitypes.NetINodeUintptr,
	uNetBlock solofsapitypes.NetBlockUintptr, netBlockIndex int32,
	uMemBlock solofsapitypes.MemBlockUintptr, memBlockIndex int32,
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
