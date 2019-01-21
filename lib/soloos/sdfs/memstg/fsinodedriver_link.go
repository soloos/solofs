package memstg

import (
	"soloos/sdfs/types"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *FsINodeDriver) Link(pointedToFsINode types.FsINodeID,
	parentID types.FsINodeID, filename string, retFsINode *types.FsINode) error {
	var (
		srcFsINode types.FsINode
		err        error
	)

	err = p.FetchFsINodeByID(pointedToFsINode, &srcFsINode)
	if err != nil {
		return err
	}

	p.PrepareFsINodeForCreate(retFsINode,
		&types.ZeroNetINodeID, parentID, filename,
		types.FSINODE_TYPE_HARD_LINK, fuse.S_IFLNK|0777)
	retFsINode.HardLinkIno = srcFsINode.Ino
	err = p.CreateINode(retFsINode)
	if err != nil {
		return err
	}

	srcFsINode.Nlink += 1
	err = p.UpdateFsINodeInDB(&srcFsINode)
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) Symlink(parentID types.FsINodeID, pointedTo string, linkName string,
	retFsINode *types.FsINode) error {
	var (
		parentFsINode types.FsINode
		err           error
	)

	err = p.FetchFsINodeByID(parentID, &parentFsINode)
	if err != nil {
		return err
	}

	p.PrepareFsINodeForCreate(retFsINode,
		&types.ZeroNetINodeID, parentID, linkName,
		types.FSINODE_TYPE_SOFT_LINK, fuse.S_IFLNK|0777)
	err = p.CreateINode(retFsINode)
	if err != nil {
		return err
	}

	err = p.FIXAttrDriver.SetXAttr(retFsINode.Ino, types.FS_XATTR_SOFT_LNKMETA_KEY, []byte(pointedTo))
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) Readlink(fsINodeID types.FsINodeID) ([]byte, error) {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.FetchFsINodeByID(fsINodeID, &fsINode)
	if err != nil {
		return nil, err
	}

	return p.FIXAttrDriver.GetXAttrData(fsINodeID, types.FS_XATTR_SOFT_LNKMETA_KEY)
}
