package memstg

import (
	fsapitypes "soloos/common/fsapi/types"
	sdfsapitypes "soloos/common/sdfsapi/types"
	"soloos/sdfs/types"
	"strings"
)

func (p *DirTreeStg) ListFsINodeByIno(ino sdfsapitypes.FsINodeID,
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

	err = p.FsINodeDriver.helper.ListFsINodeByParentIDFromDB(uFsINode.Ptr().Meta.Ino,
		isFetchAllCols, beforeLiteralFunc, literalFunc)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeStg) Rename(input *fsapitypes.RenameIn, oldName string, newName string) fsapitypes.Status {
	var (
		oldDirFsINodeID = input.NodeId
		newDirFsINodeID = input.Newdir
		uOldFsINode     types.FsINodeUintptr
		pOldFsINode     *types.FsINode
		uCheckFsINode   types.FsINodeUintptr
		pCheckFsINode   *types.FsINode
		isDirEmpty      bool
		err             error
	)

	if len([]byte(newName)) > types.FS_MAX_NAME_LENGTH {
		return types.FS_ENAMETOOLONG
	}

	uOldFsINode, err = p.FsINodeDriver.GetFsINodeByName(oldDirFsINodeID, oldName)
	defer p.FsINodeDriver.ReleaseFsINode(uOldFsINode)
	pOldFsINode = uOldFsINode.Ptr()
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	uCheckFsINode, err = p.FsINodeDriver.GetFsINodeByName(newDirFsINodeID, newName)
	defer p.FsINodeDriver.ReleaseFsINode(uCheckFsINode)
	pCheckFsINode = uCheckFsINode.Ptr()
	if err != nil {
		if err != types.ErrObjectNotExists {
			return types.ErrorToFsStatus(err)
		}
	} else {
		// newName exists
		if pCheckFsINode.Meta.Type == types.FSINODE_TYPE_DIR {
			if pOldFsINode.Meta.Type == types.FSINODE_TYPE_DIR {
				isDirEmpty, err = p.checkIsDirEmpty(&pCheckFsINode.Meta)
				if err != nil {
					return types.ErrorToFsStatus(err)
				}
				if isDirEmpty == false {
					return types.FS_ENOTEMPTY
				}
				err = p.FsINodeDriver.UnlinkFsINode(pCheckFsINode.Meta.Ino)
				if err != nil {
					return types.ErrorToFsStatus(err)
				}

			} else {
				newDirFsINodeID = pCheckFsINode.Meta.Ino
			}
		} else {
			err = p.FsINodeDriver.UnlinkFsINode(pCheckFsINode.Meta.Ino)
			if err != nil {
				return types.ErrorToFsStatus(err)
			}
		}
	}

	pOldFsINode.Meta.ParentID = newDirFsINodeID
	pOldFsINode.Meta.SetName(newName)
	err = p.FsINodeDriver.UpdateFsINodeInDB(&pOldFsINode.Meta)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.FsINodeDriver.CleanFsINodeAssitCache(oldDirFsINodeID, oldName)

	return fsapitypes.OK
}

func (p *DirTreeStg) SimpleMkdirAll(perms uint32, fsINodePath string, uid uint32, gid uint32) fsapitypes.Status {
	var (
		paths       []string
		i           int
		parentID    types.FsINodeID = types.RootFsINodeID
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

		code = p.SimpleMkdir(&fsINodeMeta, nil, parentID, perms, paths[i], uid, gid, types.FS_RDEV)
		if code != fsapitypes.OK && code != types.FS_EEXIST {
			goto DONE
		}
		parentID = fsINodeMeta.Ino
	}

DONE:
	return code
}

func (p *DirTreeStg) Mkdir(input *fsapitypes.MkdirIn, name string, out *fsapitypes.EntryOut) fsapitypes.Status {
	var (
		fsINodeMeta sdfsapitypes.FsINodeMeta
		code        fsapitypes.Status
		err         error
	)

	code = p.SimpleMkdir(&fsINodeMeta, nil, input.NodeId, input.Mode, name, input.Uid, input.Gid, types.FS_RDEV)
	if code != fsapitypes.OK {
		return code
	}

	p.SetFsEntryOutByFsINode(out, &fsINodeMeta)

	err = p.RefreshFsINodeACMtimeByIno(fsINodeMeta.ParentID)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	return fsapitypes.OK
}

func (p *DirTreeStg) SimpleMkdir(pFsINodeMeta *sdfsapitypes.FsINodeMeta,
	fsINodeID *types.FsINodeID, parentID types.FsINodeID,
	perms uint32, name string,
	uid uint32, gid uint32, rdev uint32) fsapitypes.Status {
	if len([]byte(name)) > types.FS_MAX_NAME_LENGTH {
		return types.FS_ENAMETOOLONG
	}

	var (
		uFsINode sdfsapitypes.FsINodeUintptr
		err      error
	)

	uFsINode, err = p.FsINodeDriver.GetFsINodeByName(parentID, name)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err == nil {
		return types.FS_EEXIST
	}

	if err != nil && err != types.ErrObjectNotExists {
		return types.ErrorToFsStatus(err)
	}

	err = p.createFsINode(pFsINodeMeta,
		fsINodeID, nil, parentID,
		name, types.FSINODE_TYPE_DIR, fsapitypes.S_IFDIR|perms,
		uid, gid, rdev)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	return fsapitypes.OK
}

func (p *DirTreeStg) checkIsDirEmpty(pFsINodeMeta *sdfsapitypes.FsINodeMeta) (bool, error) {
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

func (p *DirTreeStg) Rmdir(header *fsapitypes.InHeader, name string) fsapitypes.Status {
	var (
		fsINodeMeta sdfsapitypes.FsINodeMeta
		isDirEmpty  bool
		err         error
	)

	err = p.FetchFsINodeByName(&fsINodeMeta, header.NodeId, name)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	if fsINodeMeta.Type != types.FSINODE_TYPE_DIR {
		return fsapitypes.ENOTDIR
	}

	isDirEmpty, err = p.checkIsDirEmpty(&fsINodeMeta)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	if isDirEmpty == false {
		return types.FS_ENOTEMPTY
	}

	err = p.FsINodeDriver.UnlinkFsINode(fsINodeMeta.Ino)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(header.NodeId)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	return fsapitypes.OK
}

func (p *DirTreeStg) OpenDir(input *fsapitypes.OpenIn, out *fsapitypes.OpenOut) fsapitypes.Status {
	var (
		fsINodeMeta sdfsapitypes.FsINodeMeta
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, input.NodeId)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.SimpleOpen(&fsINodeMeta, input.Flags, out)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}
	return fsapitypes.OK
}

func (p *DirTreeStg) ReadDir(input *fsapitypes.ReadIn, out *fsapitypes.DirEntryList) fsapitypes.Status {
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
			if fsINodeMeta.Type == types.FSINODE_TYPE_HARD_LINK {
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
		return types.ErrorToFsStatus(err)
	}

	return fsapitypes.OK
}

func (p *DirTreeStg) ReadDirPlus(input *fsapitypes.ReadIn, out *fsapitypes.DirEntryList) fsapitypes.Status {
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
		return types.ErrorToFsStatus(err)
	}
	return fsapitypes.OK
}

func (p *DirTreeStg) ReleaseDir(input *fsapitypes.ReleaseIn) {
	// TODO make sure releaable
	p.FdTable.ReleaseFd(input.Fh)
}

func (p *DirTreeStg) FsyncDir(input *fsapitypes.FsyncIn) fsapitypes.Status {
	return fsapitypes.OK
}
