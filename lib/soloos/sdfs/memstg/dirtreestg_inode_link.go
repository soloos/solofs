package memstg

import (
	fsapitypes "soloos/fsapi/types"
	"soloos/sdfs/types"
)

func (p *DirTreeStg) Link(input *fsapitypes.LinkIn, filename string, out *fsapitypes.EntryOut) fsapitypes.Status {
	var (
		srcFsINode types.FsINode
		fsINode    types.FsINode
		err        error
	)

	err = p.FetchFsINodeByID(input.Oldnodeid, &srcFsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.FsINodeDriver.Link(&srcFsINode, input.NodeId, filename, &fsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.FetchFsINodeByIDThroughHardLink(fsINode.Ino, &fsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(srcFsINode.ParentID)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &fsINode)

	return fsapitypes.OK
}

func (p *DirTreeStg) Symlink(header *fsapitypes.InHeader, pointedTo string, linkName string, out *fsapitypes.EntryOut) fsapitypes.Status {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.FsINodeDriver.Symlink(header.NodeId, pointedTo, linkName, &fsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(header.NodeId)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &fsINode)

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
