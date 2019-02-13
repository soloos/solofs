package memstg

import (
	fsapitypes "soloos/common/fsapi/types"
	"soloos/sdfs/types"
)

func (p *DirTreeStg) SetFsAttrByFsINode(fsAttr *fsapitypes.Attr, fsINode *types.FsINode) {
	fsAttr.Ino = fsINode.Ino

	if fsINode.UNetINode != 0 {
		fsAttr.Size = fsINode.UNetINode.Ptr().Size
		fsAttr.Blocks = fsINode.UNetINode.Ptr().GetBlocks()
		fsAttr.Blksize = uint32(fsINode.UNetINode.Ptr().MemBlockCap)
		fsAttr.Padding = uint32(fsINode.UNetINode.Ptr().MemBlockCap)
	}

	fsAttr.Atime = fsINode.Atime
	fsAttr.Ctime = fsINode.Ctime
	fsAttr.Mtime = fsINode.Mtime
	fsAttr.Atimensec = fsINode.Atimensec
	fsAttr.Ctimensec = fsINode.Ctimensec
	fsAttr.Mtimensec = fsINode.Mtimensec
	fsAttr.Mode = fsINode.Mode
	fsAttr.Nlink = uint32(fsINode.Nlink)
	fsAttr.Uid = fsINode.Uid
	fsAttr.Gid = fsINode.Gid
	fsAttr.Rdev = fsINode.Rdev
}

func (p *DirTreeStg) SetFsEntryOutByFsINode(fsEntryOut *fsapitypes.EntryOut, pFsINode *types.FsINode) {
	fsEntryOut.NodeId = pFsINode.Ino
	fsEntryOut.EntryValid = p.EntryAttrValid
	fsEntryOut.EntryValidNsec = p.EntryAttrValidNsec
	fsEntryOut.AttrValid = p.EntryAttrValid
	fsEntryOut.AttrValidNsec = p.EntryAttrValidNsec
	p.SetFsAttrByFsINode(&fsEntryOut.Attr, pFsINode)
}

func (p *DirTreeStg) SetFsINodeByFsAttr(fsINode *types.FsINode, input *fsapitypes.SetAttrIn) fsapitypes.Status {
	if input.Valid&fsapitypes.FATTR_MODE != 0 {
		// fsINode.Mode = uint32(0777)&input.Mode | uint32(types.FsINodeTypeToFsType(fsINode.Type))
		fsINode.Mode = input.Mode
	}

	if input.Valid&fsapitypes.FATTR_UID != 0 {
		fsINode.Uid = input.Uid
	}
	if input.Valid&fsapitypes.FATTR_GID != 0 {
		fsINode.Gid = input.Gid
	}

	if input.Valid&fsapitypes.FATTR_SIZE != 0 {
		if fsINode.Type == types.FSINODE_TYPE_DIR {
			return fsapitypes.EISDIR
		}
		p.MemStg.NetINodeDriver.NetINodeTruncate(fsINode.UNetINode, input.Size)
	}

	now := p.FsINodeDriver.Timer.Now()
	nowt := types.DirTreeTime(now.Unix())
	nowtnsec := types.DirTreeTimeNsec(now.UnixNano())

	if input.Valid&(fsapitypes.FATTR_ATIME|fsapitypes.FATTR_MTIME|fsapitypes.FATTR_ATIME_NOW|fsapitypes.FATTR_MTIME_NOW) != 0 {
		if input.Valid&fsapitypes.FATTR_ATIME != 0 {
			if input.Valid&fsapitypes.FATTR_ATIME_NOW != 0 {
				fsINode.Atime = nowt
				fsINode.Atimensec = nowtnsec
			} else {
				fsINode.Atime = input.Atime
				fsINode.Atimensec = input.Atimensec
			}
		}

		if input.Valid&fsapitypes.FATTR_MTIME != 0 {
			if input.Valid&fsapitypes.FATTR_MTIME_NOW != 0 {
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

	return fsapitypes.OK
}

func (p *DirTreeStg) SetFsAttrOutByFsINode(fsAttrOut *fsapitypes.AttrOut, pFsINode *types.FsINode) {
	fsAttrOut.AttrValid = p.EntryAttrValid
	fsAttrOut.AttrValidNsec = p.EntryAttrValidNsec
	p.SetFsAttrByFsINode(&fsAttrOut.Attr, pFsINode)
}

func (p *DirTreeStg) GetAttr(input *fsapitypes.GetAttrIn, out *fsapitypes.AttrOut) fsapitypes.Status {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsAttrOutByFsINode(out, &fsINode)

	return fsapitypes.OK
}

func (p *DirTreeStg) SetAttr(input *fsapitypes.SetAttrIn, out *fsapitypes.AttrOut) fsapitypes.Status {
	var (
		fsINode types.FsINode
		code    fsapitypes.Status
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	code = p.SetFsINodeByFsAttr(&fsINode, input)
	if code != fsapitypes.OK {
		return code
	}

	err = p.UpdateFsINodeInDB(&fsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsAttrOutByFsINode(out, &fsINode)

	return fsapitypes.OK
}
