package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofstypes"
)

func (p *PosixFs) SetFsAttrByFsINode(fsAttr *fsapitypes.Attr, pFsINodeMeta *solofstypes.FsINodeMeta) {
	fsAttr.Ino = pFsINodeMeta.Ino

	if pFsINodeMeta.NetINodeID != solofstypes.ZeroNetINodeID ||
		pFsINodeMeta.Type == solofstypes.FSINODE_TYPE_HARD_LINK {
		var uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(pFsINodeMeta.Ino)
		defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
		var pFsINode = uFsINode.Ptr()
		if err == nil && pFsINode.UNetINode != 0 {
			fsAttr.Size = pFsINode.UNetINode.Ptr().Size
			fsAttr.Blocks = pFsINode.UNetINode.Ptr().GetBlocks()
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
