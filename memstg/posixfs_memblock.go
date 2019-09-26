package memstg

import (
	"soloos/common/solofsapitypes"
)

func (p *PosixFS) MustGetMemBlockWithReadAcquire(uNetINode solofsapitypes.NetINodeUintptr,
	memBlockIndex int32) (solofsapitypes.MemBlockUintptr, bool) {
	return p.MemStg.MemBlockDriver.MustGetMemBlockWithReadAcquire(uNetINode, memBlockIndex)
}

func (p *PosixFS) ReleaseMemBlockWithReadRelease(uMemBlock solofsapitypes.MemBlockUintptr) {
	p.MemStg.MemBlockDriver.ReleaseMemBlockWithReadRelease(uMemBlock)
}
