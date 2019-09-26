package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofsapitypes"
	"soloos/solofs/solofstypes"
)

func (p *PosixFS) SimpleGetXAttrSize(fsINodeID solofsapitypes.FsINodeID, attr string) (int, fsapitypes.Status) {
	var (
		fsINodeMeta solofsapitypes.FsINodeMeta
		sz          int
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, fsINodeID)
	if err != nil {
		return 0, solofstypes.ErrorToFsStatus(err)
	}

	sz, err = p.FsINodeDriver.FIXAttrDriver.GetXAttrSize(fsINodeMeta.Ino, attr)
	if err != nil {
		return 0, solofstypes.ErrorToFsStatus(err)
	}
	return sz, fsapitypes.OK
}

func (p *PosixFS) SimpleGetXAttrData(fsINodeID solofsapitypes.FsINodeID, attr string) ([]byte, fsapitypes.Status) {
	var (
		fsINodeMeta solofsapitypes.FsINodeMeta
		data        []byte
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, fsINodeID)
	if err != nil {
		return nil, solofstypes.ErrorToFsStatus(err)
	}

	data, err = p.FsINodeDriver.FIXAttrDriver.GetXAttrData(fsINodeMeta.Ino, attr)
	if err != nil {
		return nil, solofstypes.ErrorToFsStatus(err)
	}
	return data, fsapitypes.OK
}

func (p *PosixFS) SimpleListXAttr(fsINodeID solofsapitypes.FsINodeID) ([]byte, fsapitypes.Status) {
	var (
		fsINodeMeta solofsapitypes.FsINodeMeta
		data        []byte
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, fsINodeID)
	if err != nil {
		return nil, solofstypes.ErrorToFsStatus(err)
	}

	data, err = p.FsINodeDriver.FIXAttrDriver.ListXAttr(fsINodeMeta.Ino)
	if err != nil {
		return nil, solofstypes.ErrorToFsStatus(err)
	}
	return data, fsapitypes.OK
}

func (p *PosixFS) SimpleSetXAttr(fsINodeID solofsapitypes.FsINodeID, attr string, data []byte) fsapitypes.Status {
	var err error
	err = p.FsINodeDriver.FIXAttrDriver.SetXAttr(fsINodeID, attr, data)
	return solofstypes.ErrorToFsStatus(err)
}

func (p *PosixFS) SimpleRemoveXAttr(fsINodeID solofsapitypes.FsINodeID, attr string) fsapitypes.Status {
	var err error
	err = p.FsINodeDriver.FIXAttrDriver.RemoveXAttr(fsINodeID, attr)
	return solofstypes.ErrorToFsStatus(err)
}

// Extended attributes.
func (p *PosixFS) GetXAttrSize(header *fsapitypes.InHeader, attr string) (int, fsapitypes.Status) {
	return p.SimpleGetXAttrSize(header.NodeId, attr)
}

func (p *PosixFS) GetXAttrData(header *fsapitypes.InHeader, attr string) ([]byte, fsapitypes.Status) {
	return p.SimpleGetXAttrData(header.NodeId, attr)
}

func (p *PosixFS) ListXAttr(header *fsapitypes.InHeader) ([]byte, fsapitypes.Status) {
	return p.SimpleListXAttr(header.NodeId)
}

func (p *PosixFS) SetXAttr(input *fsapitypes.SetXAttrIn, attr string, data []byte) fsapitypes.Status {
	return p.SimpleSetXAttr(input.NodeId, attr, data)
}

func (p *PosixFS) RemoveXAttr(header *fsapitypes.InHeader, attr string) fsapitypes.Status {
	return p.SimpleRemoveXAttr(header.NodeId, attr)
}
