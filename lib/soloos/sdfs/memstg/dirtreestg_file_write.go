package memstg

import (
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *DirTreeStg) Write(input *fuse.WriteIn, data []byte) (uint32, fuse.Status) {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return 0, types.ErrorToFuseStatus(err)
	}

	err = p.MemStg.NetINodeDriver.PWriteWithMem(fsINode.UNetINode, data[:input.Size], input.Offset)
	if err != nil {
		return 0, types.ErrorToFuseStatus(err)
	}

	return input.Size, fuse.OK
}
