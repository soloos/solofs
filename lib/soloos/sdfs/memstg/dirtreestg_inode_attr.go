package memstg

import (
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *DirTreeStg) SetFuseAttrByFsINode(fuseAttr *fuse.Attr, fsINode *types.FsINode) {
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
	fuseAttr.Nlink = uint32(fsINode.Nlink)
	fuseAttr.Uid = fsINode.Uid
	fuseAttr.Gid = fsINode.Gid
	fuseAttr.Rdev = fsINode.Rdev
}

func (p *DirTreeStg) SetFuseEntryOutByFsINode(fuseEntryOut *fuse.EntryOut, pFsINode *types.FsINode) {
	fuseEntryOut.NodeId = pFsINode.Ino
	fuseEntryOut.EntryValid = p.EntryAttrValid
	fuseEntryOut.EntryValidNsec = p.EntryAttrValidNsec
	fuseEntryOut.AttrValid = p.EntryAttrValid
	fuseEntryOut.AttrValidNsec = p.EntryAttrValidNsec
	p.SetFuseAttrByFsINode(&fuseEntryOut.Attr, pFsINode)
}

func (p *DirTreeStg) SetFsINodeByFuseAttr(fsINode *types.FsINode, input *fuse.SetAttrIn) fuse.Status {
	if input.Valid&fuse.FATTR_MODE != 0 {
		// fsINode.Mode = uint32(0777)&input.Mode | uint32(types.FsINodeTypeToFuseType(fsINode.Type))
		fsINode.Mode = input.Mode
	}

	if input.Valid&fuse.FATTR_UID != 0 {
		fsINode.Uid = input.Uid
	}
	if input.Valid&fuse.FATTR_GID != 0 {
		fsINode.Gid = input.Gid
	}

	if input.Valid&fuse.FATTR_SIZE != 0 {
		if fsINode.Type == types.FSINODE_TYPE_DIR {
			return fuse.EISDIR
		}
		p.MemStg.NetINodeDriver.NetINodeTruncate(fsINode.UNetINode, input.Size)
	}

	now := p.FsINodeDriver.Timer.Now()
	nowt := types.DirTreeTime(now.Unix())
	nowtnsec := types.DirTreeTimeNsec(now.UnixNano())

	if input.Valid&(fuse.FATTR_ATIME|fuse.FATTR_MTIME|fuse.FATTR_ATIME_NOW|fuse.FATTR_MTIME_NOW) != 0 {
		if input.Valid&fuse.FATTR_ATIME != 0 {
			if input.Valid&fuse.FATTR_ATIME_NOW != 0 {
				fsINode.Atime = nowt
				fsINode.Atimensec = nowtnsec
			} else {
				fsINode.Atime = input.Atime
				fsINode.Atimensec = input.Atimensec
			}
		}

		if input.Valid&fuse.FATTR_MTIME != 0 {
			if input.Valid&fuse.FATTR_MTIME_NOW != 0 {
				fsINode.Mtime = nowt
				fsINode.Mtimensec = nowtnsec
			} else {
				fsINode.Mtime = input.Mtime
				fsINode.Mtimensec = input.Mtimensec
			}
		}
	}

	fsINode.Ctime = nowt
	fsINode.Ctimensec = nowtnsec

	return fuse.OK
}

func (p *DirTreeStg) SetFuseAttrOutByFsINode(fuseAttrOut *fuse.AttrOut, pFsINode *types.FsINode) {
	fuseAttrOut.AttrValid = p.EntryAttrValid
	fuseAttrOut.AttrValidNsec = p.EntryAttrValidNsec
	p.SetFuseAttrByFsINode(&fuseAttrOut.Attr, pFsINode)
}

func (p *DirTreeStg) GetAttr(input *fuse.GetAttrIn, out *fuse.AttrOut) fuse.Status {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.SetFuseAttrOutByFsINode(out, &fsINode)

	return fuse.OK
}

func (p *DirTreeStg) SetAttr(input *fuse.SetAttrIn, out *fuse.AttrOut) fuse.Status {
	var (
		fsINode types.FsINode
		code    fuse.Status
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	code = p.SetFsINodeByFuseAttr(&fsINode, input)
	if code != fuse.OK {
		return code
	}

	err = p.UpdateFsINodeInDB(&fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.SetFuseAttrOutByFsINode(out, &fsINode)

	return fuse.OK
}
