package metastg

import (
	"soloos/sdfs/types"
	"soloos/util"
	"strings"
)

func (p *DirTreeDriver) allocNetINode(netBlockCap int, memBlockCap int) (types.NetINodeID, error) {
	var (
		netINodeID types.NetINodeID
		err        error
	)

	util.InitUUID64(&netINodeID)
	_, err = p.helper.MustGetNetINode(netINodeID, 0, netBlockCap, memBlockCap)
	return netINodeID, err
}

func (p *DirTreeDriver) OpenFile(fsInodePath string, netBlockCap int, memBlockCap int) (types.FsINode, error) {
	var (
		paths      []string
		i          int
		parentID   types.FsINodeID = p.rootFsINode.ID
		fsINode    types.FsINode
		netINodeID types.NetINodeID
		err        error
	)

	paths = strings.Split(fsInodePath, "/")

	if paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}

	for i = 1; i < len(paths)-1; i++ {
		if paths[i] == "" {
			continue
		}
		fsINode, err = p.GetFsINodeByPathFromDB(parentID, paths[i])
		if err != nil {
			goto OPEN_FILE_DONE
		}
		parentID = fsINode.ID
	}

	fsINode, err = p.GetFsINodeByPathFromDB(parentID, paths[i])
	if err == nil {
		goto OPEN_FILE_DONE
	}

	if err == types.ErrObjectNotExists {
		netINodeID, err = p.allocNetINode(netBlockCap, memBlockCap)
		if err != nil {
			goto OPEN_FILE_DONE
		}

		fsINode = types.FsINode{
			ID:         p.AllocFsINodeID(),
			ParentID:   parentID,
			Name:       paths[i],
			Flag:       0,
			Permission: 0777,
			NetINodeID: netINodeID,
			Type:       types.FSINODE_TYPE_FILE,
		}
		err = p.InsertFsINodeInDB(fsINode)
	}

OPEN_FILE_DONE:
	if err == nil {
		err = p.prepareAndSetFsINodeCache(&fsINode)
	}
	return fsINode, err
}
