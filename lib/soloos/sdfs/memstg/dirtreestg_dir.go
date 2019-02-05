package memstg

import (
	"path/filepath"
	"soloos/sdfs/types"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
)

func (p *DirTreeStg) ListFsINodeByIno(ino types.FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(types.FsINode) bool,
) error {
	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FsINodeDriver.FetchFsINodeByIDThroughHardLink(ino, &fsINode)
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.helper.ListFsINodeByParentIDFromDB(fsINode.Ino, isFetchAllCols, beforeLiteralFunc, literalFunc)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeStg) ListFsINodeByParentPath(parentPath string,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(types.FsINode) bool,
) error {
	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByPath(parentPath, &fsINode)
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.helper.ListFsINodeByParentIDFromDB(fsINode.Ino, isFetchAllCols, beforeLiteralFunc, literalFunc)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeStg) RenameWithFullPath(oldFsINodeName, newFsINodePath string) error {
	var (
		fsINode                       types.FsINode
		oldFsINode                    types.FsINode
		parentFsINode                 types.FsINode
		tmpFsINode                    types.FsINode
		tmpParentDirPath, tmpFileName string
		err                           error
	)

	err = p.FetchFsINodeByPath(oldFsINodeName, &oldFsINode)
	if err != nil {
		return err
	}
	fsINode = oldFsINode

	tmpParentDirPath, tmpFileName = filepath.Split(newFsINodePath)
	err = p.FetchFsINodeByPath(tmpParentDirPath, &parentFsINode)
	if err != nil {
		return err
	}

	if parentFsINode.Type != types.FSINODE_TYPE_DIR {
		return types.ErrObjectNotExists
	}

	if tmpFileName == "" {
		fsINode.ParentID = parentFsINode.Ino
		// keep fsINode.Name
		goto PREPARE_PARENT_FSINODE_DONE
	}

	err = p.FetchFsINodeByPath(newFsINodePath, &tmpFsINode)
	if err != nil {
		if err == types.ErrObjectNotExists {
			fsINode.ParentID = parentFsINode.Ino
			fsINode.Name = tmpFileName
			goto PREPARE_PARENT_FSINODE_DONE
		} else {
			return types.ErrObjectNotExists
		}
	}

	if tmpFsINode.Type == types.FSINODE_TYPE_DIR {
		parentFsINode = tmpFsINode
		fsINode.ParentID = parentFsINode.Ino
		// keep fsINode.Name
		goto PREPARE_PARENT_FSINODE_DONE
	} else {
		return types.ErrObjectNotExists
	}
PREPARE_PARENT_FSINODE_DONE:

	p.FsINodeDriver.DeleteFsINodeCache(oldFsINode.ParentID, oldFsINode.Name, oldFsINode.Ino)

	err = p.UpdateFsINodeInDB(&fsINode)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeStg) Rename(input *fuse.RenameIn, oldName string, newName string) fuse.Status {
	var (
		oldDirFsINodeID = input.NodeId
		newDirFsINodeID = input.Newdir
		oldFsINode      types.FsINode
		checkFsINode    types.FsINode
		isDirEmpty      bool
		err             error
	)

	if len([]byte(newName)) > types.FS_MAX_NAME_LENGTH {
		return types.FS_ENAMETOOLONG
	}

	err = p.FetchFsINodeByName(oldDirFsINodeID, oldName, &oldFsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.FetchFsINodeByName(newDirFsINodeID, newName, &checkFsINode)
	if err != nil {
		if err != types.ErrObjectNotExists {
			return types.ErrorToFuseStatus(err)
		}
	} else {
		// newName exists
		if checkFsINode.Type == types.FSINODE_TYPE_DIR {
			if oldFsINode.Type == types.FSINODE_TYPE_DIR {
				isDirEmpty, err = p.CheckIsDirEmpty(&checkFsINode)
				if err != nil {
					return types.ErrorToFuseStatus(err)
				}
				if isDirEmpty == false {
					return types.FS_ENOTEMPTY
				}
				err = p.SimpleUnlink(&checkFsINode)
				if err != nil {
					return types.ErrorToFuseStatus(err)
				}

			} else {
				newDirFsINodeID = checkFsINode.Ino
			}
		} else {
			err = p.SimpleUnlink(&checkFsINode)
			if err != nil {
				return types.ErrorToFuseStatus(err)
			}
		}
	}

	oldFsINode.ParentID = newDirFsINodeID
	oldFsINode.Name = newName
	err = p.UpdateFsINodeInDB(&oldFsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	p.FsINodeDriver.DeleteFsINodeCache(oldDirFsINodeID, oldName, oldFsINode.Ino)
	err = p.FsINodeDriver.PrepareAndSetFsINodeCache(&oldFsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	return fuse.OK
}

func (p *DirTreeStg) MkdirAll(perms uint32, fsINodePath string, uid uint32, gid uint32) fuse.Status {
	var (
		paths    []string
		i        int
		parentID types.FsINodeID = types.RootFsINodeID
		fsINode  types.FsINode
		code     fuse.Status
	)

	paths = strings.Split(fsINodePath, "/")

	if paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}

	for i = 1; i < len(paths); i++ {
		if paths[i] == "" {
			continue
		}

		code = p.SimpleMkdir(&fsINode, nil, parentID, perms, paths[i], uid, gid, types.FS_RDEV)
		if code != fuse.OK && code != types.FS_EEXIST {
			goto DONE
		}
		parentID = fsINode.Ino
	}

DONE:
	return code
}

func (p *DirTreeStg) Mkdir(input *fuse.MkdirIn, name string, out *fuse.EntryOut) fuse.Status {
	var (
		fsINode types.FsINode
		code    fuse.Status
		err     error
	)

	code = p.SimpleMkdir(&fsINode, nil, input.NodeId, input.Mode, name, input.Uid, input.Gid, types.FS_RDEV)
	if code != fuse.OK {
		return code
	}

	p.SetFuseEntryOutByFsINode(out, &fsINode)

	err = p.RefreshFsINodeACMtimeByIno(fsINode.ParentID)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	return fuse.OK
}

func (p *DirTreeStg) SimpleMkdir(fsINode *types.FsINode,
	fsINodeID *types.FsINodeID, parentID types.FsINodeID,
	perms uint32, name string,
	uid uint32, gid uint32, rdev uint32) fuse.Status {
	if len([]byte(name)) > types.FS_MAX_NAME_LENGTH {
		return types.FS_ENAMETOOLONG
	}

	var (
		err error
	)

	err = p.FsINodeDriver.FetchFsINodeByName(parentID, name, fsINode)
	if err == nil {
		return types.FS_EEXIST
	}

	if err != nil && err != types.ErrObjectNotExists {
		return types.ErrorToFuseStatus(err)
	}

	err = p.CreateFsINode(fsINode,
		fsINodeID, nil, parentID,
		name, types.FSINODE_TYPE_DIR, fuse.S_IFDIR|perms,
		uid, gid, rdev)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	return fuse.OK
}

func (p *DirTreeStg) CheckIsDirEmpty(fsINode *types.FsINode) (bool, error) {
	var (
		isDirEmpty bool
		err        error
	)

	err = p.ListFsINodeByIno(fsINode.Ino, false,
		func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64) {
			isDirEmpty = (resultCount == 0)
			return 0, 0
		},
		func(types.FsINode) bool {
			return false
		},
	)
	if err != nil {
		return false, err
	}

	return isDirEmpty, nil
}

func (p *DirTreeStg) Rmdir(header *fuse.InHeader, name string) fuse.Status {
	var (
		fsINode    types.FsINode
		isDirEmpty bool
		err        error
	)

	err = p.FetchFsINodeByName(header.NodeId, name, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	if fsINode.Type != types.FSINODE_TYPE_DIR {
		return fuse.ENOTDIR
	}

	isDirEmpty, err = p.CheckIsDirEmpty(&fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	if isDirEmpty == false {
		return types.FS_ENOTEMPTY
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

func (p *DirTreeStg) OpenDir(input *fuse.OpenIn, out *fuse.OpenOut) fuse.Status {
	var (
		fsINode types.FsINode
		err     error
	)
	err = p.FetchFsINodeByIDThroughHardLink(input.NodeId, &fsINode)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	err = p.SimpleOpen(&fsINode, input.Flags, out)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}
	return fuse.OK
}

func (p *DirTreeStg) ReadDir(input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	var (
		fsINodeByIDThroughHardLink types.FsINode
		isAddDirEntrySuccess       bool
		err                        error
	)
	err = p.ListFsINodeByIno(input.NodeId, false,
		func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64) {
			return uint64(resultCount) - input.Offset, input.Offset
		},
		func(fsINode types.FsINode) bool {
			if fsINode.Type == types.FSINODE_TYPE_HARD_LINK {
				err = p.FetchFsINodeByIDThroughHardLink(fsINode.Ino, &fsINodeByIDThroughHardLink)
				if err != nil {
					return false
				}
				isAddDirEntrySuccess, _ = out.AddDirEntry(fuse.DirEntry{
					Mode: fsINodeByIDThroughHardLink.Mode,
					Name: fsINode.Name,
					Ino:  fsINodeByIDThroughHardLink.Ino,
				})
			} else {
				isAddDirEntrySuccess, _ = out.AddDirEntry(fuse.DirEntry{
					Mode: fsINode.Mode,
					Name: fsINode.Name,
					Ino:  fsINode.Ino,
				})

			}
			return isAddDirEntrySuccess
		},
	)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}

	return fuse.OK
}

func (p *DirTreeStg) ReadDirPlus(input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	var (
		fsINodeByIDThroughHardLink types.FsINode
		entryOut                   *fuse.EntryOut
		off                        uint64
		err                        error
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
		func(fsINode types.FsINode) bool {
			err = p.FetchFsINodeByIDThroughHardLink(fsINode.Ino, &fsINodeByIDThroughHardLink)
			if err != nil {
				return false
			}

			entryOut, off = out.AddDirLookupEntry(fuse.DirEntry{
				Mode: fsINodeByIDThroughHardLink.Mode,
				Name: fsINode.Name,
				Ino:  fsINodeByIDThroughHardLink.Ino,
			})
			if entryOut == nil {
				return false
			}

			p.SetFuseEntryOutByFsINode(entryOut, &fsINodeByIDThroughHardLink)
			return true
		},
	)
	if err != nil {
		return types.ErrorToFuseStatus(err)
	}
	return fuse.OK
}

func (p *DirTreeStg) ReleaseDir(input *fuse.ReleaseIn) {
	// TODO make sure releaable
	p.FdTable.ReleaseFd(input.Fh)
}

func (p *DirTreeStg) FsyncDir(input *fuse.FsyncIn) fuse.Status {
	return fuse.OK
}
