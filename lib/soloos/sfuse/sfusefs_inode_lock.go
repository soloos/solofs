package sfuse

import (
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

// GetLk returns existing lock information for file
func (p *SFuseFs) GetLk(input *fuse.LkIn, out *fuse.LkOut) (code fuse.Status) {
	var err error
	err = p.Client.MemDirTreeStg.GetLk(input, out)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	return fuse.OK
}

// SetLk Sets or clears the lock described by lk on file.
func (p *SFuseFs) SetLk(input *fuse.LkIn) (code fuse.Status) {
	var err error
	err = p.Client.MemDirTreeStg.SetLk(input)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	return fuse.OK
}

// SetLkw Sets or clears the lock described by lk. This call blocks until the operation can be completed.
func (p *SFuseFs) SetLkw(input *fuse.LkIn) (code fuse.Status) {
	var err error
	err = p.Client.MemDirTreeStg.SetLk(input)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	return fuse.OK
}
