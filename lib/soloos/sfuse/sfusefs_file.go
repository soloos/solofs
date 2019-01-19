package sfuse

import "github.com/hanwen/go-fuse/fuse"

func (p *SFuseFs) Read(input *fuse.ReadIn, buf []byte) (fuse.ReadResult, fuse.Status) {
	var ret ReadResult
	return ret, fuse.EPERM
}

func (p *SFuseFs) Write(input *fuse.WriteIn, data []byte) (written uint32, code fuse.Status) {
	return 0, fuse.EPERM
}
