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

	err = p.FsINodeDriver.FetchFsINodeByID(ino, &fsINode)
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

	err = p.FsINodeDriver.FetchFsINodeByPath(parentPath, &fsINode)
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.helper.ListFsINodeByParentIDFromDB(fsINode.Ino, isFetchAllCols, beforeLiteralFunc, literalFunc)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeStg) Rename(oldFsINodeName, newFsINodePath string) error {
	var (
		fsINode                       types.FsINode
		oldFsINode                    types.FsINode
		parentFsINode                 types.FsINode
		tmpFsINode                    types.FsINode
		tmpParentDirPath, tmpFileName string
		err                           error
	)

	err = p.FsINodeDriver.FetchFsINodeByPath(oldFsINodeName, &oldFsINode)
	if err != nil {
		return err
	}
	fsINode = oldFsINode

	tmpParentDirPath, tmpFileName = filepath.Split(newFsINodePath)
	err = p.FsINodeDriver.FetchFsINodeByPath(tmpParentDirPath, &parentFsINode)
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

	err = p.FsINodeDriver.FetchFsINodeByPath(newFsINodePath, &tmpFsINode)
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

	err = p.FsINodeDriver.UpdateFsINodeInDB(&fsINode)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeStg) MkdirAll(mode uint32, fsINodePath string) error {
	var (
		paths    []string
		i        int
		parentID types.FsINodeID = types.RootFsINodeID
		fsINode  types.FsINode
		err      error
	)

	paths = strings.Split(fsINodePath, "/")

	if paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}

	for i = 1; i < len(paths); i++ {
		if paths[i] == "" {
			continue
		}

		err = p.Mkdir(nil, parentID, 0777, paths[i], &fsINode)
		if err != nil && err != types.ErrObjectExists {
			goto DONE
		}
		parentID = fsINode.Ino
	}

DONE:
	return err
}

func (p *DirTreeStg) Mkdir(ino *types.FsINodeID,
	parentID types.FsINodeID, perms uint32, name string, fsINode *types.FsINode) error {
	var (
		err error
	)

	err = p.FsINodeDriver.FetchFsINodeByName(parentID, name, fsINode)
	if err == nil {
		return types.ErrObjectExists
	}

	if err != nil && err != types.ErrObjectNotExists {
		return err
	}

	now := p.FsINodeDriver.Timer.Now()
	nowt := types.DirTreeTime(now.Unix())
	nowtnsec := types.DirTreeTimeNsec(now.UnixNano())
	*fsINode = types.FsINode{
		NetINodeID: types.ZeroNetINodeID,
		ParentID:   parentID,
		Name:       name,
		Type:       types.FSINODE_TYPE_DIR,
		Atime:      nowt,
		Ctime:      nowt,
		Mtime:      nowt,
		Atimensec:  nowtnsec,
		Ctimensec:  nowtnsec,
		Mtimensec:  nowtnsec,
		Mode:       fuse.S_IFDIR | perms,
		Nlink:      1,
	}
	if ino == nil {
		fsINode.Ino = p.FsINodeDriver.helper.AllocFsINodeID()
	} else {
		fsINode.Ino = *ino
	}
	err = p.FsINodeDriver.helper.InsertFsINodeInDB(*fsINode)
	if err != nil {
		return err
	}

	return err
}

func (p *DirTreeStg) Rmdir(ino types.FsINodeID) error {
	var (
		isHasChildren = false
		err           error
	)

	err = p.ListFsINodeByIno(ino, false,
		func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64) {
			isHasChildren = resultCount > 0
			return 0, 0
		},
		func(types.FsINode) bool {
			return false
		},
	)
	if err != nil {
		return err
	}

	if isHasChildren {
		return types.ErrObjectHasChildren
	}

	err = p.FsINodeDriver.DeleteFsINodeByIno(ino)
	return err
}
