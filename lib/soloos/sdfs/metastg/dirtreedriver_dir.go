package metastg

import (
	"soloos/sdfs/types"
	"strings"
)

func (p *DirTreeDriver) ListFsINodeByParentPath(parentPath string, literalFunc func(types.FsINode) bool) error {
	var (
		fsINode types.FsINode
		err     error
	)

	fsINode, err = p.GetFsINodeByIDFromDBByPath(parentPath)
	if err != nil {
		return err
	}

	err = p.ListFsINodeByParentIDFromDB(fsINode.ID, literalFunc)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeDriver) Mkdir(fsInodePath string) (types.FsINode, error) {
	var (
		paths    []string
		fsINode  types.FsINode
		parentID types.FsINodeID
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

	for i = 1; i < len(paths)-1; i++ {
		if paths[i] == "" {
			continue
		}
		fsINode, err = p.GetFsINodeByIDFromDB(parentID, paths[i])
		if err != nil {
			return fsINode, err
		}
		parentID = fsINode.ID
	}

	fsINode, err = p.GetFsINodeByIDFromDB(parentID, paths[i])
	if err == nil {
		return fsINode, nil
	}

	fsINode = types.FsINode{
		ID:         p.AllocFsINodeID(),
		ParentID:   parentID,
		Name:       paths[len(paths)-1],
		Flag:       0,
		Permission: 0777,
		NetINodeID: types.ZeroNetINodeID,
		Type:       types.FSINODE_TYPE_DIR,
	}

	err = p.InsertFsINodeInDB(fsINode)
	return fsINode, err
}
