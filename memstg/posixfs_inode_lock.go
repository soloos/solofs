package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofsapitypes"
	"soloos/solofs/solofstypes"
)

func (p *PosixFS) setLKOutByMeta(out *fsapitypes.LkOut, meta *solofsapitypes.INodeRWMutexMeta) {
	out.Lk.Start = meta.Start
	out.Lk.End = meta.End
	out.Lk.Typ = meta.Typ
	out.Lk.Pid = meta.Pid
}

func (p *PosixFS) setMetaByLKIn(lkIn *fsapitypes.LkIn, meta *solofsapitypes.INodeRWMutexMeta) {
	meta.Start = lkIn.Lk.Start
	meta.End = lkIn.Lk.End
	meta.Typ = lkIn.Lk.Typ
	meta.Pid = lkIn.Lk.Pid
}

// GetLk returns existing lock information for file
func (p *PosixFS) GetLk(input *fsapitypes.LkIn, out *fsapitypes.LkOut) fsapitypes.Status {
	var (
		meta solofsapitypes.INodeRWMutexMeta
		err  error
	)
	err = p.FsINodeDriver.GetLk(input.NodeId, &meta)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}
	p.setLKOutByMeta(out, &meta)
	return fsapitypes.OK
}

// SetLk Sets or clears the lock described by lk on file.
func (p *PosixFS) SetLk(input *fsapitypes.LkIn) fsapitypes.Status {
	var (
		meta solofsapitypes.INodeRWMutexMeta
		err  error
	)
	p.setMetaByLKIn(input, &meta)
	err = p.FsINodeDriver.SetLk(input.NodeId, &meta)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}
	return fsapitypes.OK
}

// SetLkw Sets or clears the lock described by lk. This call blocks until the operation can be completed.
func (p *PosixFS) SetLkw(input *fsapitypes.LkIn) fsapitypes.Status {
	var (
		meta solofsapitypes.INodeRWMutexMeta
		err  error
	)
	p.setMetaByLKIn(input, &meta)
	err = p.FsINodeDriver.SetLkw(input.NodeId, &meta)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}
	return fsapitypes.OK
}
