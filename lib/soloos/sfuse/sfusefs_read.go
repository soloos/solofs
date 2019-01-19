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
