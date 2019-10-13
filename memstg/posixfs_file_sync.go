package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/solofstypes"
)

func (p *PosixFs) SimpleFlush(fsINodeIno solofstypes.FsINodeIno) error {
	var (
		uFsINode solofstypes.FsINodeUintptr
		pFsINode *solofstypes.FsINode
		err      error
	)

	uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(fsINodeIno)
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

func (p *PosixFs) Flush(input *fsapi.FlushIn) fsapi.Status {
	var err = p.SimpleFlush(input.NodeId)
	return ErrorToFsStatus(err)
}
