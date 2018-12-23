package netstg

import (
	"soloos/sdfs/types"
)

func (p *NetBlockDriver) PReadMemBlock(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr, netBlockIndex int,
	uMemBlock types.MemBlockUintptr, memBlockIndex int,
	offset int64, length int) error {
	if uNetBlock.Ptr().StorDataBackends.Len == 0 {
		return types.ErrBackendListIsEmpty
	}

	var (
		err error
	)

	err = p.dataNodeClient.PReadMemBlock(uNetINode, uNetBlock.Ptr().StorDataBackends.Arr[0],
		uNetBlock, netBlockIndex,
		uMemBlock, memBlockIndex,
		offset, length)
	if err != nil {
		return err
	}

	return nil
}
