package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/log"
	"soloos/common/solodbtypes"
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/common/timer"
	"soloos/solodb/offheap"
	"strconv"
	"sync"
	"time"
)

type FsINodeDriverHelper struct {
	solofstypes.AllocFsINodeID
	solofstypes.GetNetINode
	solofstypes.MustGetNetINode
	solofstypes.ReleaseNetINode
	solofstypes.DeleteFsINodeByIDInDB
	solofstypes.ListFsINodeByParentIDFromDB
	solofstypes.UpdateFsINodeInDB
	solofstypes.InsertFsINodeInDB
	solofstypes.FetchFsINodeByIDFromDB
	solofstypes.FetchFsINodeByNameFromDB
}

type FsINodeDriver struct {
	*soloosbase.SoloosEnv
	posixFs *PosixFs
	helper  FsINodeDriverHelper

	Timer timer.Timer

	fsINodesByIDTable offheap.LKVTableWithUint64
	fsINodesByPath    sync.Map

	SysFsINode  [2]solofstypes.FsINodeUintptr
	RootFsINode solofstypes.FsINodeUintptr

	EntryTtl           time.Duration
	EntryAttrValid     uint64
	EntryAttrValidNsec uint32

	DefaultNetBlockCap int
	DefaultMemBlockCap int
}

func (p *FsINodeDriver) Init(
	soloosEnv *soloosbase.SoloosEnv,
	posixFs *PosixFs,
	defaultNetBlockCap int,
	defaultMemBlockCap int,
	allocFsINodeID solofstypes.AllocFsINodeID,
	getNetINode solofstypes.GetNetINode,
	mustGetNetINode solofstypes.MustGetNetINode,
	releaseNetINode solofstypes.ReleaseNetINode,
	deleteFsINodeByIDInDB solofstypes.DeleteFsINodeByIDInDB,
	listFsINodeByParentIDFromDB solofstypes.ListFsINodeByParentIDFromDB,
	updateFsINodeInDB solofstypes.UpdateFsINodeInDB,
	insertFsINodeInDB solofstypes.InsertFsINodeInDB,
	fetchFsINodeByIDFromDB solofstypes.FetchFsINodeByIDFromDB,
	fetchFsINodeByNameFromDB solofstypes.FetchFsINodeByNameFromDB,
) error {
	var err error

	p.SoloosEnv = soloosEnv
	p.posixFs = posixFs

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
		int(solofstypes.FsINodeStructSize), -1, offheap.DefaultKVTableSharedCount,
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

	return nil
}

func (p *FsINodeDriver) fsINodesByIDTableInvokeBeforeReleaseObjectFunc(uObject uintptr) {
	var uFsINode = solofstypes.FsINodeUintptr(uObject)
	var pFsINode = uFsINode.Ptr()
	p.helper.ReleaseNetINode(pFsINode.UNetINode)
	p.CleanFsINodeAssitCache(pFsINode.Meta.ParentID, pFsINode.Meta.Name())
	uFsINode.Ptr().Reset()
}

func (p *FsINodeDriver) fsINodesByIDTablePrepareNewObjectFunc(uFsINode solofstypes.FsINodeUintptr,
	afterSetNewObj offheap.KVTableAfterSetNewObj) bool {
	var isNewObjectSetted bool
	if afterSetNewObj != nil {
		uFsINode.Ptr().Meta.Ino = uFsINode.Ptr().LKVTableObjectWithUint64.ID
		uFsINode.Ptr().Meta.NetINodeID = solofstypes.ZeroNetINodeID
		afterSetNewObj()
		isNewObjectSetted = true
	} else {
		isNewObjectSetted = false
	}
	return isNewObjectSetted
}

func (p *FsINodeDriver) SetHelper(
	allocFsINodeID solofstypes.AllocFsINodeID,
	getNetINode solofstypes.GetNetINode,
	mustGetNetINode solofstypes.MustGetNetINode,
	releaseNetINode solofstypes.ReleaseNetINode,
	deleteFsINodeByIDInDB solofstypes.DeleteFsINodeByIDInDB,
	listFsINodeByParentIDFromDB solofstypes.ListFsINodeByParentIDFromDB,
	updateFsINodeInDB solofstypes.UpdateFsINodeInDB,
	insertFsINodeInDB solofstypes.InsertFsINodeInDB,
	fetchFsINodeByIDFromDB solofstypes.FetchFsINodeByIDFromDB,
	fetchFsINodeByNameFromDB solofstypes.FetchFsINodeByNameFromDB,
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
		uFsINode    solofstypes.FsINodeUintptr
		fsINodeMeta solofstypes.FsINodeMeta
		ino         solofstypes.FsINodeID
		code        fsapi.Status
		err         error
	)

	ino = solofstypes.RootFsINodeID
	code = p.posixFs.SimpleMkdir(&fsINodeMeta, &ino, solofstypes.RootFsINodeParentID, 0777, "", 0, 0, solofstypes.FS_RDEV)
	if code != fsapi.OK {
		log.Warn("mkdir root error ", code)
	}

	ino, err = p.helper.AllocFsINodeID(p.posixFs.NameSpaceID)
	if err != nil {
		return err
	}
	code = p.posixFs.SimpleMkdir(&fsINodeMeta, &ino, solofstypes.RootFsINodeID, 0777, "tmp", 0, 0, solofstypes.FS_RDEV)
	if code != fsapi.OK {
		log.Warn("mkdir tmp error", code)
	}
	uFsINode, err = p.GetFsINodeByName(solofstypes.RootFsINodeParentID, "")
	// no need release: defer p.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}
	p.RootFsINode = uFsINode
	p.RootFsINode.Ptr().Acquire()
	if err != nil {
		return err
	}

	for i := 0; i < len(p.SysFsINode); i++ {
		var (
			uNewObject     offheap.LKVTableObjectUPtrWithUint64
			newInoKey      = solofstypes.FsINodeID(i)
			afterSetNewObj offheap.KVTableAfterSetNewObj
		)
		if newInoKey == p.RootFsINode.Ptr().Meta.Ino {
			continue
		}

		uNewObject, afterSetNewObj = p.fsINodesByIDTable.MustGetObject(newInoKey)
		uFsINode = solofstypes.FsINodeUintptr(uNewObject)
		if afterSetNewObj != nil {
			afterSetNewObj()
		}
		uFsINode.Ptr().Acquire()
		uFsINode.Ptr().Meta.NetINodeID = solofstypes.ZeroNetINodeID
		p.SysFsINode[i] = uFsINode
		// no need release: defer p.ReleaseFsINode(uFsINode)
	}

	return nil
}

func (p *FsINodeDriver) checkIfNeedNetINode(fsINodeType int) bool {
	return fsINodeType == solofstypes.FSINODE_TYPE_FILE
}

// ensureFsINodeValidInCache return false if fsinode invalid in cache
// if fsinode invalid, delete cache
func (p *FsINodeDriver) ensureFsINodeValidInCache(uFsINode solofstypes.FsINodeUintptr) bool {
	if p.Timer.Now().Unix()-uFsINode.Ptr().Meta.LoadInMemAt < int64(p.EntryAttrValid) {
		return true
	}

	return false
}

func (p *FsINodeDriver) updateFsINodeInCache(pFsINodeMeta *solofstypes.FsINodeMeta) error {
	var (
		uFsINode solofstypes.FsINodeUintptr
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

func (p *FsINodeDriver) commitFsINodeInCache(uFsINode solofstypes.FsINodeUintptr) error {
	var err error
	var pFsINode = uFsINode.Ptr()
	pFsINode.Meta.LoadInMemAt = p.Timer.Now().Unix()
	p.fsINodesByPath.Store(p.MakeFsINodeKey(pFsINode.Meta.ParentID, pFsINode.Meta.Name()), uFsINode)

	// ensure NetINode
	if pFsINode.UNetINode == 0 && pFsINode.Meta.NetINodeID != solofstypes.ZeroNetINodeID {
		pFsINode.UNetINode, err = p.helper.MustGetNetINode(pFsINode.Meta.NetINodeID,
			0, p.DefaultNetBlockCap, p.DefaultMemBlockCap)
		pFsINode.UNetINode.Ptr().LastCommitSize = pFsINode.UNetINode.Ptr().Size
	}
	pFsINode.IsDBMetaDataInited.Store(solodbtypes.MetaDataStateInited)

	return err
}

func (p *FsINodeDriver) DeleteFsINodeCache(uFsINode solofstypes.FsINodeUintptr,
	parentID solofstypes.FsINodeID, name string) {
	p.fsINodesByIDTable.ForceDeleteAfterReleaseDone(offheap.LKVTableObjectUPtrWithUint64(uFsINode))
}

func (p *FsINodeDriver) CleanFsINodeAssitCache(parentID solofstypes.FsINodeID, fsINodeName string) {
	p.fsINodesByPath.Delete(p.MakeFsINodeKey(parentID, fsINodeName))
}

func (p *FsINodeDriver) GetFsINodeByIDThroughHardLink(fsINodeID solofstypes.FsINodeID) (solofstypes.FsINodeUintptr, error) {
	var (
		uFsINode solofstypes.FsINodeUintptr
		err      error
	)
	for {
		uFsINode, err = p.GetFsINodeByID(fsINodeID)
		if err == nil {
			if uFsINode.Ptr().Meta.Type != solofstypes.FSINODE_TYPE_HARD_LINK {
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

func (p *FsINodeDriver) GetFsINodeByID(fsINodeID solofstypes.FsINodeID) (solofstypes.FsINodeUintptr, error) {
	if fsINodeID < solofstypes.RootFsINodeID {
		return p.SysFsINode[fsINodeID], nil
	}

	var (
		uFsINode       solofstypes.FsINodeUintptr
		pFsINode       *solofstypes.FsINode
		uObject        offheap.LKVTableObjectUPtrWithUint64
		afterSetNewObj offheap.KVTableAfterSetNewObj
		err            error
	)

	uObject, afterSetNewObj = p.fsINodesByIDTable.MustGetObject(fsINodeID)
	uFsINode = solofstypes.FsINodeUintptr(uObject)
	pFsINode = uFsINode.Ptr()
	p.fsINodesByIDTablePrepareNewObjectFunc(uFsINode, afterSetNewObj)
	if pFsINode.IsDBMetaDataInited.Load() == solodbtypes.MetaDataStateInited &&
		p.ensureFsINodeValidInCache(uFsINode) == true {
		return uFsINode, nil
	}

	pFsINode.IsDBMetaDataInited.LockContext()
	if pFsINode.IsDBMetaDataInited.Load() == solodbtypes.MetaDataStateUninited ||
		p.ensureFsINodeValidInCache(uFsINode) == false {
		pFsINode.Meta, err = p.helper.FetchFsINodeByIDFromDB(p.posixFs.NameSpaceID, fsINodeID)
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

func (p *FsINodeDriver) GetFsINodeByName(parentID solofstypes.FsINodeID,
	fsINodeName string) (solofstypes.FsINodeUintptr, error) {
	var (
		iObject  interface{}
		uFsINode solofstypes.FsINodeUintptr
		exists   bool
		err      error
	)

	iObject, exists = p.fsINodesByPath.Load(p.MakeFsINodeKey(parentID, fsINodeName))
	if iObject != nil {
		uFsINode = iObject.(solofstypes.FsINodeUintptr)
	}
	if exists && p.ensureFsINodeValidInCache(uFsINode) == true {
		uFsINode.Ptr().Acquire()
		return uFsINode, nil
	}

	// TODO only get fsinode id is ok
	var fsINodeMeta solofstypes.FsINodeMeta
	fsINodeMeta, err = p.helper.FetchFsINodeByNameFromDB(p.posixFs.NameSpaceID, parentID, fsINodeName)
	if err != nil {
		return 0, err
	}

	uFsINode, err = p.GetFsINodeByID(fsINodeMeta.Ino)
	if uFsINode != 0 {
		uFsINode.Ptr().Meta = fsINodeMeta
	}
	return uFsINode, err
}

func (p *FsINodeDriver) ReleaseFsINode(uFsINode solofstypes.FsINodeUintptr) {
	p.fsINodesByIDTable.ReleaseObject(offheap.LKVTableObjectUPtrWithUint64(uFsINode))
}

func (p *FsINodeDriver) UpdateFsINode(pFsINodeMeta *solofstypes.FsINodeMeta) error {
	var err error
	pFsINodeMeta.Ctime = solofstypes.DirTreeTime(p.Timer.Now().Unix())
	err = p.helper.UpdateFsINodeInDB(p.posixFs.NameSpaceID, *pFsINodeMeta)
	if err != nil {
		return err
	}

	err = p.updateFsINodeInCache(pFsINodeMeta)
	if err != nil {
		return err
	}

	return err
}

func (p *FsINodeDriver) RefreshFsINodeMetaACMtime(pFsINodeMeta *solofstypes.FsINodeMeta) error {
	var err error
	now := p.Timer.Now()
	nowUnixNano := now.UnixNano()
	if nowUnixNano-pFsINodeMeta.LastModifyACMTime < int64(time.Millisecond)*600 {
		return nil
	}

	nowt := solofstypes.DirTreeTime(now.Unix())
	nowtnsec := solofstypes.DirTreeTimeNsec(now.UnixNano())

	pFsINodeMeta.Atime = nowt
	pFsINodeMeta.Atimensec = nowtnsec
	pFsINodeMeta.Ctime = nowt
	pFsINodeMeta.Ctimensec = nowtnsec
	pFsINodeMeta.Mtime = nowt
	pFsINodeMeta.Mtimensec = nowtnsec

	err = p.helper.UpdateFsINodeInDB(p.posixFs.NameSpaceID, *pFsINodeMeta)
	if err != nil {
		return err
	}

	pFsINodeMeta.LastModifyACMTime = nowUnixNano
	return err
}

func (p *FsINodeDriver) RefreshFsINodeACMtime(uFsINode solofstypes.FsINodeUintptr) error {
	var (
		pFsINode = uFsINode.Ptr()
		err      error
	)

	now := p.Timer.Now()
	nowUnixNano := now.UnixNano()
	if nowUnixNano-pFsINode.Meta.LastModifyACMTime < int64(time.Millisecond)*600 {
		return nil
	}

	nowt := solofstypes.DirTreeTime(now.Unix())
	nowtnsec := solofstypes.DirTreeTimeNsec(now.UnixNano())

	pFsINode.Meta.Atime = nowt
	pFsINode.Meta.Atimensec = nowtnsec
	pFsINode.Meta.Ctime = nowt
	pFsINode.Meta.Ctimensec = nowtnsec
	pFsINode.Meta.Mtime = nowt
	pFsINode.Meta.Mtimensec = nowtnsec

	err = p.helper.UpdateFsINodeInDB(p.posixFs.NameSpaceID, pFsINode.Meta)
	if err != nil {
		return err
	}

	pFsINode.Meta.LastModifyACMTime = nowUnixNano
	return err
}

func (p *FsINodeDriver) RefreshFsINodeACMtimeByIno(fsINodeID solofstypes.FsINodeID) error {
	var (
		uFsINode solofstypes.FsINodeUintptr
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

func (p *FsINodeDriver) AllocNetINode(fsINodeMeta *solofstypes.FsINodeMeta) error {
	//TODO improve alloc NetInodeID
	solofstypes.InitTmpNetINodeID(&fsINodeMeta.NetINodeID)
	//TODO config memBlockSize netBlockSize
	var uNetINode, err = p.helper.MustGetNetINode(fsINodeMeta.NetINodeID,
		0, p.DefaultNetBlockCap, p.DefaultMemBlockCap)
	p.helper.ReleaseNetINode(uNetINode)
	return err
}

func (p *FsINodeDriver) PrepareFsINodeForCreate(fsINodeMeta *solofstypes.FsINodeMeta,
	fsINodeID *solofstypes.FsINodeID, netINodeID *solofstypes.NetINodeID, parentID solofstypes.FsINodeID,
	name string, fsINodeType int, mode uint32,
	uid uint32, gid uint32, rdev uint32,
) error {
	var err error
	now := p.Timer.Now()
	nowt := solofstypes.DirTreeTime(now.Unix())
	nowtnsec := solofstypes.DirTreeTimeNsec(now.UnixNano())
	if fsINodeID != nil {
		fsINodeMeta.Ino = *fsINodeID
	} else {
		fsINodeMeta.Ino, err = p.helper.AllocFsINodeID(p.posixFs.NameSpaceID)
		if err != nil {
			return err
		}
	}

	if netINodeID == nil {
		if fsINodeType != solofstypes.FSINODE_TYPE_FILE {
			fsINodeMeta.NetINodeID = solofstypes.ZeroNetINodeID
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

func (p *FsINodeDriver) CreateFsINode(fsINodeMeta *solofstypes.FsINodeMeta) error {
	var err error
	err = p.helper.InsertFsINodeInDB(p.posixFs.NameSpaceID, *fsINodeMeta)
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) MakeFsINodeKey(parentID solofstypes.FsINodeID, fsINodeName string) string {
	return strconv.FormatUint(parentID, 10) + "_" + fsINodeName
}
