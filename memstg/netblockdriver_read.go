package memstg

import (
	"soloos/common/solofstypes"
)

func (p *NetBlockDriver) PReadMemBlock(uNetINode solofstypes.NetINodeUintptr,
	uNetBlock solofstypes.NetBlockUintptr, netBlockIndex int32,
	uMemBlock solofstypes.MemBlockUintptr, memBlockIndex int32,
	offset uint64, length int) (int, error) {
	if uNetBlock.Ptr().StorDataBackends.Len == 0 {
		return 0, solofstypes.ErrBackendListIsEmpty
	}

	var (
		readedLen int
		err       error
	)

	// TODO choose solodn to read
	readedLen, err = p.PReadMemBlockFromNet(uNetINode,
		uNetBlock, netBlockIndex,
		uMemBlock, memBlockIndex,
		offset, length)
	if err != nil {
		return 0, err
	}

	return readedLen, nil
}
