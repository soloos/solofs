package memstg

import (
	"soloos/common/solofstypes"
)

func (p *PosixFs) GetFsINodeByID(fsINodeID solofstypes.FsINodeID) (solofstypes.FsINodeUintptr, error) {
	return p.FsINodeDriver.GetFsINodeByID(fsINodeID)
}

func (p *PosixFs) ReleaseFsINode(uFsINode solofstypes.FsINodeUintptr) {
	p.FsINodeDriver.ReleaseFsINode(uFsINode)
}
