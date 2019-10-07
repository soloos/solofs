package memstg

import (
	"soloos/common/solofstypes"
)

func (p *PosixFs) MustGetMemBlockWithReadAcquire(uNetINode solofstypes.NetINodeUintptr,
	memBlockIndex int32) (solofstypes.MemBlockUintptr, bool) {
	return p.MemStg.MemBlockDriver.MustGetMemBlockWithReadAcquire(uNetINode, memBlockIndex)
}

func (p *PosixFs) ReleaseMemBlockWithReadRelease(uMemBlock solofstypes.MemBlockUintptr) {
	p.MemStg.MemBlockDriver.ReleaseMemBlockWithReadRelease(uMemBlock)
}
