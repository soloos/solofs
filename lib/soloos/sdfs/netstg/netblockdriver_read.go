package netstg

import (
	"soloos/sdfs/types"
)

func (p *NetBlockDriver) PReadMemBlock(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr, netBlockIndex int,
	uMemBlock types.MemBlockUintptr, memBlockIndex int,
	offset uint64, length int) (int, error) {
	if uNetBlock.Ptr().StorDataBackends.Len == 0 {
		return 0, types.ErrBackendListIsEmpty
	}

	var (
		readedLen int
		err       error
	)

	readedLen, err = p.dataNodeClient.PReadMemBlock(uNetINode, uNetBlock.Ptr().StorDataBackends.Arr[0],
		uNetBlock, netBlockIndex,
		uMemBlock, memBlockIndex,
		offset, length)
	if err != nil {
		return 0, err
	}

	return readedLen, nil
}
