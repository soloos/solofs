package memstg

import (
	fsapitypes "soloos/common/fsapi/types"
	"soloos/sdfs/types"
)

type ReadResult struct {
	dataLen int
}

var _ = fsapitypes.ReadResult(&ReadResult{})

func (p ReadResult) Bytes(buf []byte) ([]byte, fsapitypes.Status) {
	return buf[:p.dataLen], fsapitypes.OK
}

func (p ReadResult) Size() int {
	return p.dataLen
}

func (p ReadResult) Done() {
}

func (p *DirTreeStg) SimpleReadWithMem(uNetINode types.NetINodeUintptr,
	data []byte, offset uint64) (int, error) {
	return p.MemStg.NetINodeDriver.PReadWithMem(uNetINode, data, offset)
}

func (p *DirTreeStg) Read(input *fsapitypes.ReadIn, buf []byte) (fsapitypes.ReadResult, fsapitypes.Status) {
	var (
		ret     ReadResult
		fsINode types.FsINode
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return ret, types.ErrorToFsStatus(err)
	}

	ret.dataLen, err = p.SimpleReadWithMem(fsINode.UNetINode, buf[:input.Size], input.Offset)
	if err != nil {
		return ret, types.ErrorToFsStatus(err)
	}

	return ret, fsapitypes.OK
}
