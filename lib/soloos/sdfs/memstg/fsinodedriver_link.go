package memstg

import (
	fsapitypes "soloos/common/fsapi/types"
	"soloos/sdfs/types"
	"sync/atomic"
)

func (p *FsINodeDriver) Link(srcFsINode *types.FsINode,
	parentID types.FsINodeID, filename string, retFsINode *types.FsINode) error {
	var (
		err error
	)

	err = p.PrepareFsINodeForCreate(retFsINode,
		nil, nil, parentID,
		filename, types.FSINODE_TYPE_HARD_LINK, srcFsINode.Mode,
		0, 0, types.FS_RDEV)
	if err != nil {
		return err
	}
	retFsINode.HardLinkIno = srcFsINode.Ino
	err = p.CreateFsINode(retFsINode)
	if err != nil {
		return err
	}

	srcFsINode.Nlink += 1
	err = p.UpdateFsINodeInDB(srcFsINode)
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) Symlink(parentID types.FsINodeID, pointedTo string, linkName string,
	retFsINode *types.FsINode) error {
	var (
		err error
	)

	err = p.PrepareFsINodeForCreate(retFsINode,
		nil, nil, parentID,
		linkName, types.FSINODE_TYPE_SOFT_LINK, fsapitypes.S_IFLNK|0777,
		0, 0, types.FS_RDEV)
	if err != nil {
		return err
	}
	err = p.CreateFsINode(retFsINode)
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

func (p *FsINodeDriver) decreaseFsINodeNLink(fsINode *types.FsINode) error {
	var (
		err error
	)

	if atomic.AddInt32(&fsINode.Nlink, -1) > 0 {
		err = p.UpdateFsINodeInDB(fsINode)
		if err != nil {
			return err
		}
		return nil
	}

	// assert fsINode.Nlink == 0
	if fsINode.Type == types.FSINODE_TYPE_HARD_LINK {
		var fsINodeHardLink types.FsINode
		err = p.FetchFsINodeByID(fsINode.HardLinkIno, &fsINodeHardLink)
		if err != nil {
			if err != types.ErrObjectNotExists {
				return err
			}
		} else {
			err = p.decreaseFsINodeNLink(&fsINodeHardLink)
			if err != nil {
				return err
			}
		}
	}

	err = p.helper.DeleteFsINodeByIDInDB(fsINode.Ino)
	if err != nil {
		return err
	}
	p.DeleteFsINodeCache(fsINode.ParentID, fsINode.Name, fsINode.Ino)
	return nil
}

func (p *FsINodeDriver) UnlinkFsINode(fsINode *types.FsINode) error {
	var (
		err error
	)

	err = p.decreaseFsINodeNLink(fsINode)
	if err != nil {
		return err
	}

	err = p.FetchFsINodeByID(fsINode.Ino, fsINode)
	if err != nil {
		if err == types.ErrObjectNotExists {
			return nil
		} else {
			return err
		}
	}

	var parentID = fsINode.ParentID
	fsINode.ParentID = types.ZombieFsINodeParentID
	err = p.UpdateFsINodeInDB(fsINode)
	if err != nil {
		return err
	}
	p.DeleteFsINodeCache(parentID, fsINode.Name, fsINode.Ino)

	return nil
}
