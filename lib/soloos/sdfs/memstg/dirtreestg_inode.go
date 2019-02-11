package memstg

import (
	fsapitypes "soloos/fsapi/types"
	"soloos/sdfs/types"
)

func (p *DirTreeStg) UpdateFsINodeInDB(pFsINode *types.FsINode) error {
	return p.FsINodeDriver.UpdateFsINodeInDB(pFsINode)
}

// maybe param inode means link
func (p *DirTreeStg) returnRealINode(pFsINode *types.FsINode) types.FsINodeID {
	// TODO complete me
	return pFsINode.Ino
}

func (p *DirTreeStg) FetchFsINodeByIDThroughHardLink(fsINodeID types.FsINodeID, fsINode *types.FsINode) error {
	return p.FsINodeDriver.FetchFsINodeByIDThroughHardLink(fsINodeID, fsINode)
}

func (p *DirTreeStg) FetchFsINodeByID(fsINodeID types.FsINodeID, fsINode *types.FsINode) error {
	return p.FsINodeDriver.FetchFsINodeByID(fsINodeID, fsINode)
}

func (p *DirTreeStg) FetchFsINodeByName(parentID types.FsINodeID, fsINodeName string, fsINode *types.FsINode) error {
	return p.FsINodeDriver.FetchFsINodeByName(parentID, fsINodeName, fsINode)
}

func (p *DirTreeStg) CreateFsINode(fsINode *types.FsINode,
	fsINodeID *types.FsINodeID, netINodeID *types.NetINodeID, parentID types.FsINodeID,
	name string, fsINodeType int, mode uint32,
	uid uint32, gid uint32, rdev uint32,
) error {
	var (
		err error
	)
	err = p.FsINodeDriver.PrepareFsINodeForCreate(fsINode,
		fsINodeID, netINodeID, parentID,
		name, fsINodeType, mode,
		uid, gid, rdev)
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.CreateFsINode(fsINode)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeStg) SimpleOpen(fsINode *types.FsINode, flags uint32, out *fsapitypes.OpenOut) error {
	out.Fh = p.FdTable.AllocFd(fsINode.Ino)
	out.OpenFlags = flags
	return nil
}

func (p *DirTreeStg) getFsINodePathLen(fsINode *types.FsINode, startFsNodeID types.FsINodeID) (int, error) {
	var (
		pathLen int
		err     error
	)
	for fsINode.Ino = startFsNodeID; fsINode.Ino != types.RootFsINodeID; fsINode.Ino = fsINode.ParentID {
		err = p.FetchFsINodeByID(fsINode.Ino, fsINode)
		if err != nil {
			return 0, err
		}
		pathLen += (len(fsINode.Name) + 1)
	}
	return pathLen, nil
}

func (p *DirTreeStg) Mknod(input *fsapitypes.MknodIn, name string, out *fsapitypes.EntryOut) fsapitypes.Status {
	if len([]byte(name)) > types.FS_MAX_NAME_LENGTH {
		return types.FS_ENAMETOOLONG
	}

	var (
		parentFsINode types.FsINode
		fsINode       types.FsINode
		fsINodeType   int
		err           error
	)

	fsINodeType = types.FsModeToFsINodeType(input.Mode)
	if fsINodeType == types.FSINODE_TYPE_UNKOWN {
		return fsapitypes.EIO
	}

	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &parentFsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.CreateFsINode(&fsINode,
		nil, nil, parentFsINode.Ino,
		name, fsINodeType, input.Mode,
		input.Uid, input.Gid, input.Rdev)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(input.NodeId)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &fsINode)

	return fsapitypes.OK
}

func (p *DirTreeStg) SimpleUnlink(fsINode *types.FsINode) error {
	return p.FsINodeDriver.UnlinkFsINode(fsINode)
}

func (p *DirTreeStg) Unlink(header *fsapitypes.InHeader, name string) fsapitypes.Status {
	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByName(header.NodeId, name, &fsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.SimpleUnlink(&fsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(header.NodeId)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	return fsapitypes.OK
}

func (p *DirTreeStg) Fsync(input *fsapitypes.FsyncIn) fsapitypes.Status {
	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	// TODO flush metadata

	return fsapitypes.OK
}

func (p *DirTreeStg) Lookup(header *fsapitypes.InHeader, name string, out *fsapitypes.EntryOut) fsapitypes.Status {
	if len([]byte(name)) > types.FS_MAX_NAME_LENGTH {
		return types.FS_ENAMETOOLONG
	}

	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByName(header.NodeId, name, &fsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.FetchFsINodeByIDThroughHardLink(fsINode.Ino, &fsINode)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &fsINode)
	return fsapitypes.OK
}

func (p *DirTreeStg) Access(input *fsapitypes.AccessIn) fsapitypes.Status {
	return fsapitypes.OK
}

func (p *DirTreeStg) Forget(nodeid, nlookup uint64) {
}

func (p *DirTreeStg) Release(input *fsapitypes.ReleaseIn) {
}

func (p *DirTreeStg) CheckPermissionChmod(uid uint32, gid uint32, fsINode *types.FsINode) bool {
	if uid == 0 || uid == fsINode.Uid {
		return true
	}

	return false
}

func (p *DirTreeStg) CheckPermissionRead(uid uint32, gid uint32, fsINode *types.FsINode) bool {
	perm := uint32(07777) & fsINode.Mode
	if uid == fsINode.Uid {
		if perm&types.FS_PERM_USER_READ != 0 {
			return true
		}
	}

	if gid == fsINode.Gid {
		if perm&types.FS_PERM_GROUP_READ != 0 {
			return true
		}
	}

	if perm&types.FS_PERM_OTHER_READ != 0 {
		return true
	}

	return false
}

func (p *DirTreeStg) CheckPermissionWrite(uid uint32, gid uint32, fsINode *types.FsINode) bool {
	perm := uint32(07777) & fsINode.Mode
	if uid == fsINode.Uid {
		if perm&types.FS_PERM_USER_WRITE != 0 {
			return true
		}
	}

	if gid == fsINode.Gid {
		if perm&types.FS_PERM_GROUP_WRITE != 0 {
			return true
		}
	}

	if perm&types.FS_PERM_OTHER_WRITE != 0 {
		return true
	}

	return false
}

func (p *DirTreeStg) CheckPermissionExecute(uid uint32, gid uint32, fsINode *types.FsINode) bool {
	perm := uint32(07777) & fsINode.Mode
	if uid == fsINode.Uid {
		if perm&types.FS_PERM_USER_EXECUTE != 0 {
			return true
		}
	}

	if gid == fsINode.Gid {
		if perm&types.FS_PERM_GROUP_EXECUTE != 0 {
			return true
		}
	}

	if perm&types.FS_PERM_OTHER_EXECUTE != 0 {
		return true
	}

	return false
}

func (p *DirTreeStg) RefreshFsINodeACMtime(fsINode *types.FsINode) error {
	return p.FsINodeDriver.RefreshFsINodeACMtime(fsINode)
}

func (p *DirTreeStg) RefreshFsINodeACMtimeByIno(fsINodeID types.FsINodeID) error {
	return p.FsINodeDriver.RefreshFsINodeACMtimeByIno(fsINodeID)
}
