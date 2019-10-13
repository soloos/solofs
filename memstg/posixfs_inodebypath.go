package memstg

import (
	"path/filepath"
	"soloos/common/solofstypes"
	"strings"
)

func (p *PosixFs) DeleteFsINodeByPath(fsINodePath string) error {
	var (
		fsINodeMeta solofstypes.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByPath(&fsINodeMeta, fsINodePath)
	if err != nil {
		if err.Error() == solofstypes.ErrObjectNotExists.Error() {
			return nil
		} else {
			return err
		}
	}

	err = p.FsINodeDriver.UnlinkFsINode(fsINodeMeta.Ino)

	return err
}

func (p *PosixFs) RenameWithFullPath(oldFsINodeName, newFsINodePath string) error {
	var (
		fsINodeMeta                   solofstypes.FsINodeMeta
		oldFsINodeMeta                solofstypes.FsINodeMeta
		parentFsINodeMeta             solofstypes.FsINodeMeta
		tmpFsINodeMeta                solofstypes.FsINodeMeta
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

	if parentFsINodeMeta.Type != solofstypes.FSINODE_TYPE_DIR {
		return solofstypes.ErrObjectNotExists
	}

	if tmpFileName == "" {
		fsINodeMeta.ParentID = parentFsINodeMeta.Ino
		// keep fsINodeMeta.Name
		goto PREPARE_PARENT_FSINODE_DONE
	}

	err = p.FetchFsINodeByPath(&tmpFsINodeMeta, newFsINodePath)
	if err != nil {
		if err.Error() == solofstypes.ErrObjectNotExists.Error() {
			fsINodeMeta.ParentID = parentFsINodeMeta.Ino
			fsINodeMeta.SetName(tmpFileName)
			goto PREPARE_PARENT_FSINODE_DONE
		} else {
			return solofstypes.ErrObjectNotExists
		}
	}

	if tmpFsINodeMeta.Type == solofstypes.FSINODE_TYPE_DIR {
		parentFsINodeMeta = tmpFsINodeMeta
		fsINodeMeta.ParentID = parentFsINodeMeta.Ino
		// keep fsINodeMeta.Name
		goto PREPARE_PARENT_FSINODE_DONE
	} else {
		return solofstypes.ErrObjectNotExists
	}
PREPARE_PARENT_FSINODE_DONE:

	err = p.FsINodeDriver.UpdateFsINode(&fsINodeMeta)
	if err != nil {
		return err
	}

	p.FsINodeDriver.CleanFsINodeAssitCache(oldFsINodeMeta.ParentID, oldFsINodeMeta.Name())

	return nil
}

func (p *PosixFs) ListFsINodeByParentPath(parentPath string,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int64) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(solofstypes.FsINodeMeta) bool,
) error {
	var (
		fsINodeMeta solofstypes.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByPath(&fsINodeMeta, parentPath)
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.helper.ListFsINodeByParentIDFromDB(p.NameSpaceID,
		fsINodeMeta.Ino,
		isFetchAllCols, beforeLiteralFunc, literalFunc)
	if err != nil {
		return err
	}

	return nil
}

func (p *PosixFs) FetchFsINodeByPath(fsINodeMeta *solofstypes.FsINodeMeta, fsINodePath string) error {
	var (
		paths    []string
		i        int
		parentID solofstypes.FsINodeIno = solofstypes.RootFsINodeIno
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
