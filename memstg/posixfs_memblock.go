package memstg

import (
	"soloos/common/solofsapitypes"
)

func (p *PosixFs) MustGetMemBlockWithReadAcquire(uNetINode solofsapitypes.NetINodeUintptr,
	memBlockIndex int32) (solofsapitypes.MemBlockUintptr, bool) {
	return p.MemStg.MemBlockDriver.MustGetMemBlockWithReadAcquire(uNetINode, memBlockIndex)
}

func (p *PosixFs) ReleaseMemBlockWithReadRelease(uMemBlock solofsapitypes.MemBlockUintptr) {
	p.MemStg.MemBlockDriver.ReleaseMemBlockWithReadRelease(uMemBlock)
}
