package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/sdfsapitypes"
	"soloos/sdfs/sdfstypes"
	"sync/atomic"
)

func (p *FsINodeDriver) Link(srcFsINodeID sdfsapitypes.FsINodeID,
	parentID sdfsapitypes.FsINodeID, filename string,
	retFsINode *sdfsapitypes.FsINodeMeta) error {
	var (
		uSrcFsINode sdfsapitypes.FsINodeUintptr
		pSrcFsINode *sdfsapitypes.FsINode
		err         error
	)

	uSrcFsINode, err = p.GetFsINodeByID(srcFsINodeID)
	defer p.ReleaseFsINode(uSrcFsINode)
	pSrcFsINode = uSrcFsINode.Ptr()
	if err != nil {
		return err
	}

	err = p.PrepareFsINodeForCreate(retFsINode,
		nil, nil, parentID,
		filename, sdfstypes.FSINODE_TYPE_HARD_LINK, pSrcFsINode.Meta.Mode,
		0, 0, sdfstypes.FS_RDEV)
	if err != nil {
		return err
	}
	retFsINode.HardLinkIno = pSrcFsINode.Meta.Ino
	err = p.CreateFsINode(retFsINode)
	if err != nil {
		return err
	}

	atomic.AddInt32(&pSrcFsINode.Meta.Nlink, 1)
	err = p.UpdateFsINodeInDB(&pSrcFsINode.Meta)
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) Symlink(parentID sdfsapitypes.FsINodeID, pointedTo string, linkName string,
	retFsINodeMeta *sdfsapitypes.FsINodeMeta) error {
	var (
		err error
	)

	err = p.PrepareFsINodeForCreate(retFsINodeMeta,
		nil, nil, parentID,
		linkName, sdfstypes.FSINODE_TYPE_SOFT_LINK, fsapitypes.S_IFLNK|0777,
		0, 0, sdfstypes.FS_RDEV)
	if err != nil {
		return err
	}
	err = p.CreateFsINode(retFsINodeMeta)
	if err != nil {
		return err
	}

	err = p.FIXAttrDriver.SetXAttr(retFsINodeMeta.Ino, sdfstypes.FS_XATTR_SOFT_LNKMETA_KEY, []byte(pointedTo))
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) Readlink(fsINodeID sdfsapitypes.FsINodeID) ([]byte, error) {
	return p.FIXAttrDriver.GetXAttrData(fsINodeID, sdfstypes.FS_XATTR_SOFT_LNKMETA_KEY)
}

func (p *FsINodeDriver) deleteFsINodeHardLinked(fsINodeID sdfsapitypes.FsINodeID) error {
	var (
		uFsINode                sdfsapitypes.FsINodeUintptr
		isFsINodeDeleted        bool
		isShouldReleaseUFsINode bool
		err                     error
	)
	uFsINode, err = p.GetFsINodeByID(fsINodeID)
	if err != nil {
		isShouldReleaseUFsINode = true

	} else {
		isFsINodeDeleted, err = p.decreaseFsINodeNLink(uFsINode)
		if err != nil {
			isShouldReleaseUFsINode = true
		}

		if isFsINodeDeleted == false {
			isShouldReleaseUFsINode = true
		}
	}

	if isShouldReleaseUFsINode {
		p.ReleaseFsINode(uFsINode)
	}

	return err
}

// decreaseFsINodeNLink return isFsINodeDeleted, decreaseError
// if FsINodeMeta.Nlink == 0 then delete FsINode else decrease(FsINode.Nlink)
func (p *FsINodeDriver) decreaseFsINodeNLink(uFsINode sdfsapitypes.FsINodeUintptr) (bool, error) {
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
	if pFsINode.Meta.Type == sdfstypes.FSINODE_TYPE_HARD_LINK {
		err = p.deleteFsINodeHardLinked(pFsINode.Meta.HardLinkIno)
		if err != nil {
			return false, err
		}

	}

	err = p.helper.DeleteFsINodeByIDInDB(p.posixFS.NameSpaceID, pFsINode.Meta.Ino)
	if err != nil {
		return false, err
	}

	p.DeleteFsINodeCache(uFsINode, pFsINode.Meta.ParentID, pFsINode.Meta.Name())

	return true, nil
}

func (p *FsINodeDriver) UnlinkFsINode(fsINodeID sdfsapitypes.FsINodeID) error {
	var (
		uFsINode                sdfsapitypes.FsINodeUintptr
		pFsINode                *sdfsapitypes.FsINode
		oldFsINodeParentID      sdfsapitypes.FsINodeID
		isFsINodeDeleted        bool
		isShouldReleaseUFsINode bool
		err                     error
	)

	uFsINode, err = p.GetFsINodeByID(fsINodeID)
	pFsINode = uFsINode.Ptr()

	if err != nil {
		if err == sdfsapitypes.ErrObjectNotExists {
			err = nil
			goto DONE
		} else {
			goto DONE
		}
	}

	oldFsINodeParentID = pFsINode.Meta.ParentID
	isFsINodeDeleted, err = p.decreaseFsINodeNLink(uFsINode)
	if err != nil {
		return err
	}

	if isFsINodeDeleted == false {
		pFsINode.Meta.ParentID = sdfsapitypes.ZombieFsINodeParentID
		err = p.UpdateFsINodeInDB(&pFsINode.Meta)
		if err != nil {
			goto DONE
		}
		isShouldReleaseUFsINode = true
	}

DONE:
	if uFsINode != 0 {
		p.CleanFsINodeAssitCache(oldFsINodeParentID, uFsINode.Ptr().Meta.Name())
		if isShouldReleaseUFsINode {
			p.ReleaseFsINode(uFsINode)
		}
	}
	return nil
}
