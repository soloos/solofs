package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofsapitypes"
	"soloos/solofs/solofstypes"
)

func (p *PosixFs) SimpleWriteWithMem(fsINodeID solofsapitypes.FsINodeID,
	data []byte, offset uint64) error {

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

	return p.MemStg.NetINodeDriver.PWriteWithMem(pFsINode.UNetINode, data, offset)
}

func (p *PosixFs) Write(input *fsapitypes.WriteIn, data []byte) (uint32, fsapitypes.Status) {
	var err error
	err = p.SimpleWriteWithMem(input.NodeId, data[:input.Size], input.Offset)
	if err != nil {
		return 0, solofstypes.ErrorToFsStatus(err)
	}

	return input.Size, fsapitypes.OK
}
