package memstg

import (
	fsapitypes "soloos/fsapi/types"
	"soloos/sdfs/types"
)

func (p *DirTreeStg) SimpleGetXAttrSize(fsINodeID types.FsINodeID, attr string) (int, fsapitypes.Status) {
	var (
		fsINode types.FsINode
		sz      int
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(fsINodeID, &fsINode)
	if err != nil {
		return 0, types.ErrorToFsStatus(err)
	}

	sz, err = p.FsINodeDriver.FIXAttrDriver.GetXAttrSize(fsINode.Ino, attr)
	if err != nil {
		return 0, types.ErrorToFsStatus(err)
	}
	return sz, fsapitypes.OK
}

func (p *DirTreeStg) SimpleGetXAttrData(fsINodeID types.FsINodeID, attr string) ([]byte, fsapitypes.Status) {
	var (
		fsINode types.FsINode
		data    []byte
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(fsINodeID, &fsINode)
	if err != nil {
		return nil, types.ErrorToFsStatus(err)
	}

	data, err = p.FsINodeDriver.FIXAttrDriver.GetXAttrData(fsINode.Ino, attr)
	if err != nil {
		return nil, types.ErrorToFsStatus(err)
	}
	return data, fsapitypes.OK
}

func (p *DirTreeStg) SimpleListXAttr(fsINodeID types.FsINodeID) ([]byte, fsapitypes.Status) {
	var (
		fsINode types.FsINode
		data    []byte
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(fsINodeID, &fsINode)
	if err != nil {
		return nil, types.ErrorToFsStatus(err)
	}

	data, err = p.FsINodeDriver.FIXAttrDriver.ListXAttr(fsINode.Ino)
	if err != nil {
		return nil, types.ErrorToFsStatus(err)
	}
	return data, fsapitypes.OK
}

func (p *DirTreeStg) SimpleSetXAttr(fsINodeID types.FsINodeID, attr string, data []byte) fsapitypes.Status {
	var err error
	err = p.FsINodeDriver.FIXAttrDriver.SetXAttr(fsINodeID, attr, data)
	return types.ErrorToFsStatus(err)
}

func (p *DirTreeStg) SimpleRemoveXAttr(fsINodeID types.FsINodeID, attr string) fsapitypes.Status {
	var err error
	err = p.FsINodeDriver.FIXAttrDriver.RemoveXAttr(fsINodeID, attr)
	return types.ErrorToFsStatus(err)
}

// Extended attributes.
func (p *DirTreeStg) GetXAttrSize(header *fsapitypes.InHeader, attr string) (int, fsapitypes.Status) {
	return p.SimpleGetXAttrSize(header.NodeId, attr)
}

func (p *DirTreeStg) GetXAttrData(header *fsapitypes.InHeader, attr string) ([]byte, fsapitypes.Status) {
	return p.SimpleGetXAttrData(header.NodeId, attr)
}

func (p *DirTreeStg) ListXAttr(header *fsapitypes.InHeader) ([]byte, fsapitypes.Status) {
	return p.SimpleListXAttr(header.NodeId)
}

func (p *DirTreeStg) SetXAttr(input *fsapitypes.SetXAttrIn, attr string, data []byte) fsapitypes.Status {
	return p.SimpleSetXAttr(input.NodeId, attr, data)
}

func (p *DirTreeStg) RemoveXAttr(header *fsapitypes.InHeader, attr string) fsapitypes.Status {
	return p.SimpleRemoveXAttr(header.NodeId, attr)
}
