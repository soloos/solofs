package memstg

import (
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *DirTreeStg) setLKOutByMeta(out *fuse.LkOut, meta *types.INodeRWMutexMeta) {
	out.Lk.Start = meta.Start
	out.Lk.End = meta.End
	out.Lk.Typ = meta.Typ
	out.Lk.Pid = meta.Pid
}

func (p *DirTreeStg) setMetaByLKIn(lkIn *fuse.LkIn, meta *types.INodeRWMutexMeta) {
	meta.Start = lkIn.Lk.Start
	meta.End = lkIn.Lk.End
	meta.Typ = lkIn.Lk.Typ
	meta.Pid = lkIn.Lk.Pid
}

// GetLk returns existing lock information for file
func (p *DirTreeStg) GetLk(input *fuse.LkIn, out *fuse.LkOut) error {
	var (
		meta types.INodeRWMutexMeta
		err  error
	)
	err = p.FsINodeDriver.GetLk(input.NodeId, &meta)
	if err != nil {
		return err
	}
	p.setLKOutByMeta(out, &meta)
	return nil
}

// SetLk Sets or clears the lock described by lk on file.
func (p *DirTreeStg) SetLk(input *fuse.LkIn) error {
	var (
		meta types.INodeRWMutexMeta
		err  error
	)
	p.setMetaByLKIn(input, &meta)
	err = p.FsINodeDriver.SetLk(input.NodeId, &meta)
	if err != nil {
		return err
	}
	return nil
}

// SetLkw Sets or clears the lock described by lk. This call blocks until the operation can be completed.
func (p *DirTreeStg) SetLkw(input *fuse.LkIn) error {
	var (
		meta types.INodeRWMutexMeta
		err  error
	)
	p.setMetaByLKIn(input, &meta)
	err = p.FsINodeDriver.SetLkw(input.NodeId, &meta)
	if err != nil {
		return err
	}
	return nil
}
