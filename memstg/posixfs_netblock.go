package memstg

import (
	"soloos/common/sdfsapitypes"
)

func (p *PosixFS) MustGetNetBlock(uNetINode sdfsapitypes.NetINodeUintptr,
	netBlockIndex int32) (sdfsapitypes.NetBlockUintptr, error) {
	return p.MemStg.NetBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
}

func (p *PosixFS) ReleaseNetBlock(uNetBlock sdfsapitypes.NetBlockUintptr) {
	p.MemStg.NetBlockDriver.ReleaseNetBlock(uNetBlock)
}

func (p *PosixFS) NetBlockSetPReadMemBlockWithDisk(preadWithDisk sdfsapitypes.PReadMemBlockWithDisk) {
	p.MemStg.netBlockDriver.SetPReadMemBlockWithDisk(preadWithDisk)
}

func (p *PosixFS) NetBlockSetUploadMemBlockWithDisk(uploadMemBlockWithDisk sdfsapitypes.UploadMemBlockWithDisk) {
	p.MemStg.netBlockDriver.SetUploadMemBlockWithDisk(uploadMemBlockWithDisk)
}
