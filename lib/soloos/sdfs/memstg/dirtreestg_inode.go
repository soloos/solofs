package memstg

import (
	"soloos/sdfs/types"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
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

func (p *DirTreeStg) FetchFsINodeByPath(fsINodePath string, fsINode *types.FsINode) error {
	var (
		paths    []string
		i        int
		parentID types.FsINodeID = types.RootFsINodeID
		err      error
	)

	paths = strings.Split(fsINodePath, "/")

	if paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}

	if len(paths) <= 1 {
		*fsINode = p.FsINodeDriver.RootFsINode
		return nil
	}

	for i = 1; i < len(paths); i++ {
		if paths[i] == "" {
			continue
		}
		err = p.FetchFsINodeByName(parentID, paths[i], fsINode)
		if err != nil {
			return err
		}
		parentID = fsINode.Ino
	}

	return err
}

func (p *DirTreeStg) DeleteFsINodeByPath(fsINodePath string) error {
	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByPath(fsINodePath, &fsINode)
	if err != nil {
		if err == types.ErrObjectNotExists {
			return nil
		} else {
			return err
		}
	}

	err = p.SimpleUnlink(&fsINode)

	return err
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

func (p *DirTreeStg) SimpleOpen(fsINode *types.FsINode, flags uint32, out *fuse.OpenOut) error {
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

func (p *DirTreeStg) Mknod(input *fuse.MknodIn, name string, out *fuse.EntryOut) fuse.Status {
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
		return fuse.EIO
	}

	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &parentFsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.CreateFsINode(&fsINode,
		nil, nil, parentFsINode.Ino,
		name, fsINodeType, input.Mode,
		input.Uid, input.Gid, input.Rdev)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(input.NodeId)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.SetFuseEntryOutByFsINode(out, &fsINode)

	return fuse.OK
}

func (p *DirTreeStg) SimpleUnlink(fsINode *types.FsINode) error {
	return p.FsINodeDriver.UnlinkFsINode(fsINode)
}

func (p *DirTreeStg) Unlink(header *fuse.InHeader, name string) fuse.Status {
	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByName(header.NodeId, name, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.SimpleUnlink(&fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(header.NodeId)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	return fuse.OK
}

func (p *DirTreeStg) Fsync(input *fuse.FsyncIn) fuse.Status {
	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	// TODO flush metadata

	return fuse.OK
}

func (p *DirTreeStg) Lookup(header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
	if len([]byte(name)) > types.FS_MAX_NAME_LENGTH {
		return types.FS_ENAMETOOLONG
	}

	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByName(header.NodeId, name, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.FetchFsINodeByIDThroughHardLink(fsINode.Ino, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.SetFuseEntryOutByFsINode(out, &fsINode)
	return fuse.OK
}

func (p *DirTreeStg) Access(input *fuse.AccessIn) fuse.Status {
	return fuse.OK
}

func (p *DirTreeStg) Forget(nodeid, nlookup uint64) {
}

func (p *DirTreeStg) Release(input *fuse.ReleaseIn) {
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
