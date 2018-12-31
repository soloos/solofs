package metastg

import (
	"soloos/sdfs/types"
	"soloos/util"
	"strconv"
	"strings"
	"sync/atomic"
)

func (p *DirTreeDriver) MakeFsINodeKey(parentID types.FsINodeID, fsINodeName string) string {
	return strconv.FormatInt(parentID, 10) + fsINodeName
}

func (p *DirTreeDriver) PrepareINodes() error {
	var err error
	p.allocINodeIDDalta = 10000 * 10
	p.lastFsINodeIDInDB, err = p.helper.FetchAndUpdateMaxID("b_fsinode", p.allocINodeIDDalta)
	if err != nil {
		return err
	}
	p.maxFsINodeID = p.lastFsINodeIDInDB

	p.fsINodes = make(map[string]types.FsINode)

	return nil
}

func (p *DirTreeDriver) AllocFsINodeID() int64 {
	var ret = atomic.AddInt64(&p.maxFsINodeID, 1)
	if p.lastFsINodeIDInDB-ret < p.allocINodeIDDalta/100 {
		util.AssertErrIsNil1(p.helper.FetchAndUpdateMaxID("b_fsinode", p.allocINodeIDDalta))
		p.lastFsINodeIDInDB += p.allocINodeIDDalta
	}
	return ret
}

func (p *DirTreeDriver) DeleteINodeByPath(fsINodePath string) error {
	var (
		fsINode types.FsINode
		err     error
	)

	fsINode, err = p.GetFsINodeByIDFromDBByPath(fsINodePath)
	if err != nil {
		if err == types.ErrObjectNotExists {
			return nil
		} else {
			return err
		}
	}

	p.fsINodesRWMutex.Lock()
	delete(p.fsINodes, p.MakeFsINodeKey(fsINode.ParentID, fsINode.Name))
	err = p.DeleteFsINodeByIDInDB(fsINode.ID)
	p.fsINodesRWMutex.Unlock()

	return err
}

func (p *DirTreeDriver) GetFsINodeByIDFromDBByPath(fsInodePath string) (types.FsINode, error) {
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

	for i = 1; i < len(paths); i++ {
		if paths[i] == "" {
			continue
		}
		fsINode, err = p.GetFsINodeByIDFromDB(parentID, paths[i])
		if err != nil {
			return fsINode, err
		}
		parentID = fsINode.ID
	}

	return fsINode, err
}
