package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/solofstypes"
)

type ReadResult struct {
	dataLen int
}

var _ = fsapi.ReadResult(&ReadResult{})

func (p ReadResult) Bytes(buf []byte) ([]byte, fsapi.Status) {
	return buf[:p.dataLen], fsapi.OK
}

func (p ReadResult) Size() int {
	return p.dataLen
}

func (p ReadResult) Done() {
}

func (p *PosixFs) SimpleReadWithMem(fsINodeIno solofstypes.FsINodeIno,
	data []byte, offset uint64) (int, error) {
	var (
		uFsINode solofstypes.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(fsINodeIno)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return -1, err
	}

	if uFsINode.Ptr().UNetINode == 0 {
		return -1, solofstypes.ErrObjectNotExists
	}

	return p.MemStg.NetINodeDriver.PReadWithMem(uFsINode.Ptr().UNetINode, data, offset)
}

func (p *PosixFs) Read(input *fsapi.ReadIn, buf []byte) (fsapi.ReadResult, fsapi.Status) {
	var (
		ret ReadResult
		err error
	)

	ret.dataLen, err = p.SimpleReadWithMem(input.NodeId, buf[:input.Size], input.Offset)
	if err != nil {
		return ret, ErrorToFsStatus(err)
	}

	return ret, fsapi.OK
}
