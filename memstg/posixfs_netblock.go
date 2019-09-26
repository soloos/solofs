package memstg

import (
	"soloos/common/solofsapitypes"
)

func (p *PosixFS) MustGetNetBlock(uNetINode solofsapitypes.NetINodeUintptr,
	netBlockIndex int32) (solofsapitypes.NetBlockUintptr, error) {
	return p.MemStg.NetBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
}

func (p *PosixFS) ReleaseNetBlock(uNetBlock solofsapitypes.NetBlockUintptr) {
	p.MemStg.NetBlockDriver.ReleaseNetBlock(uNetBlock)
}

func (p *PosixFS) NetBlockSetPReadMemBlockWithDisk(preadWithDisk solofsapitypes.PReadMemBlockWithDisk) {
	p.MemStg.netBlockDriver.SetPReadMemBlockWithDisk(preadWithDisk)
}

func (p *PosixFS) NetBlockSetUploadMemBlockWithDisk(uploadMemBlockWithDisk solofsapitypes.UploadMemBlockWithDisk) {
	p.MemStg.netBlockDriver.SetUploadMemBlockWithDisk(uploadMemBlockWithDisk)
}
