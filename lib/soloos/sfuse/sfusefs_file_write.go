package sfuse

import "github.com/hanwen/go-fuse/fuse"

func (p *SFuseFs) Write(input *fuse.WriteIn, data []byte) (written uint32, code fuse.Status) {
	return 0, fuse.EPERM
}
