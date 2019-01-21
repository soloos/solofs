package sfuse

import "github.com/hanwen/go-fuse/fuse"

type ReadResult struct {
}

var _ = fuse.ReadResult(&ReadResult{})

func (p ReadResult) Bytes(buf []byte) ([]byte, fuse.Status) {
	return nil, fuse.EPERM
}

func (p ReadResult) Size() int {
	return 0
}

func (p ReadResult) Done() {
}

func (p *SFuseFs) Read(input *fuse.ReadIn, buf []byte) (fuse.ReadResult, fuse.Status) {
	var ret ReadResult
	return ret, fuse.EPERM
}
