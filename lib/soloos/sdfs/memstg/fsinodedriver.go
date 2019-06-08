package memstg

import (
	"soloos/common/fsapitypes"
	"soloos/common/log"
	"soloos/common/sdbapitypes"
	"soloos/common/sdfsapitypes"
	"soloos/common/soloosbase"
	"soloos/common/timer"
	"soloos/sdbone/offheap"
	"soloos/sdfs/types"
	"strconv"
	"sync"
	"time"
)

type FsINodeDriverHelper struct {
	sdfsapitypes.AllocFsINodeID
	sdfsapitypes.GetNetINode
	sdfsapitypes.MustGetNetINode
	sdfsapitypes.ReleaseNetINode
	sdfsapitypes.DeleteFsINodeByIDInDB
	sdfsapitypes.ListFsINodeByParentIDFromDB
	sdfsapitypes.UpdateFsINodeInDB
	sdfsapitypes.InsertFsINodeInDB
	sdfsapitypes.FetchFsINodeByIDFromDB
	sdfsapitypes.FetchFsINodeByNameFromDB
}

type FsINodeDriver struct {
	*soloosbase.SoloOSEnv
	dirTreeStg *PosixFS
	helper     FsINodeDriverHelper

	Timer timer.Timer

	fsINodesByIDTable offheap.LKVTableWithUint64
	fsINodesByPath    sync.Map

	SysFsINode  [2]sdfsapitypes.FsINodeUintptr
	RootFsINode sdfsapitypes.FsINodeUintptr

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
	dirTreeStg *PosixFS,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
	allocFsINodeID sdfsapitypes.AllocFsINodeID,
	getNetINode sdfsapitypes.GetNetINode,
	mustGetNetINode sdfsapitypes.MustGetNetINode,
	releaseNetINode sdfsapitypes.ReleaseNetINode,
	deleteFsINodeByIDInDB sdfsapitypes.DeleteFsINodeByIDInDB,
	listFsINodeByParentIDFromDB sdfsapitypes.ListFsINodeByParentIDFromDB,
	updateFsINodeInDB sdfsapitypes.UpdateFsINodeInDB,
	insertFsINodeInDB sdfsapitypes.InsertFsINodeInDB,
	fetchFsINodeByIDFromDB sdfsapitypes.FetchFsINodeByIDFromDB,
	fetchFsINodeByNameFromDB sdfsapitypes.FetchFsINodeByNameFromDB,
	// FIXAttrDriver
	deleteFIXAttrInDB sdfsapitypes.DeleteFIXAttrInDB,
	replaceFIXAttrInDB sdfsapitypes.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB sdfsapitypes.GetFIXAttrByInoFromDB,
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
		int(sdfsapitypes.FsINodeStructSize), -1, offheap.DefaultKVTableSharedCount,
		p.fsINodesByIDTableInvokeBeforeReleaseObjectFunc)
	if err != nil {
		return err
	}

	err = p.prepareBaseDir()
	if err != nil {
		return err
	}

	p.EntryTtl = 3 * time.Second
	SplitDuration(p.EntryTtl, &p.EntryAttrValid, &p.EntryAttrValidNsec)

	err = p.SoloOSEnv.OffheapDriver.InitHKVTableWithUint64(&p.INodeRWMutexTable, "INodeRWMutex",
		int(sdfsapitypes.INodeRWMutexStructSize), -1, offheap.DefaultKVTableSharedCount,
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
	var uFsINode = sdfsapitypes.FsINodeUintptr(uObject)
	var pFsINode = uFsINode.Ptr()
	p.helper.ReleaseNetINode(pFsINode.UNetINode)
	p.CleanFsINodeAssitCache(pFsINode.Meta.ParentID, pFsINode.Meta.Name())
	uFsINode.Ptr().Reset()
}

func (p *FsINodeDriver) fsINodesByIDTablePrepareNewObjectFunc(uFsINode sdfsapitypes.FsINodeUintptr,
	afterSetNewObj offheap.KVTableAfterSetNewObj) bool {
	var isNewObjectSetted bool
	if afterSetNewObj != nil {
		uFsINode.Ptr().Meta.Ino = uFsINode.Ptr().LKVTableObjectWithUint64.ID
		uFsINode.Ptr().Meta.NetINodeID = sdfsapitypes.ZeroNetINodeID
		afterSetNewObj()
		isNewObjectSetted = true
	} else {
		isNewObjectSetted = false
	}
	return isNewObjectSetted
}

func (p *FsINodeDriver) SetHelper(
	allocFsINodeID sdfsapitypes.AllocFsINodeID,
	getNetINode sdfsapitypes.GetNetINode,
	mustGetNetINode sdfsapitypes.MustGetNetINode,
	releaseNetINode sdfsapitypes.ReleaseNetINode,
	deleteFsINodeByIDInDB sdfsapitypes.DeleteFsINodeByIDInDB,
	listFsINodeByParentIDFromDB sdfsapitypes.ListFsINodeByParentIDFromDB,
	updateFsINodeInDB sdfsapitypes.UpdateFsINodeInDB,
	insertFsINodeInDB sdfsapitypes.InsertFsINodeInDB,
	fetchFsINodeByIDFromDB sdfsapitypes.FetchFsINodeByIDFromDB,
	fetchFsINodeByNameFromDB sdfsapitypes.FetchFsINodeByNameFromDB,
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
		uFsINode    sdfsapitypes.FsINodeUintptr
		fsINodeMeta sdfsapitypes.FsINodeMeta
		ino         sdfsapitypes.FsINodeID
		code        fsapitypes.Status
		err         error
	)

	ino = sdfsapitypes.RootFsINodeID
	code = p.dirTreeStg.SimpleMkdir(&fsINodeMeta, &ino, sdfsapitypes.RootFsINodeParentID, 0777, "", 0, 0, types.FS_RDEV)
	if code != fsapitypes.OK {
		log.Warn("mkdir root error ", code)
	}

	ino = p.helper.AllocFsINodeID()
	code = p.dirTreeStg.SimpleMkdir(&fsINodeMeta, &ino, sdfsapitypes.RootFsINodeID, 0777, "tmp", 0, 0, types.FS_RDEV)
	if code != fsapitypes.OK {
		log.Warn("mkdir tmp error", code)
	}
	uFsINode, err = p.GetFsINodeByName(sdfsapitypes.RootFsINodeParentID, "")
	// no need release: defer p.ReleaseFsINode(uFsINode)
	p.RootFsINode = uFsINode
	p.RootFsINode.Ptr().Acquire()
	if err != nil {
		return err
	}

	for i := 0; i < len(p.SysFsINode); i++ {
		var (
			uNewObject     offheap.LKVTableObjectUPtrWithUint64
			newInoKey      = sdfsapitypes.FsINodeID(i)
			afterSetNewObj offheap.KVTableAfterSetNewObj
		)
		if newInoKey == p.RootFsINode.Ptr().Meta.Ino {
			continue
		}

		uNewObject, afterSetNewObj = p.fsINodesByIDTable.MustGetObject(newInoKey)
		uFsINode = sdfsapitypes.FsINodeUintptr(uNewObject)
		if afterSetNewObj != nil {
			afterSetNewObj()
		}
		uFsINode.Ptr().Acquire()
		uFsINode.Ptr().Meta.NetINodeID = sdfsapitypes.ZeroNetINodeID
		p.SysFsINode[i] = uFsINode
		// no need release: defer p.ReleaseFsINode(uFsINode)
	}

	return nil
}

func (p *FsINodeDriver) checkIfNeedNetINode(fsINodeType int) bool {
	return fsINodeType == types.FSINODE_TYPE_FILE
}

// ensureFsINodeValidInCache return false if fsinode invalid in cache
// if fsinode invalid, delete cache
func (p *FsINodeDriver) ensureFsINodeValidInCache(uFsINode sdfsapitypes.FsINodeUintptr) bool {
	if p.Timer.Now().Unix()-uFsINode.Ptr().Meta.LoadInMemAt < int64(p.EntryAttrValid) {
		return true
	}

	return false
}

func (p *FsINodeDriver) updateFsINodeInCache(pFsINodeMeta *sdfsapitypes.FsINodeMeta) error {
	var (
		uFsINode sdfsapitypes.FsINodeUintptr
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

func (p *FsINodeDriver) commitFsINodeInCache(uFsINode sdfsapitypes.FsINodeUintptr) error {
	var err error
	var pFsINode = uFsINode.Ptr()
	pFsINode.Meta.LoadInMemAt = p.Timer.Now().Unix()
	p.fsINodesByPath.Store(p.MakeFsINodeKey(pFsINode.Meta.ParentID, pFsINode.Meta.Name()), uFsINode)

	// ensure NetINode
	if pFsINode.UNetINode == 0 && pFsINode.Meta.NetINodeID != sdfsapitypes.ZeroNetINodeID {
		pFsINode.UNetINode, err = p.helper.MustGetNetINode(pFsINode.Meta.NetINodeID,
			0, p.DefaultNetBlockCap, p.DefaultMemBlockCap)
		pFsINode.UNetINode.Ptr().LastCommitSize = pFsINode.UNetINode.Ptr().Size
	}
	pFsINode.IsDBMetaDataInited.Store(sdbapitypes.MetaDataStateInited)

	return err
}

func (p *FsINodeDriver) DeleteFsINodeCache(uFsINode sdfsapitypes.FsINodeUintptr,
	parentID sdfsapitypes.FsINodeID, name string) {
	p.fsINodesByIDTable.ForceDeleteAfterReleaseDone(offheap.LKVTableObjectUPtrWithUint64(uFsINode))
}

func (p *FsINodeDriver) CleanFsINodeAssitCache(parentID sdfsapitypes.FsINodeID, fsINodeName string) {
	p.fsINodesByPath.Delete(p.MakeFsINodeKey(parentID, fsINodeName))
}

func (p *FsINodeDriver) GetFsINodeByIDThroughHardLink(fsINodeID sdfsapitypes.FsINodeID) (sdfsapitypes.FsINodeUintptr, error) {
	var (
		uFsINode sdfsapitypes.FsINodeUintptr
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

func (p *FsINodeDriver) GetFsINodeByID(fsINodeID sdfsapitypes.FsINodeID) (sdfsapitypes.FsINodeUintptr, error) {
	if fsINodeID < sdfsapitypes.RootFsINodeID {
		return p.SysFsINode[fsINodeID], nil
	}

	var (
		uFsINode       sdfsapitypes.FsINodeUintptr
		pFsINode       *sdfsapitypes.FsINode
		uObject        offheap.LKVTableObjectUPtrWithUint64
		afterSetNewObj offheap.KVTableAfterSetNewObj
		err            error
	)

	uObject, afterSetNewObj = p.fsINodesByIDTable.MustGetObject(fsINodeID)
	uFsINode = sdfsapitypes.FsINodeUintptr(uObject)
	pFsINode = uFsINode.Ptr()
	p.fsINodesByIDTablePrepareNewObjectFunc(uFsINode, afterSetNewObj)
	if pFsINode.IsDBMetaDataInited.Load() == sdbapitypes.MetaDataStateInited &&
		p.ensureFsINodeValidInCache(uFsINode) == true {
		return uFsINode, nil
	}

	pFsINode.IsDBMetaDataInited.LockContext()
	if pFsINode.IsDBMetaDataInited.Load() == sdbapitypes.MetaDataStateUninited ||
		p.ensureFsINodeValidInCache(uFsINode) == false {
		pFsINode.Meta, err = p.helper.FetchFsINodeByIDFromDB(fsINodeID)
		if err != nil {
			p.ReleaseFsINode(uFsINode)
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

func (p *FsINodeDriver) GetFsINodeByName(parentID sdfsapitypes.FsINodeID,
	fsINodeName string) (sdfsapitypes.FsINodeUintptr, error) {
	var (
		iObject  interface{}
		uFsINode sdfsapitypes.FsINodeUintptr
		exists   bool
		err      error
	)

	iObject, exists = p.fsINodesByPath.Load(p.MakeFsINodeKey(parentID, fsINodeName))
	if iObject != nil {
		uFsINode = iObject.(sdfsapitypes.FsINodeUintptr)
	}
	if exists && p.ensureFsINodeValidInCache(uFsINode) == true {
		uFsINode.Ptr().Acquire()
		return uFsINode, nil
	}

	// TODO only get fsinode id is ok
	var fsINodeMeta sdfsapitypes.FsINodeMeta
	fsINodeMeta, err = p.helper.FetchFsINodeByNameFromDB(parentID, fsINodeName)
	if err != nil {
		return 0, err
	}

	uFsINode, err = p.GetFsINodeByID(fsINodeMeta.Ino)
	if uFsINode != 0 {
		uFsINode.Ptr().Meta = fsINodeMeta
	}
	return uFsINode, err
}

func (p *FsINodeDriver) ReleaseFsINode(uFsINode sdfsapitypes.FsINodeUintptr) {
	p.fsINodesByIDTable.ReleaseObject(offheap.LKVTableObjectUPtrWithUint64(uFsINode))
}

func (p *FsINodeDriver) UpdateFsINodeInDB(pFsINodeMeta *sdfsapitypes.FsINodeMeta) error {
	var err error
	pFsINodeMeta.Ctime = sdfsapitypes.DirTreeTime(p.Timer.Now().Unix())
	err = p.helper.UpdateFsINodeInDB(*pFsINodeMeta)
	if err != nil {
		return err
	}

	err = p.updateFsINodeInCache(pFsINodeMeta)
	if err != nil {
		return err
	}

	return err
}

func (p *FsINodeDriver) RefreshFsINodeMetaACMtime(pFsINodeMeta *sdfsapitypes.FsINodeMeta) error {
	var err error
	now := p.Timer.Now()
	nowUnixNano := now.UnixNano()
	if nowUnixNano-pFsINodeMeta.LastModifyACMTime < int64(time.Millisecond)*600 {
		return nil
	}

	nowt := sdfsapitypes.DirTreeTime(now.Unix())
	nowtnsec := sdfsapitypes.DirTreeTimeNsec(now.UnixNano())

	pFsINodeMeta.Atime = nowt
	pFsINodeMeta.Atimensec = nowtnsec
	pFsINodeMeta.Ctime = nowt
	pFsINodeMeta.Ctimensec = nowtnsec
	pFsINodeMeta.Mtime = nowt
	pFsINodeMeta.Mtimensec = nowtnsec

	err = p.helper.UpdateFsINodeInDB(*pFsINodeMeta)
	if err != nil {
		return err
	}

	pFsINodeMeta.LastModifyACMTime = nowUnixNano
	return err
}

func (p *FsINodeDriver) RefreshFsINodeACMtime(uFsINode sdfsapitypes.FsINodeUintptr) error {
	var (
		pFsINode = uFsINode.Ptr()
		err      error
	)

	now := p.Timer.Now()
	nowUnixNano := now.UnixNano()
	if nowUnixNano-pFsINode.Meta.LastModifyACMTime < int64(time.Millisecond)*600 {
		return nil
	}

	nowt := sdfsapitypes.DirTreeTime(now.Unix())
	nowtnsec := sdfsapitypes.DirTreeTimeNsec(now.UnixNano())

	pFsINode.Meta.Atime = nowt
	pFsINode.Meta.Atimensec = nowtnsec
	pFsINode.Meta.Ctime = nowt
	pFsINode.Meta.Ctimensec = nowtnsec
	pFsINode.Meta.Mtime = nowt
	pFsINode.Meta.Mtimensec = nowtnsec

	err = p.helper.UpdateFsINodeInDB(pFsINode.Meta)
	if err != nil {
		return err
	}

	pFsINode.Meta.LastModifyACMTime = nowUnixNano
	return err
}

func (p *FsINodeDriver) RefreshFsINodeACMtimeByIno(fsINodeID sdfsapitypes.FsINodeID) error {
	var (
		uFsINode sdfsapitypes.FsINodeUintptr
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

func (p *FsINodeDriver) AllocNetINode(fsINodeMeta *sdfsapitypes.FsINodeMeta) error {
	//TODO improve alloc NetInodeID
	sdfsapitypes.InitTmpNetINodeID(&fsINodeMeta.NetINodeID)
	//TODO config memBlockSize netBlockSize
	var uNetINode, err = p.helper.MustGetNetINode(fsINodeMeta.NetINodeID,
		0, p.DefaultNetBlockCap, p.DefaultMemBlockCap)
	p.helper.ReleaseNetINode(uNetINode)
	return err
}

func (p *FsINodeDriver) PrepareFsINodeForCreate(fsINodeMeta *sdfsapitypes.FsINodeMeta,
	fsINodeID *sdfsapitypes.FsINodeID, netINodeID *sdfsapitypes.NetINodeID, parentID sdfsapitypes.FsINodeID,
	name string, fsINodeType int, mode uint32,
	uid uint32, gid uint32, rdev uint32,
) error {
	var err error
	now := p.Timer.Now()
	nowt := sdfsapitypes.DirTreeTime(now.Unix())
	nowtnsec := sdfsapitypes.DirTreeTimeNsec(now.UnixNano())
	if fsINodeID != nil {
		fsINodeMeta.Ino = *fsINodeID
	} else {
		fsINodeMeta.Ino = p.helper.AllocFsINodeID()
	}

	if netINodeID == nil {
		if fsINodeType != types.FSINODE_TYPE_FILE {
			fsINodeMeta.NetINodeID = sdfsapitypes.ZeroNetINodeID
		} else {
			err = p.AllocNetINode(fsINodeMeta)
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

func (p *FsINodeDriver) CreateFsINode(fsINodeMeta *sdfsapitypes.FsINodeMeta) error {
	var err error
	err = p.helper.InsertFsINodeInDB(*fsINodeMeta)
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) MakeFsINodeKey(parentID sdfsapitypes.FsINodeID, fsINodeName string) string {
	return strconv.FormatUint(parentID, 10) + "_" + fsINodeName
}
