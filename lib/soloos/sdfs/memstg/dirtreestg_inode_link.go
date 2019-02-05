package memstg

import (
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *DirTreeStg) Link(input *fuse.LinkIn, filename string, out *fuse.EntryOut) fuse.Status {
	var (
		srcFsINode types.FsINode
		fsINode    types.FsINode
		err        error
	)

	err = p.FetchFsINodeByID(input.Oldnodeid, &srcFsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.FsINodeDriver.Link(&srcFsINode, input.NodeId, filename, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.FetchFsINodeByIDThroughHardLink(fsINode.Ino, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(srcFsINode.ParentID)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.SetFuseEntryOutByFsINode(out, &fsINode)

	return fuse.OK
}

func (p *DirTreeStg) Symlink(header *fuse.InHeader, pointedTo string, linkName string, out *fuse.EntryOut) fuse.Status {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.FsINodeDriver.Symlink(header.NodeId, pointedTo, linkName, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(header.NodeId)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.SetFuseEntryOutByFsINode(out, &fsINode)

	return fuse.OK
}

func (p *DirTreeStg) Readlink(header *fuse.InHeader) ([]byte, fuse.Status) {
	var (
		out []byte
		err error
	)
	out, err = p.FsINodeDriver.Readlink(header.NodeId)
	if err != nil {
		return nil, types.ErrorToFuseStatus(err)
	}

	return out, fuse.OK
}
