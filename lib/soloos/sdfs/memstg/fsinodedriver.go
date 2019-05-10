package memstg

import (
	fsapitypes "soloos/common/fsapi/types"
	"soloos/common/log"
	sdbapitypes "soloos/common/sdbapi/types"
	sdfsapitypes "soloos/common/sdfsapi/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/common/timer"
	"soloos/sdbone/offheap"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"strconv"
	"sync"
	"time"
)

type FsINodeDriverHelper struct {
	api.AllocFsINodeID
	api.GetNetINode
	api.MustGetNetINode
	api.ReleaseNetINode
	api.DeleteFsINodeByIDInDB
	api.ListFsINodeByParentIDFromDB
	api.UpdateFsINodeInDB
	api.InsertFsINodeInDB
	api.FetchFsINodeByIDFromDB
	api.FetchFsINodeByNameFromDB
}

type FsINodeDriver struct {
	*soloosbase.SoloOSEnv
	dirTreeStg *DirTreeStg
	helper     FsINodeDriverHelper

	Timer timer.Timer

	fsINodesRWMutex   sync.RWMutex
	fsINodesByIDTable offheap.LKVTableWithUint64
	fsINodesByPath    map[string]types.FsINodeUintptr

	SysFsINode  [2]types.FsINodeUintptr
	RootFsINode types.FsINodeUintptr

	EntryTtl           time.Duration
	EntryAttrValid     uint64
	EntryAttrValidNsec uint32

	INodeRWMutexTable offheap.HKVTableWithUint64

	FIXAttrDriver FIXAttrDriver

	DefaultNetBlockCap int
	DefaultMemBlockCap int
}

func (p *FsINodeDriver) Init(
	soloOSEnv *soloosbase.SoloOSEnv,
	dirTreeStg *DirTreeStg,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
	allocFsINodeID api.AllocFsINodeID,
	getNetINode api.GetNetINode,
	mustGetNetINode api.MustGetNetINode,
	releaseNetINode api.ReleaseNetINode,
	deleteFsINodeByIDInDB api.DeleteFsINodeByIDInDB,
	listFsINodeByParentIDFromDB api.ListFsINodeByParentIDFromDB,
	updateFsINodeInDB api.UpdateFsINodeInDB,
	insertFsINodeInDB api.InsertFsINodeInDB,
	fetchFsINodeByIDFromDB api.FetchFsINodeByIDFromDB,
	fetchFsINodeByNameFromDB api.FetchFsINodeByNameFromDB,
	// FIXAttrDriver
	deleteFIXAttrInDB api.DeleteFIXAttrInDB,
	replaceFIXAttrInDB api.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB api.GetFIXAttrByInoFromDB,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
	p.dirTreeStg = dirTreeStg

	err = p.Timer.Init()
	if err != nil {
		return err
	}

	p.SetHelper(
		allocFsINodeID,
		getNetINode,
		mustGetNetINode,
		releaseNetINode,
		deleteFsINodeByIDInDB,
		listFsINodeByParentIDFromDB,
		updateFsINodeInDB,
		insertFsINodeInDB,
		fetchFsINodeByIDFromDB,
		fetchFsINodeByNameFromDB,
	)

	p.DefaultNetBlockCap = defaultNetBlockCap
	p.DefaultMemBlockCap = defaultMemBlockCap

	err = p.fsINodesByIDTable.Init("FsINode",
		int(types.FsINodeStructSize), -1, offheap.DefaultKVTableSharedCount,
		p.fsINodesByIDTableInvokeBeforeReleaseObjectFunc)
	if err != nil {
		return err
	}

	p.fsINodesByPath = make(map[string]types.FsINodeUintptr)

	err = p.prepareBaseDir()
	if err != nil {
		return err
	}

	p.EntryTtl = 3 * time.Second
	SplitDuration(p.EntryTtl, &p.EntryAttrValid, &p.EntryAttrValidNsec)

	err = p.SoloOSEnv.OffheapDriver.InitHKVTableWithUint64(&p.INodeRWMutexTable, "INodeRWMutex",
		int(types.INodeRWMutexStructSize), -1, offheap.DefaultKVTableSharedCount,
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

func (p *FsINodeDriver) fsINodesByIDTableInvokeBeforeReleaseObjectFunc(uObject uintptr) {
	var uFsINode = types.FsINodeUintptr(uObject)
	p.helper.ReleaseNetINode(uFsINode.Ptr().UNetINode)
}

func (p *FsINodeDriver) SetHelper(
	allocFsINodeID api.AllocFsINodeID,
	getNetINode api.GetNetINode,
	mustGetNetINode api.MustGetNetINode,
	releaseNetINode api.ReleaseNetINode,
	deleteFsINodeByIDInDB api.DeleteFsINodeByIDInDB,
	listFsINodeByParentIDFromDB api.ListFsINodeByParentIDFromDB,
	updateFsINodeInDB api.UpdateFsINodeInDB,
	insertFsINodeInDB api.InsertFsINodeInDB,
	fetchFsINodeByIDFromDB api.FetchFsINodeByIDFromDB,
	fetchFsINodeByNameFromDB api.FetchFsINodeByNameFromDB,
) {
	p.helper = FsINodeDriverHelper{
		AllocFsINodeID:              allocFsINodeID,
		GetNetINode:                 getNetINode,
		MustGetNetINode:             mustGetNetINode,
		ReleaseNetINode:             releaseNetINode,
		DeleteFsINodeByIDInDB:       deleteFsINodeByIDInDB,
		ListFsINodeByParentIDFromDB: listFsINodeByParentIDFromDB,
		UpdateFsINodeInDB:           updateFsINodeInDB,
		InsertFsINodeInDB:           insertFsINodeInDB,
		FetchFsINodeByIDFromDB:      fetchFsINodeByIDFromDB,
		FetchFsINodeByNameFromDB:    fetchFsINodeByNameFromDB,
	}
}

func (p *FsINodeDriver) prepareBaseDir() error {
	var (
		uFsINode    types.FsINodeUintptr
		fsINodeMeta types.FsINodeMeta
		ino         types.FsINodeID
		code        fsapitypes.Status
		err         error
	)

	ino = types.RootFsINodeID
	code = p.dirTreeStg.SimpleMkdir(&fsINodeMeta, &ino, types.RootFsINodeParentID, 0777, "", 0, 0, types.FS_RDEV)
	if code != fsapitypes.OK {
		log.Warn("mkdir root error ", code)
	}

	ino = p.helper.AllocFsINodeID()
	code = p.dirTreeStg.SimpleMkdir(&fsINodeMeta, &ino, types.RootFsINodeID, 0777, "tmp", 0, 0, types.FS_RDEV)
	if code != fsapitypes.OK {
		log.Warn("mkdir tmp error", code)
	}
	uFsINode, err = p.GetFsINodeByName(types.RootFsINodeParentID, "")
	// no need release: defer p.ReleaseFsINode(uFsINode)
	p.RootFsINode = uFsINode
	if err != nil {
		return err
	}

	for i := 0; i < len(p.SysFsINode); i++ {
		var (
			uNewObject     offheap.LKVTableObjectUPtrWithUint64
			newInoKey      = types.FsINodeID(i)
			afterSetNewObj offheap.KVTableAfterSetNewObj
		)
		if newInoKey == p.RootFsINode.Ptr().Meta.Ino {
			continue
		}

		uNewObject, afterSetNewObj = p.fsINodesByIDTable.MustGetObject(newInoKey)
		if afterSetNewObj != nil {
			afterSetNewObj()
		}
		p.SysFsINode[i] = types.FsINodeUintptr(uNewObject)
		// no need release: defer p.ReleaseFsINode(uFsINode)
	}

	return nil
}

func (p *FsINodeDriver) checkIfNeedNetINode(fsINodeType int) bool {
	return fsINodeType == types.FSINODE_TYPE_FILE
}

// ensureFsINodeValidInCache return false if fsinode invalid in cache
// if fsinode invalid, delete cache
func (p *FsINodeDriver) ensureFsINodeValidInCache(uFsINode types.FsINodeUintptr) bool {
	if p.Timer.Now().Unix()-uFsINode.Ptr().Meta.LoadInMemAt < int64(p.EntryAttrValid) {
		return true
	}

	return false
}

func (p *FsINodeDriver) updateFsINodeInCache(pFsINodeMeta *types.FsINodeMeta) error {
	var (
		uFsINode types.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.GetFsINodeByID(pFsINodeMeta.Ino)
	defer p.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}

	pFsINodeMeta.LoadInMemAt = p.Timer.Now().Unix()
	uFsINode.Ptr().Meta = *pFsINodeMeta

	return nil
}

func (p *FsINodeDriver) commitFsINodeInCache(uFsINode types.FsINodeUintptr) error {
	var err error
	var pFsINode = uFsINode.Ptr()
	pFsINode.Meta.LoadInMemAt = p.Timer.Now().Unix()
	p.fsINodesByPath[p.MakeFsINodeKey(pFsINode.Meta.ParentID, pFsINode.Meta.Name())] = uFsINode

	// ensure NetINode
	if pFsINode.UNetINode == 0 {
		pFsINode.UNetINode, err = p.helper.MustGetNetINode(pFsINode.Meta.NetINodeID,
			0, p.DefaultNetBlockCap, p.DefaultMemBlockCap)
		pFsINode.UNetINode.Ptr().LastCommitSize = pFsINode.UNetINode.Ptr().Size
	}
	pFsINode.IsDBMetaDataInited.Store(sdbapitypes.MetaDataStateInited)

	return err
}

func (p *FsINodeDriver) fsINodesByIDTablePrepareNewObjectFunc(uFsINode types.FsINodeUintptr,
	afterSetNewObj offheap.KVTableAfterSetNewObj) bool {
	var isNewObjectSetted bool
	if afterSetNewObj != nil {
		uFsINode.Ptr().Meta.Ino = uFsINode.Ptr().LKVTableObjectWithUint64.ID
		afterSetNewObj()
		isNewObjectSetted = true
	} else {
		isNewObjectSetted = false
	}
	return isNewObjectSetted
}

func (p *FsINodeDriver) DeleteFsINodeCache(uFsINode types.FsINodeUintptr,
	parentID types.FsINodeID, name string) {
	p.fsINodesRWMutex.Lock()
	p.fsINodesByIDTable.ForceDeleteAfterReleaseDone(offheap.LKVTableObjectUPtrWithUint64(uFsINode))
	delete(p.fsINodesByPath, p.MakeFsINodeKey(parentID, name))
	p.fsINodesRWMutex.Unlock()
}

func (p *FsINodeDriver) CleanFsINodeAssitCache(parentID types.FsINodeID, fsINodeName string) {
	p.fsINodesRWMutex.Lock()
	delete(p.fsINodesByPath, p.MakeFsINodeKey(parentID, fsINodeName))
	p.fsINodesRWMutex.Unlock()
}

func (p *FsINodeDriver) GetFsINodeByIDThroughHardLink(fsINodeID types.FsINodeID) (types.FsINodeUintptr, error) {
	var (
		uFsINode types.FsINodeUintptr
		err      error
	)
	for {
		uFsINode, err = p.GetFsINodeByID(fsINodeID)
		if err == nil {
			if uFsINode.Ptr().Meta.Type != types.FSINODE_TYPE_HARD_LINK {
				return uFsINode, nil
			}
		}

		p.ReleaseFsINode(uFsINode)
		if err != nil {
			return 0, err
		}
		fsINodeID = uFsINode.Ptr().Meta.HardLinkIno
	}
}

func (p *FsINodeDriver) GetFsINodeByID(fsINodeID types.FsINodeID) (types.FsINodeUintptr, error) {
	if fsINodeID < types.RootFsINodeID {
		return p.SysFsINode[fsINodeID], nil
	}

	var (
		uFsINode       types.FsINodeUintptr
		pFsINode       *types.FsINode
		uObject        offheap.LKVTableObjectUPtrWithUint64
		afterSetNewObj offheap.KVTableAfterSetNewObj
		err            error
	)

	uObject, afterSetNewObj = p.fsINodesByIDTable.MustGetObject(fsINodeID)
	uFsINode = types.FsINodeUintptr(uObject)
	pFsINode = uFsINode.Ptr()
	p.fsINodesByIDTablePrepareNewObjectFunc(uFsINode, afterSetNewObj)
	if pFsINode.IsDBMetaDataInited.Load() == sdbapitypes.MetaDataStateInited &&
		p.ensureFsINodeValidInCache(uFsINode) == true {
		return uFsINode, nil
	}

	pFsINode.IsDBMetaDataInited.LockContext()
	if pFsINode.IsDBMetaDataInited.Load() == sdbapitypes.MetaDataStateUninited {
		err = p.helper.FetchFsINodeByIDFromDB(&pFsINode.Meta)
		if err != nil {
			if err == types.ErrObjectNotExists {
				defer p.DeleteFsINodeCache(uFsINode, pFsINode.Meta.ParentID, pFsINode.Meta.Name())
			} else {
				defer p.ReleaseFsINode(uFsINode)
			}
		} else {
			err = p.commitFsINodeInCache(uFsINode)
		}
	}
	pFsINode.IsDBMetaDataInited.UnlockContext()

	if err != nil {
		return 0, err
	}

	return uFsINode, nil
}

func (p *FsINodeDriver) GetFsINodeByName(parentID types.FsINodeID, fsINodeName string) (types.FsINodeUintptr, error) {
	var (
		uFsINode types.FsINodeUintptr
		exists   bool
		err      error
	)

	p.fsINodesRWMutex.RLock()
	uFsINode, exists = p.fsINodesByPath[p.MakeFsINodeKey(parentID, fsINodeName)]
	p.fsINodesRWMutex.RUnlock()
	if exists && p.ensureFsINodeValidInCache(uFsINode) == true {
		return uFsINode, nil
	}

	var fsINodeMeta types.FsINodeMeta
	fsINodeMeta.ParentID = parentID
	fsINodeMeta.SetName(fsINodeName)
	err = p.helper.FetchFsINodeByNameFromDB(&fsINodeMeta)
	if err != nil {
		return 0, err
	}

	uFsINode, err = p.GetFsINodeByID(fsINodeMeta.Ino)
	if uFsINode != 0 {
		uFsINode.Ptr().Meta = fsINodeMeta
	}
	return uFsINode, err
}

func (p *FsINodeDriver) ReleaseFsINode(uFsINode types.FsINodeUintptr) {
	p.fsINodesByIDTable.ReleaseObject(offheap.LKVTableObjectUPtrWithUint64(uFsINode))
}

func (p *FsINodeDriver) UpdateFsINodeInDB(pFsINodeMeta *types.FsINodeMeta) error {
	var err error
	pFsINodeMeta.Ctime = types.DirTreeTime(p.Timer.Now().Unix())
	err = p.helper.UpdateFsINodeInDB(pFsINodeMeta)
	if err != nil {
		return err
	}

	err = p.updateFsINodeInCache(pFsINodeMeta)
	if err != nil {
		return err
	}

	return err
}

func (p *FsINodeDriver) RefreshFsINodeMetaACMtime(fsINodeMeta *types.FsINodeMeta) error {
	var err error
	now := p.Timer.Now()
	nowUnixNano := now.UnixNano()
	if nowUnixNano-fsINodeMeta.LastModifyACMTime < int64(time.Millisecond)*600 {
		return nil
	}

	nowt := types.DirTreeTime(now.Unix())
	nowtnsec := types.DirTreeTimeNsec(now.UnixNano())

	fsINodeMeta.Atime = nowt
	fsINodeMeta.Atimensec = nowtnsec
	fsINodeMeta.Ctime = nowt
	fsINodeMeta.Ctimensec = nowtnsec
	fsINodeMeta.Mtime = nowt
	fsINodeMeta.Mtimensec = nowtnsec

	err = p.helper.UpdateFsINodeInDB(fsINodeMeta)
	if err != nil {
		return err
	}

	fsINodeMeta.LastModifyACMTime = nowUnixNano
	return err
}

func (p *FsINodeDriver) RefreshFsINodeACMtime(uFsINode types.FsINodeUintptr) error {
	var (
		pFsINode = uFsINode.Ptr()
		err      error
	)

	now := p.Timer.Now()
	nowUnixNano := now.UnixNano()
	if nowUnixNano-pFsINode.Meta.LastModifyACMTime < int64(time.Millisecond)*600 {
		return nil
	}

	nowt := types.DirTreeTime(now.Unix())
	nowtnsec := types.DirTreeTimeNsec(now.UnixNano())

	pFsINode.Meta.Atime = nowt
	pFsINode.Meta.Atimensec = nowtnsec
	pFsINode.Meta.Ctime = nowt
	pFsINode.Meta.Ctimensec = nowtnsec
	pFsINode.Meta.Mtime = nowt
	pFsINode.Meta.Mtimensec = nowtnsec

	err = p.helper.UpdateFsINodeInDB(&pFsINode.Meta)
	if err != nil {
		return err
	}

	pFsINode.Meta.LastModifyACMTime = nowUnixNano
	return err
}

func (p *FsINodeDriver) RefreshFsINodeACMtimeByIno(fsINodeID types.FsINodeID) error {
	var (
		uFsINode types.FsINodeUintptr
		err      error
	)

	uFsINode, err = p.GetFsINodeByID(fsINodeID)
	defer p.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}

	p.RefreshFsINodeACMtime(uFsINode)
	return err
}

func (p *FsINodeDriver) AllocNetINodeID(fsINodeMeta *types.FsINodeMeta) error {
	//TODO improve alloc NetInodeID
	sdfsapitypes.InitTmpNetINodeID(&fsINodeMeta.NetINodeID)
	//TODO config memBlockSize netBlockSize
	var uNetINode, err = p.helper.MustGetNetINode(fsINodeMeta.NetINodeID,
		0, p.DefaultNetBlockCap, p.DefaultMemBlockCap)
	p.helper.ReleaseNetINode(uNetINode)
	return err
}

func (p *FsINodeDriver) PrepareFsINodeForCreate(fsINodeMeta *types.FsINodeMeta,
	fsINodeID *types.FsINodeID, netINodeID *types.NetINodeID, parentID types.FsINodeID,
	name string, fsINodeType int, mode uint32,
	uid uint32, gid uint32, rdev uint32,
) error {
	var err error
	now := p.Timer.Now()
	nowt := types.DirTreeTime(now.Unix())
	nowtnsec := types.DirTreeTimeNsec(now.UnixNano())
	if fsINodeID != nil {
		fsINodeMeta.Ino = *fsINodeID
	} else {
		fsINodeMeta.Ino = p.helper.AllocFsINodeID()
	}

	if netINodeID == nil {
		if fsINodeType != types.FSINODE_TYPE_FILE {
			fsINodeMeta.NetINodeID = sdfsapitypes.ZeroNetINodeID
		} else {
			err = p.AllocNetINodeID(fsINodeMeta)
			if err != nil {
				return err
			}
		}
	} else {
		fsINodeMeta.NetINodeID = *netINodeID
	}

	fsINodeMeta.ParentID = parentID
	fsINodeMeta.SetName(name)
	fsINodeMeta.Type = fsINodeType
	fsINodeMeta.Atime = nowt
	fsINodeMeta.Ctime = nowt
	fsINodeMeta.Mtime = nowt
	fsINodeMeta.Atimensec = nowtnsec
	fsINodeMeta.Ctimensec = nowtnsec
	fsINodeMeta.Mtimensec = nowtnsec
	fsINodeMeta.Mode = mode
	fsINodeMeta.Nlink = 1
	fsINodeMeta.Uid = uid
	fsINodeMeta.Gid = gid
	fsINodeMeta.Rdev = rdev

	return nil
}

func (p *FsINodeDriver) CreateFsINode(fsINodeMeta *types.FsINodeMeta) error {
	var err error
	err = p.helper.InsertFsINodeInDB(fsINodeMeta)
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) MakeFsINodeKey(parentID types.FsINodeID, fsINodeName string) string {
	return strconv.FormatUint(parentID, 10) + "_" + fsINodeName
}
