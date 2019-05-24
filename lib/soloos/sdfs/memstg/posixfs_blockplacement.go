package memstg

import (
	"soloos/common/sdfsapitypes"
)

func (p *PosixFS) SetFsINodeBlockPlacement(fsINodeID sdfsapitypes.FsINodeID,
	policy sdfsapitypes.MemBlockPlacementPolicy) error {
	var (
		uFsINode sdfsapitypes.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.FsINodeDriver.GetFsINodeByID(fsINodeID)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}

	if uFsINode.Ptr().UNetINode == 0 {
		return sdfsapitypes.ErrNetINodeNotExists
	}

	uFsINode.Ptr().UNetINode.Ptr().MemBlockPlacementPolicy = policy

	return nil
}
