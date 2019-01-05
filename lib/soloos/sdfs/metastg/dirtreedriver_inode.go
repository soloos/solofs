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

	p.fsINodesByPath = make(map[string]types.FsINode)
	p.fsINodesByID = make(map[types.FsINodeID]types.FsINode)

	p.rootFsINode, err = p.GetFsINodeByPathFromDB(-1, "")
	if err != nil {
		return err
	}

	p.Mkdir("/tmp")

	return nil
}

func (p *DirTreeDriver) ensureFsINodeHasNetINode(fsINode *types.FsINode) error {
	if fsINode.Type != types.FSINODE_TYPE_FILE {
		return nil
	}

	if fsINode.UNetINode != 0 {
		return nil
	}

	var err error
	fsINode.UNetINode, err = p.helper.GetNetINode(fsINode.NetINodeID)
	return err
}

func (p *DirTreeDriver) prepareAndSetFsINodeCache(fsINode *types.FsINode) error {
	var err error
	err = p.ensureFsINodeHasNetINode(fsINode)
	if err != nil {
		return err
	}

	p.fsINodesByPathRWMutex.Lock()
	p.fsINodesByPath[p.MakeFsINodeKey(fsINode.ParentID, fsINode.Name)] = *fsINode
	p.fsINodesByPathRWMutex.Unlock()

	p.fsINodesByIDRWMutex.Lock()
	p.fsINodesByID[fsINode.ID] = *fsINode
	p.fsINodesByIDRWMutex.Unlock()

	return nil
}

func (p *DirTreeDriver) deleteFsINodeCache(parentID types.FsINodeID, fsINodeName string, fsINodeID types.FsINodeID) {
	p.fsINodesByPathRWMutex.Lock()
	delete(p.fsINodesByPath, p.MakeFsINodeKey(parentID, fsINodeName))
	p.fsINodesByPathRWMutex.Unlock()

	p.fsINodesByIDRWMutex.Lock()
	delete(p.fsINodesByID, fsINodeID)
	p.fsINodesByIDRWMutex.Unlock()
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

	fsINode, err = p.GetFsINodeByPath(fsINodePath)
	if err != nil {
		if err == types.ErrObjectNotExists {
			return nil
		} else {
			return err
		}
	}

	err = p.DeleteFsINodeByIDInDB(fsINode.ID)
	p.deleteFsINodeCache(fsINode.ParentID, fsINode.Name, fsINode.ID)

	return err
}

func (p *DirTreeDriver) GetFsINodeByID(fsINodeID types.FsINodeID) (types.FsINode, error) {
	var (
		fsINode types.FsINode
		err     error
	)
	fsINode, err = p.GetFsINodeByIDFromDB(fsINodeID)
	if err != nil {
		return fsINode, err
	}

	return fsINode, err
}

func (p *DirTreeDriver) GetFsINodeByPath(fsInodePath string) (types.FsINode, error) {
	var (
		paths    []string
		i        int
		parentID types.FsINodeID = p.rootFsINode.ID
		fsINode  types.FsINode
		err      error
	)

	paths = strings.Split(fsInodePath, "/")

	if paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}

	if len(paths) <= 1 {
		return p.rootFsINode, nil
	}

	for i = 1; i < len(paths); i++ {
		if paths[i] == "" {
			continue
		}
		fsINode, err = p.GetFsINodeByPathFromDB(parentID, paths[i])
		if err != nil {
			return fsINode, err
		}
		parentID = fsINode.ID
	}

	return fsINode, err
}
