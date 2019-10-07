package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofstypes"
)

func (p *PosixFs) Link(input *fsapitypes.LinkIn, filename string, out *fsapitypes.EntryOut) fsapitypes.Status {
	var (
		srcFsINodeMeta     solofstypes.FsINodeMeta
		srcFsINodeID       = input.Oldnodeid
		newFsINodeParentID = input.NodeId
		newFsINodeMeta     solofstypes.FsINodeMeta
		err                error
	)

	err = p.FsINodeDriver.Link(srcFsINodeID, newFsINodeParentID, filename, &newFsINodeMeta)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.FetchFsINodeByIDThroughHardLink(&newFsINodeMeta, newFsINodeMeta.Ino)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.FetchFsINodeByID(&srcFsINodeMeta, srcFsINodeID)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(srcFsINodeMeta.ParentID)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &newFsINodeMeta)

	return fsapitypes.OK
}

func (p *PosixFs) Symlink(header *fsapitypes.InHeader, pointedTo string, linkName string, out *fsapitypes.EntryOut) fsapitypes.Status {
	var (
		fsINodeMeta solofstypes.FsINodeMeta
		err         error
	)
	err = p.FsINodeDriver.Symlink(header.NodeId, pointedTo, linkName, &fsINodeMeta)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(header.NodeId)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &fsINodeMeta)

	return fsapitypes.OK
}

func (p *PosixFs) Readlink(header *fsapitypes.InHeader) ([]byte, fsapitypes.Status) {
	var (
		out []byte
		err error
	)
	out, err = p.FsINodeDriver.Readlink(header.NodeId)
	if err != nil {
		return nil, ErrorToFsStatus(err)
	}

	return out, fsapitypes.OK
}
