package memstg

import (
	fsapitypes "soloos/fsapi/types"
	"soloos/sdfs/types"
)

func (p *DirTreeStg) SimpleWriteWithMem(uNetINode types.NetINodeUintptr,
	data []byte, offset uint64) error {
	return p.MemStg.NetINodeDriver.PWriteWithMem(uNetINode, data, offset)
}

func (p *DirTreeStg) Write(input *fsapitypes.WriteIn, data []byte) (uint32, fsapitypes.Status) {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return 0, types.ErrorToFsStatus(err)
	}

	err = p.SimpleWriteWithMem(fsINode.UNetINode, data[:input.Size], input.Offset)
	if err != nil {
		return 0, types.ErrorToFsStatus(err)
	}

	return input.Size, fsapitypes.OK
}
