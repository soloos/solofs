package memstg

import (
	"path/filepath"
	"soloos/common/sdfsapitypes"
	"soloos/sdfs/sdfstypes"
	"strings"
)

func (p *PosixFS) DeleteFsINodeByPath(fsINodePath string) error {
	var (
		fsINodeMeta sdfsapitypes.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByPath(&fsINodeMeta, fsINodePath)
	if err != nil {
		if err == sdfsapitypes.ErrObjectNotExists {
			return nil
		} else {
			return err
		}
	}

	err = p.FsINodeDriver.UnlinkFsINode(fsINodeMeta.Ino)

	return err
}

func (p *PosixFS) RenameWithFullPath(oldFsINodeName, newFsINodePath string) error {
	var (
		fsINodeMeta                   sdfsapitypes.FsINodeMeta
		oldFsINodeMeta                sdfsapitypes.FsINodeMeta
		parentFsINodeMeta             sdfsapitypes.FsINodeMeta
		tmpFsINodeMeta                sdfsapitypes.FsINodeMeta
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

	if parentFsINodeMeta.Type != sdfstypes.FSINODE_TYPE_DIR {
		return sdfsapitypes.ErrObjectNotExists
	}

	if tmpFileName == "" {
		fsINodeMeta.ParentID = parentFsINodeMeta.Ino
		// keep fsINodeMeta.Name
		goto PREPARE_PARENT_FSINODE_DONE
	}

	err = p.FetchFsINodeByPath(&tmpFsINodeMeta, newFsINodePath)
	if err != nil {
		if err == sdfsapitypes.ErrObjectNotExists {
			fsINodeMeta.ParentID = parentFsINodeMeta.Ino
			fsINodeMeta.SetName(tmpFileName)
			goto PREPARE_PARENT_FSINODE_DONE
		} else {
			return sdfsapitypes.ErrObjectNotExists
		}
	}

	if tmpFsINodeMeta.Type == sdfstypes.FSINODE_TYPE_DIR {
		parentFsINodeMeta = tmpFsINodeMeta
		fsINodeMeta.ParentID = parentFsINodeMeta.Ino
		// keep fsINodeMeta.Name
		goto PREPARE_PARENT_FSINODE_DONE
	} else {
		return sdfsapitypes.ErrObjectNotExists
	}
PREPARE_PARENT_FSINODE_DONE:

	err = p.FsINodeDriver.UpdateFsINodeInDB(&fsINodeMeta)
	if err != nil {
		return err
	}

	p.FsINodeDriver.CleanFsINodeAssitCache(oldFsINodeMeta.ParentID, oldFsINodeMeta.Name())

	return nil
}

func (p *PosixFS) ListFsINodeByParentPath(parentPath string,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(sdfsapitypes.FsINodeMeta) bool,
) error {
	var (
		fsINodeMeta sdfsapitypes.FsINodeMeta
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

func (p *PosixFS) FetchFsINodeByPath(fsINodeMeta *sdfsapitypes.FsINodeMeta, fsINodePath string) error {
	var (
		paths    []string
		i        int
		parentID sdfsapitypes.FsINodeID = sdfsapitypes.RootFsINodeID
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
