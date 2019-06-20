package memstg

import (
	"soloos/common/sdfsapitypes"
)

func (p *PosixFS) MustGetMemBlockWithReadAcquire(uNetINode sdfsapitypes.NetINodeUintptr,
	memBlockIndex int32) (sdfsapitypes.MemBlockUintptr, bool) {
	return p.MemStg.MemBlockDriver.MustGetMemBlockWithReadAcquire(uNetINode, memBlockIndex)
}

func (p *PosixFS) ReleaseMemBlockWithReadRelease(uMemBlock sdfsapitypes.MemBlockUintptr) {
	p.MemStg.MemBlockDriver.ReleaseMemBlockWithReadRelease(uMemBlock)
}
