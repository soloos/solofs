package memstg

import (
	"soloos/common/solofsapitypes"
)

func (p *PosixFs) MustGetNetBlock(uNetINode solofsapitypes.NetINodeUintptr,
	netBlockIndex int32) (solofsapitypes.NetBlockUintptr, error) {
	return p.MemStg.NetBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
}

func (p *PosixFs) ReleaseNetBlock(uNetBlock solofsapitypes.NetBlockUintptr) {
	p.MemStg.NetBlockDriver.ReleaseNetBlock(uNetBlock)
}

func (p *PosixFs) NetBlockSetPReadMemBlockWithDisk(preadWithDisk solofsapitypes.PReadMemBlockWithDisk) {
	p.MemStg.netBlockDriver.SetPReadMemBlockWithDisk(preadWithDisk)
}

func (p *PosixFs) NetBlockSetUploadMemBlockWithDisk(uploadMemBlockWithDisk solofsapitypes.UploadMemBlockWithDisk) {
	p.MemStg.netBlockDriver.SetUploadMemBlockWithDisk(uploadMemBlockWithDisk)
}
