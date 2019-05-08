package memstg

import (
	fsapitypes "soloos/common/fsapi/types"
	sdfsapitypes "soloos/common/sdfsapi/types"
	"soloos/sdfs/types"
)

func (p *DirTreeStg) SimpleWriteWithMem(fsINodeID sdfsapitypes.FsINodeID,
	data []byte, offset uint64) error {

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

	return p.MemStg.NetINodeDriver.PWriteWithMem(pFsINode.UNetINode, data, offset)
}

func (p *DirTreeStg) Write(input *fsapitypes.WriteIn, data []byte) (uint32, fsapitypes.Status) {
	var err error
	err = p.SimpleWriteWithMem(input.NodeId, data[:input.Size], input.Offset)
	if err != nil {
		return 0, types.ErrorToFsStatus(err)
	}

	return input.Size, fsapitypes.OK
}
