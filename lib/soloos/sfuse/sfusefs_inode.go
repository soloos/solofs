package sfuse

import (
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *SFuseFs) setFuseAttrByFsINode(fuseAttr *fuse.Attr, pFsINode *types.FsINode) {
	fuseAttr.Ino = pFsINode.Ino

	if pFsINode.UNetINode != 0 {
		fuseAttr.Size = pFsINode.UNetINode.Ptr().Size
		fuseAttr.Blocks = pFsINode.UNetINode.Ptr().GetBlocks()
		fuseAttr.Blksize = uint32(pFsINode.UNetINode.Ptr().MemBlockCap)
		fuseAttr.Padding = uint32(pFsINode.UNetINode.Ptr().MemBlockCap)
	}

	fuseAttr.Atime = pFsINode.Atime
	fuseAttr.Ctime = pFsINode.Ctime
	fuseAttr.Mtime = pFsINode.Mtime
	fuseAttr.Atimensec = pFsINode.Atimensec
	fuseAttr.Ctimensec = pFsINode.Ctimensec
	fuseAttr.Mtimensec = pFsINode.Mtimensec
	fuseAttr.Mode = pFsINode.Mode
	fuseAttr.Nlink = pFsINode.Nlink
	fuseAttr.Uid = pFsINode.Uid
	fuseAttr.Gid = pFsINode.Gid
	fuseAttr.Rdev = pFsINode.Rdev
}

func (p *SFuseFs) setFuseEntryOutByFsINode(fuseEntryOut *fuse.EntryOut, pFsINode *types.FsINode) {
	fuseEntryOut.NodeId = pFsINode.Ino
	fuseEntryOut.EntryValid = p.Client.MemDirTreeStg.EntryAttrValid
	fuseEntryOut.EntryValidNsec = p.Client.MemDirTreeStg.EntryAttrValidNsec
	fuseEntryOut.AttrValid = p.Client.MemDirTreeStg.EntryAttrValid
	fuseEntryOut.AttrValidNsec = p.Client.MemDirTreeStg.EntryAttrValidNsec
	p.setFuseAttrByFsINode(&fuseEntryOut.Attr, pFsINode)
}

func (p *SFuseFs) setFsINodeByFuseAttr(pFsINode *types.FsINode, fuseAttr *fuse.SetAttrInCommon) {
	if pFsINode.UNetINode != 0 {
		pFsINode.UNetINode.Ptr().Size = fuseAttr.Size
	}
	pFsINode.Atime = fuseAttr.Atime
	pFsINode.Ctime = fuseAttr.Ctime
	pFsINode.Mtime = fuseAttr.Mtime
	pFsINode.Atimensec = fuseAttr.Atimensec
	pFsINode.Ctimensec = fuseAttr.Ctimensec
	pFsINode.Mtimensec = fuseAttr.Mtimensec
	pFsINode.Mode = fuseAttr.Mode
	pFsINode.Uid = fuseAttr.Uid
	pFsINode.Gid = fuseAttr.Gid
}

func (p *SFuseFs) setFuseAttrOutByFsINode(fuseAttrOut *fuse.AttrOut, pFsINode *types.FsINode) {
	fuseAttrOut.AttrValid = p.Client.MemDirTreeStg.EntryAttrValid
	fuseAttrOut.AttrValidNsec = p.Client.MemDirTreeStg.EntryAttrValidNsec
	p.setFuseAttrByFsINode(&fuseAttrOut.Attr, pFsINode)
}

func (p *SFuseFs) FetchFsINodeByName(parentFsINodeID types.FsINodeID, name string, fsINode *types.FsINode) error {
	return p.Client.MemDirTreeStg.FetchFsINodeByName(parentFsINodeID, name, fsINode)
}

func (p *SFuseFs) FetchFsINodeByID(fsINodeID types.FsINodeID, fsINode *types.FsINode) error {
	return p.Client.MemDirTreeStg.FetchFsINodeByID(fsINodeID, fsINode)
}

// Attributes.
func (p *SFuseFs) GetAttr(input *fuse.GetAttrIn, out *fuse.AttrOut) (code fuse.Status) {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.FetchFsINodeByID(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.setFuseAttrOutByFsINode(out, &fsINode)

	return fuse.OK
}

func (p *SFuseFs) SetAttr(input *fuse.SetAttrIn, out *fuse.AttrOut) (code fuse.Status) {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.FetchFsINodeByID(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.setFsINodeByFuseAttr(&fsINode, &input.SetAttrInCommon)
	err = p.Client.MemDirTreeStg.UpdateFsINodeInDB(&fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	return fuse.OK
}

func (p *SFuseFs) Lookup(header *fuse.InHeader, name string, out *fuse.EntryOut) (status fuse.Status) {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.FetchFsINodeByName(header.NodeId, name, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.setFuseEntryOutByFsINode(out, &fsINode)
	return fuse.OK
}

// Modifying structure.
func (p *SFuseFs) Mknod(input *fuse.MknodIn, name string, out *fuse.EntryOut) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Unlink(header *fuse.InHeader, name string) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Rename(input *fuse.RenameIn, oldName string, newName string) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Access(input *fuse.AccessIn) (code fuse.Status) {
	return fuse.OK
}

func (p *SFuseFs) Forget(nodeid, nlookup uint64) {
}

// File handling.
func (p *SFuseFs) Create(input *fuse.CreateIn, name string, out *fuse.CreateOut) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Open(input *fuse.OpenIn, out *fuse.OpenOut) (status fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Release(input *fuse.ReleaseIn) {
}

func (p *SFuseFs) Flush(input *fuse.FlushIn) fuse.Status {
	return fuse.EPERM
}

func (p *SFuseFs) Fsync(input *fuse.FsyncIn) (code fuse.Status) {
	return fuse.EPERM
}

func (p *SFuseFs) Fallocate(input *fuse.FallocateIn) (code fuse.Status) {
	return fuse.EPERM
}
