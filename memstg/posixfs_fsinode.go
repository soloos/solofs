package memstg

import (
	"soloos/common/solofstypes"
)

func (p *PosixFs) GetFsINodeByID(fsINodeIno solofstypes.FsINodeIno) (solofstypes.FsINodeUintptr, error) {
	return p.FsINodeDriver.GetFsINodeByID(fsINodeIno)
}

func (p *PosixFs) ReleaseFsINode(uFsINode solofstypes.FsINodeUintptr) {
	p.FsINodeDriver.ReleaseFsINode(uFsINode)
}
