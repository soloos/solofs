package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/solofsapitypes"
	"soloos/solofs/solofstypes"
)

func (p *PosixFs) FetchFsINodeByID(pFsINodeMeta *solofsapitypes.FsINodeMeta,
	fsINodeID solofsapitypes.FsINodeID) error {
	var (
		uFsINode solofsapitypes.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.FsINodeDriver.GetFsINodeByID(fsINodeID)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}

	*pFsINodeMeta = uFsINode.Ptr().Meta
	return nil
}

func (p *PosixFs) FetchFsINodeByName(pFsINodeMeta *solofsapitypes.FsINodeMeta,
	parentID solofsapitypes.FsINodeID, fsINodeName string) error {
	var (
		uFsINode solofsapitypes.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.FsINodeDriver.GetFsINodeByName(parentID, fsINodeName)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}

	*pFsINodeMeta = uFsINode.Ptr().Meta
	return nil
}

func (p *PosixFs) FetchFsINodeByIDThroughHardLink(pFsINodeMeta *solofsapitypes.FsINodeMeta,
	fsINodeID solofsapitypes.FsINodeID) error {
	var (
		uFsINode solofsapitypes.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(fsINodeID)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}

	*pFsINodeMeta = uFsINode.Ptr().Meta
	return nil
}

func (p *PosixFs) createFsINode(pFsINodeMeta *solofsapitypes.FsINodeMeta,
	fsINodeID *solofsapitypes.FsINodeID,
	netINodeID *solofsapitypes.NetINodeID, parentID solofsapitypes.FsINodeID,
	name string, fsINodeType int, mode uint32,
	uid uint32, gid uint32, rdev uint32,
) error {
	var (
		err error
	)
	err = p.FsINodeDriver.PrepareFsINodeForCreate(pFsINodeMeta,
		fsINodeID, netINodeID, parentID,
		name, fsINodeType, mode,
		uid, gid, rdev)
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.CreateFsINode(pFsINodeMeta)
	if err != nil {
		return err
	}

	return nil
}

func (p *PosixFs) SimpleOpen(fsINodeMeta *solofsapitypes.FsINodeMeta,
	flags uint32, out *fsapitypes.OpenOut) error {
	out.Fh = p.FdTable.AllocFd(fsINodeMeta.Ino)
	out.OpenFlags = flags
	return nil
}

func (p *PosixFs) Mknod(input *fsapitypes.MknodIn, name string, out *fsapitypes.EntryOut) fsapitypes.Status {
	if len([]byte(name)) > solofstypes.FS_MAX_NAME_LENGTH {
		return solofstypes.FS_ENAMETOOLONG
	}

	var (
		parentFsINodeMeta solofsapitypes.FsINodeMeta
		fsINodeMeta       solofsapitypes.FsINodeMeta
		fsINodeType       int
		err               error
	)

	fsINodeType = solofstypes.FsModeToFsINodeType(input.Mode)
	if fsINodeType == solofstypes.FSINODE_TYPE_UNKOWN {
		return fsapitypes.EIO
	}

	err = p.FetchFsINodeByIDThroughHardLink(&parentFsINodeMeta, input.NodeId)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	err = p.createFsINode(&fsINodeMeta,
		nil, nil, parentFsINodeMeta.Ino,
		name, fsINodeType, input.Mode,
		input.Uid, input.Gid, input.Rdev)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(input.NodeId)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &fsINodeMeta)

	return fsapitypes.OK
}

func (p *PosixFs) Unlink(header *fsapitypes.InHeader, name string) fsapitypes.Status {
	var (
		fsINodeMeta solofsapitypes.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByName(&fsINodeMeta, header.NodeId, name)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	err = p.FsINodeDriver.UnlinkFsINode(fsINodeMeta.Ino)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(header.NodeId)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	return fsapitypes.OK
}

func (p *PosixFs) Fsync(input *fsapitypes.FsyncIn) fsapitypes.Status {
	var (
		fsINodeMeta solofsapitypes.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, input.NodeId)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	// TODO flush metadata

	return fsapitypes.OK
}

func (p *PosixFs) Lookup(header *fsapitypes.InHeader, name string, out *fsapitypes.EntryOut) fsapitypes.Status {
	if len([]byte(name)) > solofstypes.FS_MAX_NAME_LENGTH {
		return solofstypes.FS_ENAMETOOLONG
	}

	var (
		fsINodeMeta solofsapitypes.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByName(&fsINodeMeta, header.NodeId, name)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, fsINodeMeta.Ino)
	if err != nil {
		return solofstypes.ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &fsINodeMeta)
	return fsapitypes.OK
}

func (p *PosixFs) Access(input *fsapitypes.AccessIn) fsapitypes.Status {
	return fsapitypes.OK
}

func (p *PosixFs) Forget(nodeid, nlookup uint64) {
}

func (p *PosixFs) Release(input *fsapitypes.ReleaseIn) {
}

func (p *PosixFs) CheckPermissionChmod(uid uint32, gid uint32,
	fsINodeMeta *solofsapitypes.FsINodeMeta) bool {

	if uid == 0 || uid == fsINodeMeta.Uid {
		return true
	}

	return false
}

func (p *PosixFs) CheckPermissionRead(uid uint32, gid uint32,
	fsINodeMeta *solofsapitypes.FsINodeMeta) bool {

	perm := uint32(07777) & fsINodeMeta.Mode
	if uid == fsINodeMeta.Uid {
		if perm&solofstypes.FS_PERM_USER_READ != 0 {
			return true
		}
	}

	if gid == fsINodeMeta.Gid {
		if perm&solofstypes.FS_PERM_GROUP_READ != 0 {
			return true
		}
	}

	if perm&solofstypes.FS_PERM_OTHER_READ != 0 {
		return true
	}

	return false
}

func (p *PosixFs) CheckPermissionWrite(uid uint32, gid uint32,
	fsINodeMeta *solofsapitypes.FsINodeMeta) bool {

	perm := uint32(07777) & fsINodeMeta.Mode
	if uid == fsINodeMeta.Uid {
		if perm&solofstypes.FS_PERM_USER_WRITE != 0 {
			return true
		}
	}

	if gid == fsINodeMeta.Gid {
		if perm&solofstypes.FS_PERM_GROUP_WRITE != 0 {
			return true
		}
	}

	if perm&solofstypes.FS_PERM_OTHER_WRITE != 0 {
		return true
	}

	return false
}

func (p *PosixFs) CheckPermissionExecute(uid uint32, gid uint32,
	fsINodeMeta *solofsapitypes.FsINodeMeta) bool {

	perm := uint32(07777) & fsINodeMeta.Mode
	if uid == fsINodeMeta.Uid {
		if perm&solofstypes.FS_PERM_USER_EXECUTE != 0 {
			return true
		}
	}

	if gid == fsINodeMeta.Gid {
		if perm&solofstypes.FS_PERM_GROUP_EXECUTE != 0 {
			return true
		}
	}

	if perm&solofstypes.FS_PERM_OTHER_EXECUTE != 0 {
		return true
	}

	return false
}

func (p *PosixFs) RefreshFsINodeACMtimeByIno(fsINodeID solofsapitypes.FsINodeID) error {
	return p.FsINodeDriver.RefreshFsINodeACMtimeByIno(fsINodeID)
}

func (p *PosixFs) TruncateINode(pFsINode *solofsapitypes.FsINodeMeta, size uint64) error {
	var (
		uFsINode solofsapitypes.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.FsINodeDriver.GetFsINodeByID(pFsINode.Ino)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}
	if uFsINode.Ptr().UNetINode == 0 {
		return nil
	}

	return p.MemStg.NetINodeDriver.NetINodeTruncate(uFsINode.Ptr().UNetINode, size)
}
