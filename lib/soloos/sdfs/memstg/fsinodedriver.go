package memstg

import (
	fsapitypes "soloos/common/fsapi/types"
	"soloos/common/log"
	"soloos/common/timer"
	"soloos/common/util"
	"soloos/sdbone/offheap"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"strconv"
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
	fsINodesByID    map[types.FsINodeID]types.FsINode
	fsINodesByPath  map[string]types.FsINode

	SysFsINode  [2]types.FsINode
	RootFsINode types.FsINode

	EntryTtl           time.Duration
	EntryAttrValid     uint64
	EntryAttrValidNsec uint32

	INodeRWMutexTable offheap.HKVTableWithUint64

	FIXAttrDriver FIXAttrDriver

	DefaultNetBlockCap int
	DefaultMemBlockCap int
}

func (p *FsINodeDriver) Init(
	offheapDriver *offheap.OffheapDriver,
	dirTreeStg *DirTreeStg,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
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

	p.DefaultNetBlockCap = defaultNetBlockCap
	p.DefaultMemBlockCap = defaultMemBlockCap

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

	err = offheapDriver.InitHKVTableWithUint64(&p.INodeRWMutexTable, "INodeRWMutex",
		int(types.INodeRWMutexStructSize), -1, types.DefaultKVTableSharedCount,
		nil, nil)
	if err != nil {
		return err
	}

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
		code    fsapitypes.Status
		err     error
	)

	ino = types.RootFsINodeID
	code = p.dirTreeStg.SimpleMkdir(&fsINode, &ino, types.RootFsINodeParentID, 0777, "", 0, 0, types.FS_RDEV)
	if code != fsapitypes.OK {
		log.Warn("mkdir root error ", code)
	}

	ino = p.helper.AllocFsINodeID()
	code = p.dirTreeStg.SimpleMkdir(&fsINode, &ino, types.RootFsINodeID, 0777, "tmp", 0, 0, types.FS_RDEV)
	if code != fsapitypes.OK {
		log.Warn("mkdir tmp error", code)
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
	fsINode.UNetINode, err = p.helper.GetNetINodeWithReadAcquire(true, fsINode.NetINodeID)
	if err != nil {
		return err
	}
	fsINode.UNetINode.Ptr().LastCommitSize = fsINode.UNetINode.Ptr().Size

	return nil
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
	p.fsINodesByPath[p.MakeFsINodeKey(fsINode.ParentID, fsINode.Name)] = *fsINode
	p.fsINodesByID[fsINode.Ino] = *fsINode
	p.fsINodesRWMutex.Unlock()
}

func (p *FsINodeDriver) DeleteFsINodeCache(parentID types.FsINodeID, fsINodeName string, fsINodeID types.FsINodeID) {
	p.fsINodesRWMutex.Lock()
	delete(p.fsINodesByPath, p.MakeFsINodeKey(parentID, fsINodeName))
	delete(p.fsINodesByID, fsINodeID)
	p.fsINodesRWMutex.Unlock()
}

func (p *FsINodeDriver) FetchFsINodeByIDThroughHardLink(fsINodeID types.FsINodeID, fsINode *types.FsINode) error {
	var err error
	for {
		err = p.FetchFsINodeByID(fsINodeID, fsINode)
		if err != nil {
			return err
		}

		if fsINode.Type != types.FSINODE_TYPE_HARD_LINK {
			return nil
		}

		fsINodeID = fsINode.HardLinkIno
	}
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
	*fsINode, exists = p.fsINodesByPath[p.MakeFsINodeKey(parentID, fsINodeName)]
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

func (p *FsINodeDriver) UpdateFsINodeInDB(fsINode *types.FsINode) error {
	var err error
	fsINode.Ctime = types.DirTreeTime(p.Timer.Now().Unix())
	err = p.helper.UpdateFsINodeInDB(*fsINode)
	p.SetFsINodeCache(fsINode)
	return err
}

func (p *FsINodeDriver) RefreshFsINodeACMtime(fsINode *types.FsINode) error {
	var err error
	now := p.Timer.Now()
	nowUnixNano := now.UnixNano()
	if nowUnixNano-fsINode.LastModifyACMTime < int64(time.Millisecond)*600 {
		return nil
	}

	nowt := types.DirTreeTime(now.Unix())
	nowtnsec := types.DirTreeTimeNsec(now.UnixNano())

	fsINode.Atime = nowt
	fsINode.Atimensec = nowtnsec
	fsINode.Ctime = nowt
	fsINode.Ctimensec = nowtnsec
	fsINode.Mtime = nowt
	fsINode.Mtimensec = nowtnsec

	err = p.helper.UpdateFsINodeInDB(*fsINode)
	if err != nil {
		return err
	}

	p.SetFsINodeCache(fsINode)
	fsINode.LastModifyACMTime = nowUnixNano
	return err
}

func (p *FsINodeDriver) RefreshFsINodeACMtimeByIno(fsINodeID types.FsINodeID) error {
	var (
		fsINode types.FsINode
		err     error
	)

	err = p.FetchFsINodeByID(fsINodeID, &fsINode)
	if err != nil {
		return err
	}

	p.RefreshFsINodeACMtime(&fsINode)
	return err
}

func (p *FsINodeDriver) AllocNetINodeID(fsINode *types.FsINode) error {
	var err error
	//TODO improve alloc NetInodeID
	util.InitUUID64(&fsINode.NetINodeID)
	//TODO config memBlockSize netBlockSize
	fsINode.UNetINode, err = p.helper.MustGetNetINodeWithReadAcquire(fsINode.NetINodeID,
		0, p.DefaultNetBlockCap, p.DefaultMemBlockCap)
	return err
}

func (p *FsINodeDriver) PrepareFsINodeForCreate(fsINode *types.FsINode,
	fsINodeID *types.FsINodeID, netINodeID *types.NetINodeID, parentID types.FsINodeID,
	name string, fsINodeType int, mode uint32,
	uid uint32, gid uint32, rdev uint32,
) error {
	var err error
	now := p.Timer.Now()
	nowt := types.DirTreeTime(now.Unix())
	nowtnsec := types.DirTreeTimeNsec(now.UnixNano())
	if fsINodeID != nil {
		fsINode.Ino = *fsINodeID
	} else {
		fsINode.Ino = p.helper.AllocFsINodeID()
	}

	if netINodeID == nil {
		if fsINodeType != types.FSINODE_TYPE_FILE {
			fsINode.NetINodeID = types.ZeroNetINodeID
		} else {
			err = p.AllocNetINodeID(fsINode)
			if err != nil {
				return err
			}
		}
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
	fsINode.Uid = uid
	fsINode.Gid = gid
	fsINode.Rdev = rdev
	return nil
}

func (p *FsINodeDriver) CreateFsINode(fsINode *types.FsINode) error {
	var err error
	err = p.helper.InsertFsINodeInDB(*fsINode)
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) MakeFsINodeKey(parentID types.FsINodeID, fsINodeName string) string {
	return strconv.FormatUint(parentID, 10) + fsINodeName
}
