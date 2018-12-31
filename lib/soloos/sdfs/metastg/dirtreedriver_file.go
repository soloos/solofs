package metastg

import (
	"soloos/sdfs/types"
	"strings"
)

func (p *DirTreeDriver) OpenFile(fsInodePath string) (types.FsINode, error) {
	var (
		paths    []string
		i        int
		parentID types.FsINodeID
		fsINode  types.FsINode
		err      error
	)

	paths = strings.Split(fsInodePath, "/")

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
	if err == types.ErrObjectNotExists {
		fsINode = types.FsINode{
			ID:         p.AllocFsINodeID(),
			ParentID:   parentID,
			Name:       paths[i],
			Flag:       0,
			Permission: 0777,
			NetINodeID: types.ZeroNetINodeID,
			Type:       types.FSINODE_TYPE_FILE,
		}
		err = p.InsertFsINodeInDB(fsINode)
		return fsINode, err
	}

	return fsINode, err
}
