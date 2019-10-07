package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/solofstypes"
	"strings"
)

func (p *PosixFs) ListFsINodeByIno(ino solofstypes.FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int64) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(solofstypes.FsINodeMeta) bool,
) error {
	var (
		uFsINode solofstypes.FsINodeUintptr
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

func (p *PosixFs) Rename(input *fsapi.RenameIn, oldName string, newName string) fsapi.Status {
	var (
		oldDirFsINodeID = input.NodeId
		newDirFsINodeID = input.Newdir
		uOldFsINode     solofstypes.FsINodeUintptr
		pOldFsINode     *solofstypes.FsINode
		uCheckFsINode   solofstypes.FsINodeUintptr
		pCheckFsINode   *solofstypes.FsINode
		isDirEmpty      bool
		err             error
	)

	if len([]byte(newName)) > solofstypes.FS_MAX_NAME_LENGTH {
		return solofstypes.FS_ENAMETOOLONG
	}

	uOldFsINode, err = p.FsINodeDriver.GetFsINodeByName(oldDirFsINodeID, oldName)
	defer p.FsINodeDriver.ReleaseFsINode(uOldFsINode)
	pOldFsINode = uOldFsINode.Ptr()
	if err != nil {
		return ErrorToFsStatus(err)
	}

	uCheckFsINode, err = p.FsINodeDriver.GetFsINodeByName(newDirFsINodeID, newName)
	defer p.FsINodeDriver.ReleaseFsINode(uCheckFsINode)
	pCheckFsINode = uCheckFsINode.Ptr()
	if err != nil {
		if err.Error() != solofstypes.ErrObjectNotExists.Error() {
			return ErrorToFsStatus(err)
		}
	} else {
		// newName exists
		if pCheckFsINode.Meta.Type == solofstypes.FSINODE_TYPE_DIR {
			if pOldFsINode.Meta.Type == solofstypes.FSINODE_TYPE_DIR {
				isDirEmpty, err = p.checkIsDirEmpty(&pCheckFsINode.Meta)
				if err != nil {
					return ErrorToFsStatus(err)
				}
				if isDirEmpty == false {
					return solofstypes.FS_ENOTEMPTY
				}
				err = p.FsINodeDriver.UnlinkFsINode(pCheckFsINode.Meta.Ino)
				if err != nil {
					return ErrorToFsStatus(err)
				}

			} else {
				newDirFsINodeID = pCheckFsINode.Meta.Ino
			}
		} else {
			err = p.FsINodeDriver.UnlinkFsINode(pCheckFsINode.Meta.Ino)
			if err != nil {
				return ErrorToFsStatus(err)
			}
		}
	}

	pOldFsINode.Meta.ParentID = newDirFsINodeID
	pOldFsINode.Meta.SetName(newName)
	err = p.FsINodeDriver.UpdateFsINode(&pOldFsINode.Meta)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	p.FsINodeDriver.CleanFsINodeAssitCache(oldDirFsINodeID, oldName)

	return fsapi.OK
}

func (p *PosixFs) SimpleMkdirAll(perms uint32, fsINodePath string, uid uint32, gid uint32) fsapi.Status {
	var (
		paths       []string
		i           int
		parentID    solofstypes.FsINodeID = solofstypes.RootFsINodeID
		fsINodeMeta solofstypes.FsINodeMeta
		code        fsapi.Status
	)

	paths = strings.Split(fsINodePath, "/")

	if paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}

	for i = 1; i < len(paths); i++ {
		if paths[i] == "" {
			continue
		}

		code = p.SimpleMkdir(&fsINodeMeta, nil, parentID, perms, paths[i], uid, gid, solofstypes.FS_RDEV)
		if code != fsapi.OK && code != solofstypes.FS_EEXIST {
			goto DONE
		}
		parentID = fsINodeMeta.Ino
	}

DONE:
	if code == solofstypes.FS_EEXIST {
		code = fsapi.OK
	}
	return code
}

func (p *PosixFs) Mkdir(input *fsapi.MkdirIn, name string, out *fsapi.EntryOut) fsapi.Status {
	var (
		fsINodeMeta solofstypes.FsINodeMeta
		code        fsapi.Status
		err         error
	)

	code = p.SimpleMkdir(&fsINodeMeta, nil, input.NodeId, input.Mode, name, input.Uid, input.Gid, solofstypes.FS_RDEV)
	if code != fsapi.OK {
		return code
	}

	p.SetFsEntryOutByFsINode(out, &fsINodeMeta)

	err = p.RefreshFsINodeACMtimeByIno(fsINodeMeta.ParentID)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	return fsapi.OK
}

func (p *PosixFs) SimpleMkdir(pFsINodeMeta *solofstypes.FsINodeMeta,
	fsINodeID *solofstypes.FsINodeID, parentID solofstypes.FsINodeID,
	perms uint32, name string,
	uid uint32, gid uint32, rdev uint32) fsapi.Status {
	if len([]byte(name)) > solofstypes.FS_MAX_NAME_LENGTH {
		return solofstypes.FS_ENAMETOOLONG
	}

	var (
		uFsINode solofstypes.FsINodeUintptr
		err      error
	)

	uFsINode, err = p.FsINodeDriver.GetFsINodeByName(parentID, name)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err == nil {
		*pFsINodeMeta = uFsINode.Ptr().Meta
		return solofstypes.FS_EEXIST
	}

	if err != nil && err.Error() != solofstypes.ErrObjectNotExists.Error() {
		return ErrorToFsStatus(err)
	}

	err = p.createFsINode(pFsINodeMeta,
		fsINodeID, nil, parentID,
		name, solofstypes.FSINODE_TYPE_DIR, fsapi.S_IFDIR|perms,
		uid, gid, rdev)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	return fsapi.OK
}

func (p *PosixFs) checkIsDirEmpty(pFsINodeMeta *solofstypes.FsINodeMeta) (bool, error) {
	var (
		isDirEmpty bool
		err        error
	)

	err = p.ListFsINodeByIno(pFsINodeMeta.Ino, false,
		func(resultCount int64) (fetchRowsLimit uint64, fetchRowsOffset uint64) {
			isDirEmpty = (resultCount == 0)
			return 0, 0
		},
		func(solofstypes.FsINodeMeta) bool {
			return false
		},
	)
	if err != nil {
		return false, err
	}

	return isDirEmpty, nil
}

func (p *PosixFs) Rmdir(header *fsapi.InHeader, name string) fsapi.Status {
	var (
		fsINodeMeta solofstypes.FsINodeMeta
		isDirEmpty  bool
		err         error
	)

	err = p.FetchFsINodeByName(&fsINodeMeta, header.NodeId, name)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	if fsINodeMeta.Type != solofstypes.FSINODE_TYPE_DIR {
		return fsapi.ENOTDIR
	}

	isDirEmpty, err = p.checkIsDirEmpty(&fsINodeMeta)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	if isDirEmpty == false {
		return solofstypes.FS_ENOTEMPTY
	}

	err = p.FsINodeDriver.UnlinkFsINode(fsINodeMeta.Ino)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(header.NodeId)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	return fsapi.OK
}

func (p *PosixFs) OpenDir(input *fsapi.OpenIn, out *fsapi.OpenOut) fsapi.Status {
	var (
		fsINodeMeta solofstypes.FsINodeMeta
		err         error
	)
	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, input.NodeId)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.SimpleOpen(&fsINodeMeta, input.Flags, out)
	if err != nil {
		return ErrorToFsStatus(err)
	}
	return fsapi.OK
}

func (p *PosixFs) ReadDir(input *fsapi.ReadIn, out *fsapi.DirEntryList) fsapi.Status {
	var (
		fsINodeMetaByIDThroughHardLink solofstypes.FsINodeMeta
		isAddDirEntrySuccess           bool
		err                            error
	)
	err = p.ListFsINodeByIno(input.NodeId, false,
		func(resultCount int64) (fetchRowsLimit uint64, fetchRowsOffset uint64) {
			return uint64(resultCount) - input.Offset, input.Offset
		},
		func(fsINodeMeta solofstypes.FsINodeMeta) bool {
			if fsINodeMeta.Type == solofstypes.FSINODE_TYPE_HARD_LINK {
				err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMetaByIDThroughHardLink, fsINodeMeta.Ino)
				if err != nil {
					return false
				}
				isAddDirEntrySuccess, _ = out.AddDirEntry(fsapi.DirEntry{
					Mode: fsINodeMetaByIDThroughHardLink.Mode,
					Name: fsINodeMeta.Name(),
					Ino:  fsINodeMetaByIDThroughHardLink.Ino,
				})
			} else {
				isAddDirEntrySuccess, _ = out.AddDirEntry(fsapi.DirEntry{
					Mode: fsINodeMeta.Mode,
					Name: fsINodeMeta.Name(),
					Ino:  fsINodeMeta.Ino,
				})

			}
			return isAddDirEntrySuccess
		},
	)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	return fsapi.OK
}

func (p *PosixFs) ReadDirPlus(input *fsapi.ReadIn, out *fsapi.DirEntryList) fsapi.Status {
	var (
		fsINodeMetaByIDThroughHardLink solofstypes.FsINodeMeta
		entryOut                       *fsapi.EntryOut
		err                            error
	)
	err = p.ListFsINodeByIno(input.NodeId, true,
		func(resultCount int64) (uint64, uint64) {
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
		func(fsINodeMeta solofstypes.FsINodeMeta) bool {
			err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMetaByIDThroughHardLink, fsINodeMeta.Ino)
			if err != nil {
				return false
			}

			entryOut, _ = out.AddDirLookupEntry(fsapi.DirEntry{
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
		return ErrorToFsStatus(err)
	}
	return fsapi.OK
}

func (p *PosixFs) ReleaseDir(input *fsapi.ReleaseIn) {
	// TODO make sure releaable
	p.FdTable.ReleaseFd(input.Fh)
}

func (p *PosixFs) FsyncDir(input *fsapi.FsyncIn) fsapi.Status {
	return fsapi.OK
}
