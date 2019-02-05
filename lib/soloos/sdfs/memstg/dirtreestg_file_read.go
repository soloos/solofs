package memstg

import (
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

type ReadResult struct {
	dataLen int
}

var _ = fuse.ReadResult(&ReadResult{})

func (p ReadResult) Bytes(buf []byte) ([]byte, fuse.Status) {
	return buf[:p.dataLen], fuse.OK
}

func (p ReadResult) Size() int {
	return p.dataLen
}

func (p ReadResult) Done() {
}

func (p *DirTreeStg) Read(input *fuse.ReadIn, buf []byte) (fuse.ReadResult, fuse.Status) {
	var (
		ret     ReadResult
		fsINode types.FsINode
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return ret, types.ErrorToFuseStatus(err)
	}

	ret.dataLen, err = p.MemStg.NetINodeDriver.PReadWithMem(fsINode.UNetINode, buf[:input.Size], input.Offset)
	if err != nil {
		return ret, types.ErrorToFuseStatus(err)
	}

	return ret, fuse.OK
}
