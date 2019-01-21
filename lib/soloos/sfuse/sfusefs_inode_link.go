package sfuse

import (
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *SFuseFs) Link(input *fuse.LinkIn, filename string, out *fuse.EntryOut) (code fuse.Status) {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.Client.MemDirTreeStg.Link(input.Oldnodeid, input.NodeId, filename, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.setFuseEntryOutByFsINode(out, &fsINode)
	return fuse.OK
}

func (p *SFuseFs) Symlink(header *fuse.InHeader, pointedTo string, linkName string, out *fuse.EntryOut) fuse.Status {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.Client.MemDirTreeStg.Symlink(header.NodeId, pointedTo, linkName, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.setFuseEntryOutByFsINode(out, &fsINode)
	return fuse.OK
}

func (p *SFuseFs) Readlink(header *fuse.InHeader) ([]byte, fuse.Status) {
	var (
		out []byte
		err error
	)
	out, err = p.Client.MemDirTreeStg.Readlink(header.NodeId)
	if err != nil {
		return nil, types.ErrorToFuseStatus(err)
	}

	return out, fuse.OK
}
