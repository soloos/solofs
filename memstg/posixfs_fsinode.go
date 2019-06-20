package memstg

import (
	"soloos/common/sdfsapitypes"
)

func (p *PosixFS) GetFsINodeByID(fsINodeID sdfsapitypes.FsINodeID) (sdfsapitypes.FsINodeUintptr, error) {
	return p.FsINodeDriver.GetFsINodeByID(fsINodeID)
}

func (p *PosixFS) ReleaseFsINode(uFsINode sdfsapitypes.FsINodeUintptr) {
	p.FsINodeDriver.ReleaseFsINode(uFsINode)
}
