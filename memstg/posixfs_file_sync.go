package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofsapitypes"
	"soloos/solofs/solofstypes"
)

func (p *PosixFs) SimpleFlush(fsINodeID solofsapitypes.FsINodeID) error {
	var (
		uFsINode solofsapitypes.FsINodeUintptr
		pFsINode *solofsapitypes.FsINode
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

func (p *PosixFs) Flush(input *fsapitypes.FlushIn) fsapitypes.Status {
	var err = p.SimpleFlush(input.NodeId)
	return solofstypes.ErrorToFsStatus(err)
}
