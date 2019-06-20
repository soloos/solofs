package memstg

import (
	"soloos/common/sdfsapitypes"
)

func (p *NetBlockDriver) PReadMemBlock(uNetINode sdfsapitypes.NetINodeUintptr,
	uNetBlock sdfsapitypes.NetBlockUintptr, netBlockIndex int32,
	uMemBlock sdfsapitypes.MemBlockUintptr, memBlockIndex int32,
	offset uint64, length int) (int, error) {
	if uNetBlock.Ptr().StorDataBackends.Len == 0 {
		return 0, sdfsapitypes.ErrBackendListIsEmpty
	}

	var (
		readedLen int
		err       error
	)

	// TODO choose datanode to read
	readedLen, err = p.dataNodeClient.PReadMemBlock(uNetINode,
		uNetBlock, netBlockIndex,
		uMemBlock, memBlockIndex,
		offset, length)
	if err != nil {
		return 0, err
	}

	return readedLen, nil
}
