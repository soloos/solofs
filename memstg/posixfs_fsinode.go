package memstg

import (
	"soloos/common/solofsapitypes"
)

func (p *PosixFs) GetFsINodeByID(fsINodeID solofsapitypes.FsINodeID) (solofsapitypes.FsINodeUintptr, error) {
	return p.FsINodeDriver.GetFsINodeByID(fsINodeID)
}

func (p *PosixFs) ReleaseFsINode(uFsINode solofsapitypes.FsINodeUintptr) {
	p.FsINodeDriver.ReleaseFsINode(uFsINode)
}
