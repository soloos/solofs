package memstg

import (
	"soloos/common/solofsapitypes"
)

func (p *NetBlockDriver) PReadMemBlock(uNetINode solofsapitypes.NetINodeUintptr,
	uNetBlock solofsapitypes.NetBlockUintptr, netBlockIndex int32,
	uMemBlock solofsapitypes.MemBlockUintptr, memBlockIndex int32,
	offset uint64, length int) (int, error) {
	if uNetBlock.Ptr().StorDataBackends.Len == 0 {
		return 0, solofsapitypes.ErrBackendListIsEmpty
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
