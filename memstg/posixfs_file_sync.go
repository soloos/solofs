package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/sdfsapitypes"
	"soloos/sdfs/sdfstypes"
)

func (p *PosixFS) SimpleFlush(fsINodeID sdfsapitypes.FsINodeID) error {
	var (
		uFsINode sdfsapitypes.FsINodeUintptr
		pFsINode *sdfsapitypes.FsINode
		err      error
	)

	uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(fsINodeID)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	pFsINode = uFsINode.Ptr()
	if err != nil {
		return err
	}

	if pFsINode.UNetINode != 0 {
		err = p.MemStg.NetINodeDriver.Sync(pFsINode.UNetINode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PosixFS) Flush(input *fsapitypes.FlushIn) fsapitypes.Status {
	var err = p.SimpleFlush(input.NodeId)
	return sdfstypes.ErrorToFsStatus(err)
}
