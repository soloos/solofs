package memstg

import (
	"path/filepath"
	"soloos/sdfs/types"
	"strings"
)

func (p *DirTreeStg) DeleteFsINodeByPath(fsINodePath string) error {
	var (
		fsINodeMeta types.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByPath(&fsINodeMeta, fsINodePath)
	if err != nil {
		if err == types.ErrObjectNotExists {
			return nil
		} else {
			return err
		}
	}

	err = p.FsINodeDriver.UnlinkFsINode(fsINodeMeta.Ino)

	return err
}

func (p *DirTreeStg) RenameWithFullPath(oldFsINodeName, newFsINodePath string) error {
	var (
		fsINodeMeta                   types.FsINodeMeta
		oldFsINodeMeta                types.FsINodeMeta
		parentFsINodeMeta             types.FsINodeMeta
		tmpFsINodeMeta                types.FsINodeMeta
		tmpParentDirPath, tmpFileName string
		err                           error
	)

	err = p.FetchFsINodeByPath(&oldFsINodeMeta, oldFsINodeName)
	if err != nil {
		return err
	}
	fsINodeMeta = oldFsINodeMeta

	tmpParentDirPath, tmpFileName = filepath.Split(newFsINodePath)
	err = p.FetchFsINodeByPath(&parentFsINodeMeta, tmpParentDirPath)
	if err != nil {
		return err
	}

	if parentFsINodeMeta.Type != types.FSINODE_TYPE_DIR {
		return types.ErrObjectNotExists
	}

	if tmpFileName == "" {
		fsINodeMeta.ParentID = parentFsINodeMeta.Ino
		// keep fsINodeMeta.Name
		goto PREPARE_PARENT_FSINODE_DONE
	}

	err = p.FetchFsINodeByPath(&tmpFsINodeMeta, newFsINodePath)
	if err != nil {
		if err == types.ErrObjectNotExists {
			fsINodeMeta.ParentID = parentFsINodeMeta.Ino
			fsINodeMeta.SetName(tmpFileName)
			goto PREPARE_PARENT_FSINODE_DONE
		} else {
			return types.ErrObjectNotExists
		}
	}

	if tmpFsINodeMeta.Type == types.FSINODE_TYPE_DIR {
		parentFsINodeMeta = tmpFsINodeMeta
		fsINodeMeta.ParentID = parentFsINodeMeta.Ino
		// keep fsINodeMeta.Name
		goto PREPARE_PARENT_FSINODE_DONE
	} else {
		return types.ErrObjectNotExists
	}
PREPARE_PARENT_FSINODE_DONE:

	err = p.FsINodeDriver.UpdateFsINodeInDB(&fsINodeMeta)
	if err != nil {
		return err
	}

	p.FsINodeDriver.CleanFsINodeAssitCache(oldFsINodeMeta.ParentID, oldFsINodeMeta.Name())

	return nil
}

func (p *DirTreeStg) ListFsINodeByParentPath(parentPath string,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(types.FsINodeMeta) bool,
) error {
	var (
		fsINodeMeta types.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByPath(&fsINodeMeta, parentPath)
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.helper.ListFsINodeByParentIDFromDB(fsINodeMeta.Ino,
		isFetchAllCols, beforeLiteralFunc, literalFunc)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeStg) FetchFsINodeByPath(fsINodeMeta *types.FsINodeMeta, fsINodePath string) error {
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
		*fsINodeMeta = p.FsINodeDriver.RootFsINode.Ptr().Meta
		return nil
	}

	for i = 1; i < len(paths); i++ {
		if paths[i] == "" {
			continue
		}
		err = p.FetchFsINodeByName(fsINodeMeta, parentID, paths[i])
		if err != nil {
			return err
		}
		parentID = fsINodeMeta.Ino
	}

	return err
}
