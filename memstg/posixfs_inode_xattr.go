package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofsapitypes"
	"soloos/solofs/solofstypes"
)

// Extended attributes.
func (p *PosixFs) GetXAttrSize(header *fsapitypes.InHeader, attr string) (int, fsapitypes.Status) {
	var (
		fsINodeID   = header.NodeId
		fsINodeMeta solofsapitypes.FsINodeMeta
		sz          int
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, fsINodeID)
	if err != nil {
		return 0, solofstypes.ErrorToFsStatus(err)
	}

	sz, err = p.FIXAttrDriver.GetXAttrSize(fsINodeMeta.Ino, attr)
	if err != nil {
		return 0, solofstypes.ErrorToFsStatus(err)
	}
	return sz, fsapitypes.OK
}

func (p *PosixFs) GetXAttrData(header *fsapitypes.InHeader, attr string) ([]byte, fsapitypes.Status) {
	var (
		fsINodeID   = header.NodeId
		fsINodeMeta solofsapitypes.FsINodeMeta
		data        []byte
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, fsINodeID)
	if err != nil {
		return nil, solofstypes.ErrorToFsStatus(err)
	}

	data, err = p.FIXAttrDriver.GetXAttrData(fsINodeMeta.Ino, attr)
	if err != nil {
		return nil, solofstypes.ErrorToFsStatus(err)
	}
	return data, fsapitypes.OK
}

func (p *PosixFs) ListXAttr(header *fsapitypes.InHeader) ([]byte, fsapitypes.Status) {
	var (
		fsINodeID   = header.NodeId
		fsINodeMeta solofsapitypes.FsINodeMeta
		data        []byte
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, fsINodeID)
	if err != nil {
		return nil, solofstypes.ErrorToFsStatus(err)
	}

	data, err = p.FIXAttrDriver.ListXAttr(fsINodeMeta.Ino)
	if err != nil {
		return nil, solofstypes.ErrorToFsStatus(err)
	}
	return data, fsapitypes.OK
}

func (p *PosixFs) SetXAttr(input *fsapitypes.SetXAttrIn, attr string, data []byte) fsapitypes.Status {
	var fsINodeID = input.NodeId
	var err error
	err = p.FIXAttrDriver.SetXAttr(fsINodeID, attr, data)
	return solofstypes.ErrorToFsStatus(err)
}

func (p *PosixFs) RemoveXAttr(header *fsapitypes.InHeader, attr string) fsapitypes.Status {
	var fsINodeID = header.NodeId
	var err error
	err = p.FIXAttrDriver.RemoveXAttr(fsINodeID, attr)
	return solofstypes.ErrorToFsStatus(err)
}
