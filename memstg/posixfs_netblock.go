package memstg

import (
	"soloos/common/solofstypes"
)

func (p *PosixFs) MustGetNetBlock(uNetINode solofstypes.NetINodeUintptr,
	netBlockIndex int32) (solofstypes.NetBlockUintptr, error) {
	return p.MemStg.NetBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
}

func (p *PosixFs) ReleaseNetBlock(uNetBlock solofstypes.NetBlockUintptr) {
	p.MemStg.NetBlockDriver.ReleaseNetBlock(uNetBlock)
}

func (p *PosixFs) NetBlockSetPReadMemBlockWithDisk(preadWithDisk solofstypes.PReadMemBlockWithDisk) {
	p.MemStg.netBlockDriver.SetPReadMemBlockWithDisk(preadWithDisk)
}

func (p *PosixFs) NetBlockSetUploadMemBlockWithDisk(uploadMemBlockWithDisk solofstypes.UploadMemBlockWithDisk) {
	p.MemStg.netBlockDriver.SetUploadMemBlockWithDisk(uploadMemBlockWithDisk)
}
