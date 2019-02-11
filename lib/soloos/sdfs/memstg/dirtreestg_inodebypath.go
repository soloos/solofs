package memstg

import (
	"path/filepath"
	"soloos/sdfs/types"
	"strings"
)

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
