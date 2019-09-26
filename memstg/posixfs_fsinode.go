package memstg

import (
	"soloos/common/solofsapitypes"
)

func (p *PosixFS) GetFsINodeByID(fsINodeID solofsapitypes.FsINodeID) (solofsapitypes.FsINodeUintptr, error) {
	return p.FsINodeDriver.GetFsINodeByID(fsINodeID)
}

func (p *PosixFS) ReleaseFsINode(uFsINode solofsapitypes.FsINodeUintptr) {
	p.FsINodeDriver.ReleaseFsINode(uFsINode)
}
