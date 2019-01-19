package sfuse

import (
	"soloos/log"
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *SFuseFs) setFuseAttrByFsINode(fuseAttr *fuse.Attr, fsINode *types.FsINode) {
	fuseAttr.Ino = fsINode.Ino

	if fsINode.UNetINode != 0 {
		fuseAttr.Size = fsINode.UNetINode.Ptr().Size
		fuseAttr.Blocks = fsINode.UNetINode.Ptr().GetBlocks()
		fuseAttr.Blksize = uint32(fsINode.UNetINode.Ptr().MemBlockCap)
		fuseAttr.Padding = uint32(fsINode.UNetINode.Ptr().MemBlockCap)
	}
	fuseAttr.Atime = fsINode.Atime
	fuseAttr.Ctime = fsINode.Ctime
	fuseAttr.Mtime = fsINode.Mtime
	fuseAttr.Atimensec = fsINode.Atimensec
	fuseAttr.Ctimensec = fsINode.Ctimensec
	fuseAttr.Mtimensec = fsINode.Mtimensec
	fuseAttr.Mode = fsINode.Mode
	fuseAttr.Nlink = fsINode.Nlink
	fuseAttr.Uid = fsINode.Uid
	fuseAttr.Gid = fsINode.Gid
	fuseAttr.Rdev = fsINode.Rdev
}

func (p *SFuseFs) setFsINodeByFuseAttr(fsINode *types.FsINode, fuseAttr *fuse.SetAttrInCommon) {
	if fsINode.UNetINode != 0 {
		fsINode.UNetINode.Ptr().Size = fuseAttr.Size
	}
	fsINode.Atime = fuseAttr.Atime
	fsINode.Ctime = fuseAttr.Ctime
	fsINode.Mtime = fuseAttr.Mtime
	fsINode.Atimensec = fuseAttr.Atimensec
	fsINode.Ctimensec = fuseAttr.Ctimensec
	fsINode.Mtimensec = fuseAttr.Mtimensec
	fsINode.Mode = fuseAttr.Mode
	fsINode.Uid = fuseAttr.Uid
	fsINode.Gid = fuseAttr.Gid
}

func (p *SFuseFs) setFuseAttrOutByFsINode(fuseAttrOut *fuse.AttrOut, fsINode *types.FsINode) {
	fuseAttrOut.AttrValid = p.Client.DirTreeDriver.EntryAttrValid
	fuseAttrOut.AttrValidNsec = p.Client.DirTreeDriver.EntryAttrValidNsec
	p.setFuseAttrByFsINode(&fuseAttrOut.Attr, fsINode)
}

// Attributes.
func (p *SFuseFs) GetAttr(input *fuse.GetAttrIn, out *fuse.AttrOut) (code fuse.Status) {
	var (
		fsINode types.FsINode
		err     error
	)
	fsINode, err = p.Client.DirTreeDriver.GetFsINodeByID(input.NodeId)
	if err != nil {
		return fuse.EPERM
	}

	p.setFuseAttrOutByFsINode(out, &fsINode)
	log.Error(out.Ino, err)
	log.Error(out.Mode, err)
	log.Error(out.Nlink, err)
	log.Error(out.Rdev, err)

	return fuse.OK
}

func (p *SFuseFs) SetAttr(input *fuse.SetAttrIn, out *fuse.AttrOut) (code fuse.Status) {
	var (
		fsINode types.FsINode
		err     error
	)
	fsINode, err = p.Client.DirTreeDriver.GetFsINodeByID(input.NodeId)
	if err != nil {
		return fuse.EPERM
	}

	p.setFsINodeByFuseAttr(&fsINode, &input.SetAttrInCommon)
	err = p.Client.DirTreeDriver.UpdateFsINodeInDB(&fsINode)
	if err != nil {
		return fuse.EPERM
	}

	return fuse.OK
}

func (p *SFuseFs) Lookup(header *fuse.InHeader, name string, out *fuse.EntryOut) (status fuse.Status) {
	log.Error("fuck you shit")
	return fuse.EPERM
}
