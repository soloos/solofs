package memstg

import (
	fsapitypes "soloos/common/fsapi/types"
	"soloos/sdfs/types"
)

func (p *DirTreeStg) Link(input *fsapitypes.LinkIn, filename string, out *fsapitypes.EntryOut) fsapitypes.Status {
	var (
		srcFsINodeMeta types.FsINodeMeta
		newFsINodeMeta types.FsINodeMeta
		err            error
	)

	err = p.FetchFsINodeByID(&srcFsINodeMeta, input.Oldnodeid)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.FsINodeDriver.Link(&srcFsINodeMeta, input.NodeId, filename, &newFsINodeMeta)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.FetchFsINodeByIDThroughHardLink(&newFsINodeMeta, newFsINodeMeta.Ino)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(srcFsINodeMeta.ParentID)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &newFsINodeMeta)

	return fsapitypes.OK
}

func (p *DirTreeStg) Symlink(header *fsapitypes.InHeader, pointedTo string, linkName string, out *fsapitypes.EntryOut) fsapitypes.Status {
	var (
		fsINodeMeta types.FsINodeMeta
		err         error
	)
	err = p.FsINodeDriver.Symlink(header.NodeId, pointedTo, linkName, &fsINodeMeta)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(header.NodeId)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &fsINodeMeta)

	return fsapitypes.OK
}

func (p *DirTreeStg) Readlink(header *fsapitypes.InHeader) ([]byte, fsapitypes.Status) {
	var (
		out []byte
		err error
	)
	out, err = p.FsINodeDriver.Readlink(header.NodeId)
	if err != nil {
		return nil, types.ErrorToFsStatus(err)
	}

	return out, fsapitypes.OK
}
