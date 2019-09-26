package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofsapitypes"
	"soloos/solofs/solofstypes"
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

func (p *PosixFS) SimpleReadWithMem(fsINodeID solofsapitypes.FsINodeID,
	data []byte, offset uint64) (int, error) {
	var (
		uFsINode solofsapitypes.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(fsINodeID)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return -1, err
	}

	if uFsINode.Ptr().UNetINode == 0 {
		return -1, solofsapitypes.ErrObjectNotExists
	}

	return p.MemStg.NetINodeDriver.PReadWithMem(uFsINode.Ptr().UNetINode, data, offset)
}

func (p *PosixFS) Read(input *fsapitypes.ReadIn, buf []byte) (fsapitypes.ReadResult, fsapitypes.Status) {
	var (
		ret ReadResult
		err error
	)

	ret.dataLen, err = p.SimpleReadWithMem(input.NodeId, buf[:input.Size], input.Offset)
	if err != nil {
		return ret, solofstypes.ErrorToFsStatus(err)
	}

	return ret, fsapitypes.OK
}
