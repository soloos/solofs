package memstg

import (
	fsapitypes "soloos/common/fsapi/types"
	sdfsapitypes "soloos/common/sdfsapi/types"
	"soloos/sdfs/types"
)

func (p *DirTreeStg) SetFsAttrByFsINode(fsAttr *fsapitypes.Attr, pFsINodeMeta *types.FsINodeMeta) {
	fsAttr.Ino = pFsINodeMeta.Ino

	if pFsINodeMeta.NetINodeID != sdfsapitypes.ZeroNetINodeID ||
		pFsINodeMeta.Type == types.FSINODE_TYPE_HARD_LINK {
		var uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(pFsINodeMeta.Ino)
		defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
		var pFsINode = uFsINode.Ptr()
		if err == nil {
			fsAttr.Size = pFsINode.UNetINode.Ptr().Size
			fsAttr.Blocks = pFsINode.UNetINode.Ptr().GetBlocks()
			fsAttr.Blksize = uint32(pFsINode.UNetINode.Ptr().MemBlockCap)
			fsAttr.Padding = uint32(pFsINode.UNetINode.Ptr().MemBlockCap)
		}
	}

	fsAttr.Atime = pFsINodeMeta.Atime
	fsAttr.Ctime = pFsINodeMeta.Ctime
	fsAttr.Mtime = pFsINodeMeta.Mtime
	fsAttr.Atimensec = pFsINodeMeta.Atimensec
	fsAttr.Ctimensec = pFsINodeMeta.Ctimensec
	fsAttr.Mtimensec = pFsINodeMeta.Mtimensec
	fsAttr.Mode = pFsINodeMeta.Mode
	fsAttr.Nlink = uint32(pFsINodeMeta.Nlink)
	fsAttr.Uid = pFsINodeMeta.Uid
	fsAttr.Gid = pFsINodeMeta.Gid
	fsAttr.Rdev = pFsINodeMeta.Rdev
}

func (p *DirTreeStg) SetFsEntryOutByFsINode(fsEntryOut *fsapitypes.EntryOut,
	pFsINodeMeta *types.FsINodeMeta) {

	fsEntryOut.NodeId = pFsINodeMeta.Ino
	fsEntryOut.EntryValid = p.EntryAttrValid
	fsEntryOut.EntryValidNsec = p.EntryAttrValidNsec
	fsEntryOut.AttrValid = p.EntryAttrValid
	fsEntryOut.AttrValidNsec = p.EntryAttrValidNsec
	p.SetFsAttrByFsINode(&fsEntryOut.Attr, pFsINodeMeta)
}

func (p *DirTreeStg) SetFsINodeByFsAttr(pFsINodeMeta *types.FsINodeMeta,
	input *fsapitypes.SetAttrIn) fsapitypes.Status {

	if input.Valid&fsapitypes.FATTR_MODE != 0 {
		// pFsINodeMeta.Mode = uint32(0777)&input.Mode | uint32(types.FsINodeTypeToFsType(pFsINodeMeta.Type))
		pFsINodeMeta.Mode = input.Mode
	}

	if input.Valid&fsapitypes.FATTR_UID != 0 {
		pFsINodeMeta.Uid = input.Uid
	}
	if input.Valid&fsapitypes.FATTR_GID != 0 {
		pFsINodeMeta.Gid = input.Gid
	}

	if input.Valid&fsapitypes.FATTR_SIZE != 0 {
		if pFsINodeMeta.Type == types.FSINODE_TYPE_DIR {
			return fsapitypes.EISDIR
		}
		var err = p.TruncateINode(pFsINodeMeta, input.Size)
		if err != nil {
			return fsapitypes.EIO
		}
	}

	now := p.FsINodeDriver.Timer.Now()
	nowt := types.DirTreeTime(now.Unix())
	nowtnsec := types.DirTreeTimeNsec(now.UnixNano())

	if input.Valid&(fsapitypes.FATTR_ATIME|fsapitypes.FATTR_MTIME|fsapitypes.FATTR_ATIME_NOW|fsapitypes.FATTR_MTIME_NOW) != 0 {
		if input.Valid&fsapitypes.FATTR_ATIME != 0 {
			if input.Valid&fsapitypes.FATTR_ATIME_NOW != 0 {
				pFsINodeMeta.Atime = nowt
				pFsINodeMeta.Atimensec = nowtnsec
			} else {
				pFsINodeMeta.Atime = input.Atime
				pFsINodeMeta.Atimensec = input.Atimensec
			}
		}

		if input.Valid&fsapitypes.FATTR_MTIME != 0 {
			if input.Valid&fsapitypes.FATTR_MTIME_NOW != 0 {
				pFsINodeMeta.Mtime = nowt
				pFsINodeMeta.Mtimensec = nowtnsec
			} else {
				pFsINodeMeta.Mtime = input.Mtime
				pFsINodeMeta.Mtimensec = input.Mtimensec
			}
		}
	}

	pFsINodeMeta.Ctime = nowt
	pFsINodeMeta.Ctimensec = nowtnsec

	return fsapitypes.OK
}

func (p *DirTreeStg) SetFsAttrOutByFsINode(fsAttrOut *fsapitypes.AttrOut,
	pFsINodeMeta *types.FsINodeMeta) {

	fsAttrOut.AttrValid = p.EntryAttrValid
	fsAttrOut.AttrValidNsec = p.EntryAttrValidNsec
	p.SetFsAttrByFsINode(&fsAttrOut.Attr, pFsINodeMeta)
}

func (p *DirTreeStg) GetAttr(input *fsapitypes.GetAttrIn, out *fsapitypes.AttrOut) fsapitypes.Status {
	var (
		fsINodeMeta types.FsINodeMeta
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, input.NodeId)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsAttrOutByFsINode(out, &fsINodeMeta)

	return fsapitypes.OK
}

func (p *DirTreeStg) SetAttr(input *fsapitypes.SetAttrIn, out *fsapitypes.AttrOut) fsapitypes.Status {
	var (
		fsINodeMeta types.FsINodeMeta
		code        fsapitypes.Status
		err         error
	)

	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, input.NodeId)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	code = p.SetFsINodeByFsAttr(&fsINodeMeta, input)
	if code != fsapitypes.OK {
		return code
	}

	err = p.FsINodeDriver.UpdateFsINodeInDB(&fsINodeMeta)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsAttrOutByFsINode(out, &fsINodeMeta)

	return fsapitypes.OK
}
