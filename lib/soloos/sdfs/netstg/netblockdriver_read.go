package netstg

import (
	"soloos/sdfs/types"
)

func (p *NetBlockDriver) PRead(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr, netBlockIndex int,
	uMemBlock types.MemBlockUintptr, memBlockIndex int,
	offset, length int) error {
	if uNetBlock.Ptr().StorDataBackends.Len == 0 {
		return types.ErrBackendListIsEmpty
	}

	var (
		err error
	)

	err = p.dataNodeClient.PRead(uNetBlock.Ptr().StorDataBackends.Arr[0],
		uNetBlock, netBlockIndex,
		uMemBlock, memBlockIndex,
		offset, length)
	if err != nil {
		return err
	}

	return nil
}
