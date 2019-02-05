package memstg

import (
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *DirTreeStg) SimpleGetXAttrSize(fsINodeID types.FsINodeID, attr string) (int, fuse.Status) {
	var (
		fsINode types.FsINode
		sz      int
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(fsINodeID, &fsINode)
	if err != nil {
		return 0, types.ErrorToFuseStatus(err)
	}

	sz, err = p.FsINodeDriver.FIXAttrDriver.GetXAttrSize(fsINode.Ino, attr)
	if err != nil {
		return 0, types.ErrorToFuseStatus(err)
	}
	return sz, fuse.OK
}

func (p *DirTreeStg) SimpleGetXAttrData(fsINodeID types.FsINodeID, attr string) ([]byte, fuse.Status) {
	var (
		fsINode types.FsINode
		data    []byte
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(fsINodeID, &fsINode)
	if err != nil {
		return nil, types.ErrorToFuseStatus(err)
	}

	data, err = p.FsINodeDriver.FIXAttrDriver.GetXAttrData(fsINode.Ino, attr)
	if err != nil {
		return nil, types.ErrorToFuseStatus(err)
	}
	return data, fuse.OK
}

func (p *DirTreeStg) SimpleListXAttr(fsINodeID types.FsINodeID) ([]byte, fuse.Status) {
	var (
		fsINode types.FsINode
		data    []byte
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(fsINodeID, &fsINode)
	if err != nil {
		return nil, types.ErrorToFuseStatus(err)
	}

	data, err = p.FsINodeDriver.FIXAttrDriver.ListXAttr(fsINode.Ino)
	if err != nil {
		return nil, types.ErrorToFuseStatus(err)
	}
	return data, fuse.OK
}

func (p *DirTreeStg) SimpleSetXAttr(fsINodeID types.FsINodeID, attr string, data []byte) fuse.Status {
	var err error
	err = p.FsINodeDriver.FIXAttrDriver.SetXAttr(fsINodeID, attr, data)
	return types.ErrorToFuseStatus(err)
}

func (p *DirTreeStg) SimpleRemoveXAttr(fsINodeID types.FsINodeID, attr string) fuse.Status {
	var err error
	err = p.FsINodeDriver.FIXAttrDriver.RemoveXAttr(fsINodeID, attr)
	return types.ErrorToFuseStatus(err)
}

// Extended attributes.
func (p *DirTreeStg) GetXAttrSize(header *fuse.InHeader, attr string) (int, fuse.Status) {
	return p.SimpleGetXAttrSize(header.NodeId, attr)
}

func (p *DirTreeStg) GetXAttrData(header *fuse.InHeader, attr string) ([]byte, fuse.Status) {
	return p.SimpleGetXAttrData(header.NodeId, attr)
}

func (p *DirTreeStg) ListXAttr(header *fuse.InHeader) ([]byte, fuse.Status) {
	return p.SimpleListXAttr(header.NodeId)
}

func (p *DirTreeStg) SetXAttr(input *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
	return p.SimpleSetXAttr(input.NodeId, attr, data)
}

func (p *DirTreeStg) RemoveXAttr(header *fuse.InHeader, attr string) fuse.Status {
	return p.SimpleRemoveXAttr(header.NodeId, attr)
}
