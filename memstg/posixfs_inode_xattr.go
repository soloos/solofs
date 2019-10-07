package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/solofstypes"
)

// Extended attributes.
func (p *PosixFs) GetXAttrSize(header *fsapi.InHeader, attr string) (int, fsapi.Status) {
	var (
		fsINodeID   = header.NodeId
		fsINodeMeta solofstypes.FsINodeMeta
		sz          int
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, fsINodeID)
	if err != nil {
		return 0, ErrorToFsStatus(err)
	}

	sz, err = p.FIXAttrDriver.GetXAttrSize(fsINodeMeta.Ino, attr)
	if err != nil {
		return 0, ErrorToFsStatus(err)
	}
	return sz, fsapi.OK
}

func (p *PosixFs) GetXAttrData(header *fsapi.InHeader, attr string) ([]byte, fsapi.Status) {
	var (
		fsINodeID   = header.NodeId
		fsINodeMeta solofstypes.FsINodeMeta
		data        []byte
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, fsINodeID)
	if err != nil {
		return nil, ErrorToFsStatus(err)
	}

	data, err = p.FIXAttrDriver.GetXAttrData(fsINodeMeta.Ino, attr)
	if err != nil {
		return nil, ErrorToFsStatus(err)
	}
	return data, fsapi.OK
}

func (p *PosixFs) ListXAttr(header *fsapi.InHeader) ([]byte, fsapi.Status) {
	var (
		fsINodeID   = header.NodeId
		fsINodeMeta solofstypes.FsINodeMeta
		data        []byte
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, fsINodeID)
	if err != nil {
		return nil, ErrorToFsStatus(err)
	}

	data, err = p.FIXAttrDriver.ListXAttr(fsINodeMeta.Ino)
	if err != nil {
		return nil, ErrorToFsStatus(err)
	}
	return data, fsapi.OK
}

func (p *PosixFs) SetXAttr(input *fsapi.SetXAttrIn, attr string, data []byte) fsapi.Status {
	var fsINodeID = input.NodeId
	var err error
	err = p.FIXAttrDriver.SetXAttr(fsINodeID, attr, data)
	return ErrorToFsStatus(err)
}

func (p *PosixFs) RemoveXAttr(header *fsapi.InHeader, attr string) fsapi.Status {
	var fsINodeID = header.NodeId
	var err error
	err = p.FIXAttrDriver.RemoveXAttr(fsINodeID, attr)
	return ErrorToFsStatus(err)
}
