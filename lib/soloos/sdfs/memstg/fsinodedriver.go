package memstg

import (
	"soloos/log"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/timer"
	"soloos/util"
	"soloos/util/offheap"
	"strings"
	"sync"
	"time"
)

type FsINodeDriverHelper struct {
	api.AllocFsINodeID
	api.GetNetINodeWithReadAcquire
	api.MustGetNetINodeWithReadAcquire
	api.DeleteFsINodeByIDInDB
	api.ListFsINodeByParentIDFromDB
	api.UpdateFsINodeInDB
	api.InsertFsINodeInDB
	api.GetFsINodeByIDFromDB
	api.GetFsINodeByNameFromDB
}

type FsINodeDriver struct {
	Timer timer.Timer

	dirTreeStg *DirTreeStg
	helper     FsINodeDriverHelper

	fsINodesRWMutex sync.RWMutex
	fsINodesPool    types.FsINodePool
	fsINodesByID    map[types.FsINodeID]types.FsINode
	fsINodesByPath  map[string]types.FsINode

	SysFsINode  [2]types.FsINode
	RootFsINode types.FsINode

	EntryTtl           time.Duration
	EntryAttrValid     uint64
	EntryAttrValidNsec uint32

	INodeRWMutexPool types.INodeRWMutexPool

	FIXAttrDriver FIXAttrDriver
}

func (p *FsINodeDriver) Init(
	offheapDriver *offheap.OffheapDriver,
	dirTreeStg *DirTreeStg,
	allocFsINodeID api.AllocFsINodeID,
	getNetINodeWithReadAcquire api.GetNetINodeWithReadAcquire,
	mustGetNetINodeWithReadAcquire api.MustGetNetINodeWithReadAcquire,
	deleteFsINodeByIDInDB api.DeleteFsINodeByIDInDB,
	listFsINodeByParentIDFromDB api.ListFsINodeByParentIDFromDB,
	updateFsINodeInDB api.UpdateFsINodeInDB,
	insertFsINodeInDB api.InsertFsINodeInDB,
	getFsINodeByIDFromDB api.GetFsINodeByIDFromDB,
	getFsINodeByNameFromDB api.GetFsINodeByNameFromDB,
	// FIXAttrDriver
	deleteFIXAttrInDB api.DeleteFIXAttrInDB,
	replaceFIXAttrInDB api.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB api.GetFIXAttrByInoFromDB,
) error {
	var err error

	err = p.Timer.Init()
	if err != nil {
		return err
	}

	p.dirTreeStg = dirTreeStg
	p.SetHelper(
		allocFsINodeID,
		getNetINodeWithReadAcquire,
		mustGetNetINodeWithReadAcquire,
		deleteFsINodeByIDInDB,
		listFsINodeByParentIDFromDB,
		updateFsINodeInDB,
		insertFsINodeInDB,
		getFsINodeByIDFromDB,
		getFsINodeByNameFromDB,
	)

	p.fsINodesByID = make(map[types.FsINodeID]types.FsINode)
	p.fsINodesByPath = make(map[string]types.FsINode)

	err = p.prepareBaseDir()
	if err != nil {
		return err
	}

	p.EntryTtl = 3 * time.Second
	SplitDuration(p.EntryTtl, &p.EntryAttrValid, &p.EntryAttrValidNsec)

	p.INodeRWMutexPool.Init(offheapDriver)

	err = p.FIXAttrDriver.Init(
		deleteFIXAttrInDB,
		replaceFIXAttrInDB,
		getFIXAttrByInoFromDB,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) SetHelper(
	allocFsINodeID api.AllocFsINodeID,
	getNetINodeWithReadAcquire api.GetNetINodeWithReadAcquire,
	mustGetNetINodeWithReadAcquire api.MustGetNetINodeWithReadAcquire,
	deleteFsINodeByIDInDB api.DeleteFsINodeByIDInDB,
	listFsINodeByParentIDFromDB api.ListFsINodeByParentIDFromDB,
	updateFsINodeInDB api.UpdateFsINodeInDB,
	insertFsINodeInDB api.InsertFsINodeInDB,
	getFsINodeByIDFromDB api.GetFsINodeByIDFromDB,
	getFsINodeByNameFromDB api.GetFsINodeByNameFromDB,
) {
	p.helper = FsINodeDriverHelper{
		AllocFsINodeID:                 allocFsINodeID,
		GetNetINodeWithReadAcquire:     getNetINodeWithReadAcquire,
		MustGetNetINodeWithReadAcquire: mustGetNetINodeWithReadAcquire,
		DeleteFsINodeByIDInDB:          deleteFsINodeByIDInDB,
		ListFsINodeByParentIDFromDB:    listFsINodeByParentIDFromDB,
		UpdateFsINodeInDB:              updateFsINodeInDB,
		InsertFsINodeInDB:              insertFsINodeInDB,
		GetFsINodeByIDFromDB:           getFsINodeByIDFromDB,
		GetFsINodeByNameFromDB:         getFsINodeByNameFromDB,
	}
}

func (p *FsINodeDriver) prepareBaseDir() error {
	var (
		fsINode types.FsINode
		ino     types.FsINodeID
		err     error
	)

	ino = types.RootFsINodeID
	err = p.dirTreeStg.Mkdir(&ino, types.RootFsINodeParentID, 0777, "", &fsINode)
	if err != nil {
		log.Warn(err)
	}

	ino = p.helper.AllocFsINodeID()
	err = p.dirTreeStg.Mkdir(&ino, types.RootFsINodeID, 0777, "tmp", &fsINode)
	if err != nil {
		log.Warn(err)
	}

	err = p.FetchFsINodeByName(types.RootFsINodeParentID, "", &p.RootFsINode)
	if err != nil {
		return err
	}

	for i := 0; i < len(p.SysFsINode); i++ {
		p.SysFsINode[i].Ino = types.FsINodeID(i)
	}

	return nil
}

func (p *FsINodeDriver) checkIfNeedNetINode(fsINodeType int) bool {
	return fsINodeType == types.FSINODE_TYPE_FILE
}

func (p *FsINodeDriver) ensureFsINodeHasNetINode(fsINode *types.FsINode) error {
	if p.checkIfNeedNetINode(fsINode.Type) == false {
		return nil
	}

	if fsINode.UNetINode != 0 {
		return nil
	}

	var err error
	fsINode.UNetINode, err = p.helper.GetNetINodeWithReadAcquire(fsINode.NetINodeID)
	return err
}

// ensureFsINodeValidInCache return false if fsinode invalid in cache
// if fsinode invalid, delete cache
func (p *FsINodeDriver) ensureFsINodeValidInCache(fsINode *types.FsINode) bool {
	if p.Timer.Now().Unix()-fsINode.LoadInMemAt < int64(p.EntryAttrValid) {
		return true
	}

	p.DeleteFsINodeCache(fsINode.ParentID, fsINode.Name, fsINode.Ino)
	return false
}

func (p *FsINodeDriver) PrepareAndSetFsINodeCache(fsINode *types.FsINode) error {
	var err error
	err = p.ensureFsINodeHasNetINode(fsINode)
	if err != nil {
		return err
	}

	p.SetFsINodeCache(fsINode)
	return nil
}

func (p *FsINodeDriver) SetFsINodeCache(fsINode *types.FsINode) {
	fsINode.LoadInMemAt = p.Timer.Now().Unix()

	p.fsINodesRWMutex.Lock()
	p.fsINodesByPath[p.fsINodesPool.MakeFsINodeKey(fsINode.ParentID, fsINode.Name)] = *fsINode
	p.fsINodesByID[fsINode.Ino] = *fsINode
	p.fsINodesRWMutex.Unlock()
}

func (p *FsINodeDriver) DeleteFsINodeCache(parentID types.FsINodeID, fsINodeName string, fsINodeID types.FsINodeID) {
	p.fsINodesRWMutex.Lock()
	delete(p.fsINodesByPath, p.fsINodesPool.MakeFsINodeKey(parentID, fsINodeName))
	delete(p.fsINodesByID, fsINodeID)
	p.fsINodesRWMutex.Unlock()
}

func (p *FsINodeDriver) DeleteFsINodeByPath(fsINodePath string) error {
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

	err = p.helper.DeleteFsINodeByIDInDB(fsINode.Ino)
	p.DeleteFsINodeCache(fsINode.ParentID, fsINode.Name, fsINode.Ino)

	return err
}

func (p *FsINodeDriver) DeleteFsINodeByIno(ino types.FsINodeID) error {
	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByID(ino, &fsINode)
	if err != nil {
		if err == types.ErrObjectNotExists {
			return nil
		} else {
			return err
		}
	}

	err = p.helper.DeleteFsINodeByIDInDB(fsINode.Ino)
	p.DeleteFsINodeCache(fsINode.ParentID, fsINode.Name, fsINode.Ino)

	return err
}

func (p *FsINodeDriver) FetchFsINodeByPath(fsINodePath string, fsINode *types.FsINode) error {
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
		*fsINode = p.RootFsINode
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

func (p *FsINodeDriver) FetchFsINodeByID(fsINodeID types.FsINodeID, fsINode *types.FsINode) error {
	var (
		exists bool
		err    error
	)

	p.fsINodesRWMutex.RLock()
	*fsINode, exists = p.fsINodesByID[fsINodeID]
	p.fsINodesRWMutex.RUnlock()
	if exists && p.ensureFsINodeValidInCache(fsINode) == true {
		return nil
	}

	if fsINodeID < types.RootFsINodeID {
		*fsINode = p.SysFsINode[fsINodeID]
		return nil
	}

	*fsINode, err = p.helper.GetFsINodeByIDFromDB(fsINodeID)
	if err != nil {
		return err
	}

	err = p.PrepareAndSetFsINodeCache(fsINode)

	return err
}

func (p *FsINodeDriver) FetchFsINodeByName(parentID types.FsINodeID, fsINodeName string, fsINode *types.FsINode) error {
	var (
		exists bool
		err    error
	)

	p.fsINodesRWMutex.RLock()
	*fsINode, exists = p.fsINodesByPath[p.fsINodesPool.MakeFsINodeKey(parentID, fsINodeName)]
	p.fsINodesRWMutex.RUnlock()
	if exists && p.ensureFsINodeValidInCache(fsINode) == true {
		return nil
	}

	*fsINode, err = p.helper.GetFsINodeByNameFromDB(parentID, fsINodeName)
	if err != nil {
		return err
	}

	err = p.PrepareAndSetFsINodeCache(fsINode)

	return err
}

func (p *FsINodeDriver) UpdateFsINodeInDB(pFsINode *types.FsINode) error {
	var err error
	pFsINode.Mtime = types.DirTreeTime(p.Timer.Now().Unix())
	err = p.helper.UpdateFsINodeInDB(*pFsINode)
	p.DeleteFsINodeCache(pFsINode.ParentID, pFsINode.Name, pFsINode.Ino)
	return err
}

func (p *FsINodeDriver) AllocNetINodeID(fsINode *types.FsINode) {
	//TODO
	util.InitUUID64(&fsINode.NetINodeID)
}

func (p *FsINodeDriver) PrepareFsINodeForCreate(fsINode *types.FsINode,
	netINodeID *types.NetINodeID,
	parentID types.FsINodeID,
	name string, fsINodeType int, mode uint32,
) {
	now := p.Timer.Now()
	nowt := types.DirTreeTime(now.Unix())
	nowtnsec := types.DirTreeTimeNsec(now.UnixNano())
	fsINode.Ino = p.helper.AllocFsINodeID()
	if netINodeID == nil {
		p.AllocNetINodeID(fsINode)
	} else {
		fsINode.NetINodeID = *netINodeID
	}
	fsINode.ParentID = parentID
	fsINode.Name = name
	fsINode.Type = fsINodeType
	fsINode.Atime = nowt
	fsINode.Ctime = nowt
	fsINode.Mtime = nowt
	fsINode.Atimensec = nowtnsec
	fsINode.Ctimensec = nowtnsec
	fsINode.Mtimensec = nowtnsec
	fsINode.Mode = mode
	fsINode.Nlink = 1
}

func (p *FsINodeDriver) CreateINode(fsINode *types.FsINode) error {
	var err error
	err = p.helper.InsertFsINodeInDB(*fsINode)
	if err != nil {
		return err
	}

	return nil

}
