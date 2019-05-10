package memstg

import (
	fsapitypes "soloos/common/fsapi/types"
	sdfsapitypes "soloos/common/sdfsapi/types"
	"soloos/sdfs/types"
	"sync/atomic"
)

func (p *FsINodeDriver) Link(srcFsINodeMeta *types.FsINodeMeta,
	parentID types.FsINodeID, filename string,
	retFsINode *types.FsINodeMeta) error {
	var (
		err error
	)

	err = p.PrepareFsINodeForCreate(retFsINode,
		nil, nil, parentID,
		filename, types.FSINODE_TYPE_HARD_LINK, srcFsINodeMeta.Mode,
		0, 0, types.FS_RDEV)
	if err != nil {
		return err
	}
	retFsINode.HardLinkIno = srcFsINodeMeta.Ino
	err = p.CreateFsINode(retFsINode)
	if err != nil {
		return err
	}

	srcFsINodeMeta.Nlink += 1
	err = p.UpdateFsINodeInDB(srcFsINodeMeta)
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) Symlink(parentID types.FsINodeID, pointedTo string, linkName string,
	retFsINodeMeta *types.FsINodeMeta) error {
	var (
		err error
	)

	err = p.PrepareFsINodeForCreate(retFsINodeMeta,
		nil, nil, parentID,
		linkName, types.FSINODE_TYPE_SOFT_LINK, fsapitypes.S_IFLNK|0777,
		0, 0, types.FS_RDEV)
	if err != nil {
		return err
	}
	err = p.CreateFsINode(retFsINodeMeta)
	if err != nil {
		return err
	}

	err = p.FIXAttrDriver.SetXAttr(retFsINodeMeta.Ino, types.FS_XATTR_SOFT_LNKMETA_KEY, []byte(pointedTo))
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) Readlink(fsINodeID types.FsINodeID) ([]byte, error) {
	return p.FIXAttrDriver.GetXAttrData(fsINodeID, types.FS_XATTR_SOFT_LNKMETA_KEY)
}

// decreaseFsINodeNLink return isFsINodeDeleted, decreaseError
// if FsINodeMeta.Nlink == 0 then delete FsINode else decrease(FsINode.Nlink)
func (p *FsINodeDriver) decreaseFsINodeNLink(uFsINode types.FsINodeUintptr) (bool, error) {
	var (
		pFsINode = uFsINode.Ptr()
		err      error
	)

	if atomic.AddInt32(&pFsINode.Meta.Nlink, -1) > 0 {
		err = p.UpdateFsINodeInDB(&pFsINode.Meta)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	// assert fsINode.Nlink == 0
	if pFsINode.Meta.Type == types.FSINODE_TYPE_HARD_LINK {
		var uFsINodeHardLink types.FsINodeUintptr
		uFsINodeHardLink, err = p.GetFsINodeByID(pFsINode.Meta.HardLinkIno)
		defer p.ReleaseFsINode(uFsINodeHardLink)
		if err != nil {
			if err != types.ErrObjectNotExists {
				return false, err
			}
		} else {
			_, err = p.decreaseFsINodeNLink(uFsINodeHardLink)
			if err != nil {
				return false, err
			}
		}
	}

	err = p.helper.DeleteFsINodeByIDInDB(pFsINode.Meta.Ino)
	if err != nil {
		return false, err
	}

	p.DeleteFsINodeCache(uFsINode, pFsINode.Meta.ParentID, pFsINode.Meta.Name())

	return true, nil
}

func (p *FsINodeDriver) UnlinkFsINode(fsINodeID types.FsINodeID) error {
	var (
		uFsINode               types.FsINodeUintptr
		pFsINode               *types.FsINode
		isFsINodeDeleted       bool
		isShouldDeleteUFsINode bool
		err                    error
	)

	uFsINode, err = p.GetFsINodeByID(fsINodeID)
	pFsINode = uFsINode.Ptr()
	defer func(uFsINode types.FsINodeUintptr, parentID types.FsINodeID) {
		if isShouldDeleteUFsINode {
			p.DeleteFsINodeCache(uFsINode, pFsINode.Meta.ParentID, pFsINode.Meta.Name())
		} else {
			p.ReleaseFsINode(uFsINode)
		}
	}(uFsINode, pFsINode.Meta.ParentID)

	if err != nil {
		if err == types.ErrObjectNotExists {
			return nil
		} else {
			return err
		}
	}

	isFsINodeDeleted, err = p.decreaseFsINodeNLink(uFsINode)
	if err != nil {
		return err
	}

	if isFsINodeDeleted == false {
		pFsINode.Meta.ParentID = sdfsapitypes.ZombieFsINodeParentID
		err = p.UpdateFsINodeInDB(&pFsINode.Meta)
		if err != nil {
			return err
		}
		isShouldDeleteUFsINode = true
	}

	return nil
}
