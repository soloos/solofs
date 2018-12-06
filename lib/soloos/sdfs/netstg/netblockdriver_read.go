package netstg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

func (p *NetBlockDriver) PRead(uINode types.INodeUintptr,
	uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int,
	offset, length int) error {
	if uNetBlock.Ptr().DataNodes.Len == 0 {
		return types.ErrBackendListIsEmpty
	}

	var (
		resp snettypes.Response
		err  error
	)

	err = p.dataNodeClient.PRead(uNetBlock.Ptr().DataNodes.Arr[0],
		uNetBlock,
		offset, length,
		&resp)
	if err != nil {
		return err
	}

	return nil
}
