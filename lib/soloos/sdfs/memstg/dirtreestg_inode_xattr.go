package memstg

import (
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *DirTreeStg) GetXAttrSize(fsINodeID types.FsINodeID, attr string) (int, fuse.Status) {
	var (
		sz  int
		err error
	)
	sz, err = p.FsINodeDriver.FIXAttrDriver.GetXAttrSize(fsINodeID, attr)
	if err != nil {
		return 0, types.ErrorToFuseStatus(err)
	}
	return sz, fuse.OK
}

func (p *DirTreeStg) GetXAttrData(fsINodeID types.FsINodeID, attr string) ([]byte, fuse.Status) {
	var (
		data []byte
		err  error
	)
	data, err = p.FsINodeDriver.FIXAttrDriver.GetXAttrData(fsINodeID, attr)
	if err != nil {
		return nil, types.ErrorToFuseStatus(err)
	}
	return data, fuse.OK
}

func (p *DirTreeStg) ListXAttr(fsINodeID types.FsINodeID) ([]byte, fuse.Status) {
	var (
		data []byte
		err  error
	)
	data, err = p.FsINodeDriver.FIXAttrDriver.ListXAttr(fsINodeID)
	if err != nil {
		return nil, types.ErrorToFuseStatus(err)
	}
	return data, fuse.OK
}

func (p *DirTreeStg) SetXAttr(fsINodeID types.FsINodeID, attr string, data []byte) fuse.Status {
	var err error
	err = p.FsINodeDriver.FIXAttrDriver.SetXAttr(fsINodeID, attr, data)
	return types.ErrorToFuseStatus(err)
}

func (p *DirTreeStg) RemoveXAttr(fsINodeID types.FsINodeID, attr string) fuse.Status {
	var err error
	err = p.FsINodeDriver.FIXAttrDriver.RemoveXAttr(fsINodeID, attr)
	return types.ErrorToFuseStatus(err)
}
