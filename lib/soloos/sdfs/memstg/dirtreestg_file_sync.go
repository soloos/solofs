package memstg

import (
	fsapitypes "soloos/common/fsapi/types"
	"soloos/sdfs/types"
)

func (p *DirTreeStg) SimpleFlush(uNetINode types.NetINodeUintptr) error {
	return p.MemStg.NetINodeDriver.Sync(uNetINode)
}

func (p *DirTreeStg) Flush(input *fsapitypes.FlushIn) fsapitypes.Status {
	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	if fsINode.UNetINode != 0 {
		err = p.SimpleFlush(fsINode.UNetINode)
		if err != nil {
			return types.ErrorToFsStatus(err)
		}
	}

	return fsapitypes.OK
}
