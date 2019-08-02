package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/sdfsapitypes"
	"soloos/sdfs/sdfstypes"
	"strings"
)

func (p *PosixFS) ListFsINodeByIno(ino sdfsapitypes.FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(sdfsapitypes.FsINodeMeta) bool,
) error {
	var (
		uFsINode sdfsapitypes.FsINodeUintptr
		err      error
	)

	uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(ino)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.helper.ListFsINodeByParentIDFromDB(p.NameSpaceID, uFsINode.Ptr().Meta.Ino,
		isFetchAllCols, beforeLiteralFunc, literalFunc)
	if err != nil {
		return err
	}

	return nil
}

func (p *PosixFS) Rename(input *fsapitypes.RenameIn, oldName string, newName string) fsapitypes.Status {
	var (
		oldDirFsINodeID = input.NodeId
		newDirFsINodeID = input.Newdir
		uOldFsINode     sdfsapitypes.FsINodeUintptr
		pOldFsINode     *sdfsapitypes.FsINode
		uCheckFsINode   sdfsapitypes.FsINodeUintptr
		pCheckFsINode   *sdfsapitypes.FsINode
		isDirEmpty      bool
		err             error
	)

	if len([]byte(newName)) > sdfstypes.FS_MAX_NAME_LENGTH {
		return sdfstypes.FS_ENAMETOOLONG
	}

	uOldFsINode, err = p.FsINodeDriver.GetFsINodeByName(oldDirFsINodeID, oldName)
	defer p.FsINodeDriver.ReleaseFsINode(uOldFsINode)
	pOldFsINode = uOldFsINode.Ptr()
	if err != nil {
		return sdfstypes.ErrorToFsStatus(err)
	}

	uCheckFsINode, err = p.FsINodeDriver.GetFsINodeByName(newDirFsINodeID, newName)
	defer p.FsINodeDriver.ReleaseFsINode(uCheckFsINode)
	pCheckFsINode = uCheckFsINode.Ptr()
	if err != nil {
		if err != sdfsapitypes.ErrObjectNotExists {
			return sdfstypes.ErrorToFsStatus(err)
		}
	} else {
		// newName exists
		if pCheckFsINode.Meta.Type == sdfstypes.FSINODE_TYPE_DIR {
			if pOldFsINode.Meta.Type == sdfstypes.FSINODE_TYPE_DIR {
				isDirEmpty, err = p.checkIsDirEmpty(&pCheckFsINode.Meta)
				if err != nil {
					return sdfstypes.ErrorToFsStatus(err)
				}
				if isDirEmpty == false {
					return sdfstypes.FS_ENOTEMPTY
				}
				err = p.FsINodeDriver.UnlinkFsINode(pCheckFsINode.Meta.Ino)
				if err != nil {
					return sdfstypes.ErrorToFsStatus(err)
				}

			} else {
				newDirFsINodeID = pCheckFsINode.Meta.Ino
			}
		} else {
			err = p.FsINodeDriver.UnlinkFsINode(pCheckFsINode.Meta.Ino)
			if err != nil {
				return sdfstypes.ErrorToFsStatus(err)
			}
		}
	}

	pOldFsINode.Meta.ParentID = newDirFsINodeID
	pOldFsINode.Meta.SetName(newName)
	err = p.FsINodeDriver.UpdateFsINodeInDB(&pOldFsINode.Meta)
	if err != nil {
		return sdfstypes.ErrorToFsStatus(err)
	}

	p.FsINodeDriver.CleanFsINodeAssitCache(oldDirFsINodeID, oldName)

	return fsapitypes.OK
}

func (p *PosixFS) SimpleMkdirAll(perms uint32, fsINodePath string, uid uint32, gid uint32) fsapitypes.Status {
	var (
		paths       []string
		i           int
		parentID    sdfsapitypes.FsINodeID = sdfsapitypes.RootFsINodeID
		fsINodeMeta sdfsapitypes.FsINodeMeta
		code        fsapitypes.Status
	)

	paths = strings.Split(fsINodePath, "/")

	if paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}

	for i = 1; i < len(paths); i++ {
		if paths[i] == "" {
			continue
		}

		code = p.SimpleMkdir(&fsINodeMeta, nil, parentID, perms, paths[i], uid, gid, sdfstypes.FS_RDEV)
		if code != fsapitypes.OK && code != sdfstypes.FS_EEXIST {
			goto DONE
		}
		parentID = fsINodeMeta.Ino
	}

DONE:
	if code == sdfstypes.FS_EEXIST {
		code = fsapitypes.OK
	}
	return code
}

func (p *PosixFS) Mkdir(input *fsapitypes.MkdirIn, name string, out *fsapitypes.EntryOut) fsapitypes.Status {
	var (
		fsINodeMeta sdfsapitypes.FsINodeMeta
		code        fsapitypes.Status
		err         error
	)

	code = p.SimpleMkdir(&fsINodeMeta, nil, input.NodeId, input.Mode, name, input.Uid, input.Gid, sdfstypes.FS_RDEV)
	if code != fsapitypes.OK {
		return code
	}

	p.SetFsEntryOutByFsINode(out, &fsINodeMeta)

	err = p.RefreshFsINodeACMtimeByIno(fsINodeMeta.ParentID)
	if err != nil {
		return sdfstypes.ErrorToFsStatus(err)
	}

	return fsapitypes.OK
}

func (p *PosixFS) SimpleMkdir(pFsINodeMeta *sdfsapitypes.FsINodeMeta,
	fsINodeID *sdfsapitypes.FsINodeID, parentID sdfsapitypes.FsINodeID,
	perms uint32, name string,
	uid uint32, gid uint32, rdev uint32) fsapitypes.Status {
	if len([]byte(name)) > sdfstypes.FS_MAX_NAME_LENGTH {
		return sdfstypes.FS_ENAMETOOLONG
	}

	var (
		uFsINode sdfsapitypes.FsINodeUintptr
		err      error
	)

	uFsINode, err = p.FsINodeDriver.GetFsINodeByName(parentID, name)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err == nil {
		*pFsINodeMeta = uFsINode.Ptr().Meta
		return sdfstypes.FS_EEXIST
	}

	if err != nil && err != sdfsapitypes.ErrObjectNotExists {
		return sdfstypes.ErrorToFsStatus(err)
	}

	err = p.createFsINode(pFsINodeMeta,
		fsINodeID, nil, parentID,
		name, sdfstypes.FSINODE_TYPE_DIR, fsapitypes.S_IFDIR|perms,
		uid, gid, rdev)
	if err != nil {
		return sdfstypes.ErrorToFsStatus(err)
	}

	return fsapitypes.OK
}

func (p *PosixFS) checkIsDirEmpty(pFsINodeMeta *sdfsapitypes.FsINodeMeta) (bool, error) {
	var (
		isDirEmpty bool
		err        error
	)

	err = p.ListFsINodeByIno(pFsINodeMeta.Ino, false,
		func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64) {
			isDirEmpty = (resultCount == 0)
			return 0, 0
		},
		func(sdfsapitypes.FsINodeMeta) bool {
			return false
		},
	)
	if err != nil {
		return false, err
	}

	return isDirEmpty, nil
}

func (p *PosixFS) Rmdir(header *fsapitypes.InHeader, name string) fsapitypes.Status {
	var (
		fsINodeMeta sdfsapitypes.FsINodeMeta
		isDirEmpty  bool
		err         error
	)

	err = p.FetchFsINodeByName(&fsINodeMeta, header.NodeId, name)
	if err != nil {
		return sdfstypes.ErrorToFsStatus(err)
	}

	if fsINodeMeta.Type != sdfstypes.FSINODE_TYPE_DIR {
		return fsapitypes.ENOTDIR
	}

	isDirEmpty, err = p.checkIsDirEmpty(&fsINodeMeta)
	if err != nil {
		return sdfstypes.ErrorToFsStatus(err)
	}

	if isDirEmpty == false {
		return sdfstypes.FS_ENOTEMPTY
	}

	err = p.FsINodeDriver.UnlinkFsINode(fsINodeMeta.Ino)
	if err != nil {
		return sdfstypes.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(header.NodeId)
	if err != nil {
		return sdfstypes.ErrorToFsStatus(err)
	}

	return fsapitypes.OK
}

func (p *PosixFS) OpenDir(input *fsapitypes.OpenIn, out *fsapitypes.OpenOut) fsapitypes.Status {
	var (
		fsINodeMeta sdfsapitypes.FsINodeMeta
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, input.NodeId)
	if err != nil {
		return sdfstypes.ErrorToFsStatus(err)
	}

	err = p.SimpleOpen(&fsINodeMeta, input.Flags, out)
	if err != nil {
		return sdfstypes.ErrorToFsStatus(err)
	}
	return fsapitypes.OK
}

func (p *PosixFS) ReadDir(input *fsapitypes.ReadIn, out *fsapitypes.DirEntryList) fsapitypes.Status {
	var (
		fsINodeMetaByIDThroughHardLink sdfsapitypes.FsINodeMeta
		isAddDirEntrySuccess           bool
		err                            error
	)
	err = p.ListFsINodeByIno(input.NodeId, false,
		func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64) {
			return uint64(resultCount) - input.Offset, input.Offset
		},
		func(fsINodeMeta sdfsapitypes.FsINodeMeta) bool {
			if fsINodeMeta.Type == sdfstypes.FSINODE_TYPE_HARD_LINK {
				err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMetaByIDThroughHardLink, fsINodeMeta.Ino)
				if err != nil {
					return false
				}
				isAddDirEntrySuccess, _ = out.AddDirEntry(fsapitypes.DirEntry{
					Mode: fsINodeMetaByIDThroughHardLink.Mode,
					Name: fsINodeMeta.Name(),
					Ino:  fsINodeMetaByIDThroughHardLink.Ino,
				})
			} else {
				isAddDirEntrySuccess, _ = out.AddDirEntry(fsapitypes.DirEntry{
					Mode: fsINodeMeta.Mode,
					Name: fsINodeMeta.Name(),
					Ino:  fsINodeMeta.Ino,
				})

			}
			return isAddDirEntrySuccess
		},
	)
	if err != nil {
		return sdfstypes.ErrorToFsStatus(err)
	}

	return fsapitypes.OK
}

func (p *PosixFS) ReadDirPlus(input *fsapitypes.ReadIn, out *fsapitypes.DirEntryList) fsapitypes.Status {
	var (
		fsINodeMetaByIDThroughHardLink sdfsapitypes.FsINodeMeta
		entryOut                       *fsapitypes.EntryOut
		err                            error
	)
	err = p.ListFsINodeByIno(input.NodeId, true,
		func(resultCount int) (uint64, uint64) {
			var fetchRowsLimit uint64
			if uint64(resultCount) > input.Offset {
				fetchRowsLimit = uint64(resultCount) - input.Offset
				if fetchRowsLimit > 1024 {
					fetchRowsLimit = 1024
				}
			} else {
				fetchRowsLimit = 0
			}
			return fetchRowsLimit, input.Offset
		},
		func(fsINodeMeta sdfsapitypes.FsINodeMeta) bool {
			err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMetaByIDThroughHardLink, fsINodeMeta.Ino)
			if err != nil {
				return false
			}

			entryOut, _ = out.AddDirLookupEntry(fsapitypes.DirEntry{
				Mode: fsINodeMetaByIDThroughHardLink.Mode,
				Name: fsINodeMeta.Name(),
				Ino:  fsINodeMetaByIDThroughHardLink.Ino,
			})
			if entryOut == nil {
				return false
			}

			p.SetFsEntryOutByFsINode(entryOut, &fsINodeMetaByIDThroughHardLink)
			return true
		},
	)
	if err != nil {
		return sdfstypes.ErrorToFsStatus(err)
	}
	return fsapitypes.OK
}

func (p *PosixFS) ReleaseDir(input *fsapitypes.ReleaseIn) {
	// TODO make sure releaable
	p.FdTable.ReleaseFd(input.Fh)
}

func (p *PosixFS) FsyncDir(input *fsapitypes.FsyncIn) fsapitypes.Status {
	return fsapitypes.OK
}
