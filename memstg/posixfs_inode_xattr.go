package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofsapitypes"
	"soloos/solofs/solofstypes"
)

func (p *PosixFs) SimpleGetXAttrSize(fsINodeID solofsapitypes.FsINodeID, attr string) (int, fsapitypes.Status) {
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

func (p *PosixFs) SimpleGetXAttrData(fsINodeID solofsapitypes.FsINodeID, attr string) ([]byte, fsapitypes.Status) {
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

func (p *PosixFs) SimpleListXAttr(fsINodeID solofsapitypes.FsINodeID) ([]byte, fsapitypes.Status) {
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

func (p *PosixFs) SimpleSetXAttr(fsINodeID solofsapitypes.FsINodeID, attr string, data []byte) fsapitypes.Status {
	var err error
	err = p.FsINodeDriver.FIXAttrDriver.SetXAttr(fsINodeID, attr, data)
	return solofstypes.ErrorToFsStatus(err)
}

func (p *PosixFs) SimpleRemoveXAttr(fsINodeID solofsapitypes.FsINodeID, attr string) fsapitypes.Status {
	var err error
	err = p.FsINodeDriver.FIXAttrDriver.RemoveXAttr(fsINodeID, attr)
	return solofstypes.ErrorToFsStatus(err)
}

// Extended attributes.
func (p *PosixFs) GetXAttrSize(header *fsapitypes.InHeader, attr string) (int, fsapitypes.Status) {
	return p.SimpleGetXAttrSize(header.NodeId, attr)
}

func (p *PosixFs) GetXAttrData(header *fsapitypes.InHeader, attr string) ([]byte, fsapitypes.Status) {
	return p.SimpleGetXAttrData(header.NodeId, attr)
}

func (p *PosixFs) ListXAttr(header *fsapitypes.InHeader) ([]byte, fsapitypes.Status) {
	return p.SimpleListXAttr(header.NodeId)
}

func (p *PosixFs) SetXAttr(input *fsapitypes.SetXAttrIn, attr string, data []byte) fsapitypes.Status {
	return p.SimpleSetXAttr(input.NodeId, attr, data)
}

func (p *PosixFs) RemoveXAttr(header *fsapitypes.InHeader, attr string) fsapitypes.Status {
	return p.SimpleRemoveXAttr(header.NodeId, attr)
}
