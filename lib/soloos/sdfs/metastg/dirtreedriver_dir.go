package metastg

import (
	"path/filepath"
	"soloos/sdfs/types"
	"strings"
)

func (p *DirTreeDriver) Rename(oldFsINodeName, newFsINodePath string) error {
	var (
		fsINode                       types.FsINode
		oldFsINode                    types.FsINode
		parentFsINode                 types.FsINode
		tmpFsINode                    types.FsINode
		tmpParentDirPath, tmpFileName string
		err                           error
	)

	oldFsINode, err = p.GetFsINodeByPath(oldFsINodeName)
	if err != nil {
		return err
	}
	fsINode = oldFsINode

	tmpParentDirPath, tmpFileName = filepath.Split(newFsINodePath)
	parentFsINode, err = p.GetFsINodeByPath(tmpParentDirPath)
	if err != nil {
		return err
	}

	if parentFsINode.Type != types.FSINODE_TYPE_DIR {
		return types.ErrObjectNotExists
	}

	if tmpFileName == "" {
		fsINode.ParentID = parentFsINode.ID
		// keep fsINode.Name
		goto PREPARE_PARENT_FSINODE_DONE
	}

	tmpFsINode, err = p.GetFsINodeByPath(newFsINodePath)
	if err != nil {
		if err == types.ErrObjectNotExists {
			fsINode.ParentID = parentFsINode.ID
			fsINode.Name = tmpFileName
			goto PREPARE_PARENT_FSINODE_DONE
		} else {
			return types.ErrObjectNotExists
		}
	}

	if tmpFsINode.Type == types.FSINODE_TYPE_DIR {
		parentFsINode = tmpFsINode
		fsINode.ParentID = parentFsINode.ID
		// keep fsINode.Name
		goto PREPARE_PARENT_FSINODE_DONE
	} else {
		return types.ErrObjectNotExists
	}
PREPARE_PARENT_FSINODE_DONE:

	p.deleteFsINodeCache(oldFsINode.ParentID, oldFsINode.Name, oldFsINode.ID)

	err = p.UpdateFsINodeInDB(fsINode)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeDriver) ListFsINodeByParentPath(parentPath string,
	beforeLiteralFunc func(resultCount int) bool,
	literalFunc func(types.FsINode) bool,
) error {
	var (
		fsINode types.FsINode
		err     error
	)

	fsINode, err = p.GetFsINodeByPath(parentPath)
	if err != nil {
		return err
	}

	err = p.ListFsINodeByParentIDFromDB(fsINode.ID, beforeLiteralFunc, literalFunc)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeDriver) Mkdir(fsInodePath string) (types.FsINode, error) {
	var (
		paths    []string
		fsINode  types.FsINode
		parentID types.FsINodeID = p.rootFsINode.ID
		i        int
		err      error
	)

	paths = strings.Split(fsInodePath, "/")
	if len(paths) == 0 {
		return fsINode, types.ErrObjectNotExists
	}

	if paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}

	for i = 1; i < len(paths); i++ {
		if paths[i] == "" {
			continue
		}
		fsINode, err = p.GetFsINodeByPathFromDB(parentID, paths[i])
		if err != nil {
			if err != types.ErrObjectNotExists {
				return fsINode, err
			}

			fsINode = types.FsINode{
				ID:         p.AllocFsINodeID(),
				ParentID:   parentID,
				Name:       paths[i],
				Flag:       0,
				Permission: 0777,
				NetINodeID: types.ZeroNetINodeID,
				Type:       types.FSINODE_TYPE_DIR,
			}
			err = p.InsertFsINodeInDB(fsINode)
			if err != nil {
				return fsINode, err
			}
		}
		parentID = fsINode.ID
	}

	return fsINode, err
}
