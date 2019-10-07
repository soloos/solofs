package memstg

import (
	"soloos/common/solodbtypes"
	"soloos/common/solofsapi"
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/solodb/offheap"
)

type NetINodeDriverHelper struct {
	solofstypes.PrepareNetINodeMetaDataOnlyLoadDB
	solofstypes.PrepareNetINodeMetaDataWithStorDB
	solofstypes.NetINodeCommitSizeInDB
}

type NetINodeDriver struct {
	*soloosbase.SoloosEnv
	helper NetINodeDriverHelper

	solonnClient *solofsapi.SolonnClient

	netINodeTable offheap.LKVTableWithBytes64

	memBlockDriver *MemBlockDriver
	netBlockDriver *NetBlockDriver
}

func (p *NetINodeDriver) Init(soloosEnv *soloosbase.SoloosEnv,
	netBlockDriver *NetBlockDriver,
	memBlockDriver *MemBlockDriver,
	// for NetINodeDriverHelper
	solonnClient *solofsapi.SolonnClient,
	prepareNetINodeMetaDataOnlyLoadDB solofstypes.PrepareNetINodeMetaDataOnlyLoadDB,
	prepareNetINodeMetaDataWithStorDB solofstypes.PrepareNetINodeMetaDataWithStorDB,
	netINodeCommitSizeInDB solofstypes.NetINodeCommitSizeInDB,
) error {
	var err error

	p.SoloosEnv = soloosEnv

	p.SetHelper(
		prepareNetINodeMetaDataOnlyLoadDB, prepareNetINodeMetaDataWithStorDB,
		netINodeCommitSizeInDB)

	p.solonnClient = solonnClient

	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver

	err = p.OffheapDriver.InitLKVTableWithBytes64(&p.netINodeTable, "NetINode",
		int(solofstypes.NetINodeStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetINodeDriver) netINodeTablePrepareNewObjectFunc(uNetINode solofstypes.NetINodeUintptr,
	afterSetNewObj offheap.KVTableAfterSetNewObj) bool {
	var isNewObjectSetted bool
	if afterSetNewObj != nil {
		uNetINode.Ptr().NetINodeID = uNetINode.Ptr().LKVTableObjectWithBytes64.ID
		afterSetNewObj()
		isNewObjectSetted = true
	} else {
		isNewObjectSetted = false
	}
	return isNewObjectSetted
}

func (p *NetINodeDriver) SetHelper(
	prepareNetINodeMetaDataOnlyLoadDB solofstypes.PrepareNetINodeMetaDataOnlyLoadDB,
	prepareNetINodeMetaDataWithStorDB solofstypes.PrepareNetINodeMetaDataWithStorDB,
	netINodeCommitSizeInDB solofstypes.NetINodeCommitSizeInDB,
) {
	p.helper.PrepareNetINodeMetaDataOnlyLoadDB = prepareNetINodeMetaDataOnlyLoadDB
	p.helper.PrepareNetINodeMetaDataWithStorDB = prepareNetINodeMetaDataWithStorDB
	p.helper.NetINodeCommitSizeInDB = netINodeCommitSizeInDB
}

func (p *NetINodeDriver) SetSolonnClient(solonnClient *solofsapi.SolonnClient) {
	p.solonnClient = solonnClient
}

func (p *NetINodeDriver) NetINodeTruncate(uNetINode solofstypes.NetINodeUintptr, size uint64) error {
	return p.helper.NetINodeCommitSizeInDB(uNetINode, size)
}

func (p *NetINodeDriver) GetNetINode(netINodeID solofstypes.NetINodeID) (solofstypes.NetINodeUintptr, error) {
	var (
		uObject        offheap.LKVTableObjectUPtrWithBytes64
		uNetINode      solofstypes.NetINodeUintptr
		pNetINode      *solofstypes.NetINode
		afterSetNewObj offheap.KVTableAfterSetNewObj
		err            error
	)
	uObject, afterSetNewObj = p.netINodeTable.MustGetObject(netINodeID)
	p.netINodeTablePrepareNewObjectFunc(solofstypes.NetINodeUintptr(uObject), afterSetNewObj)
	uNetINode = solofstypes.NetINodeUintptr(uObject)
	pNetINode = uNetINode.Ptr()
	if pNetINode.IsDBMetaDataInited.Load() == solodbtypes.MetaDataStateUninited {
		pNetINode.IsDBMetaDataInited.LockContext()
		if pNetINode.IsDBMetaDataInited.Load() == solodbtypes.MetaDataStateUninited {
			err = p.helper.PrepareNetINodeMetaDataOnlyLoadDB(uNetINode)
		}
		pNetINode.IsDBMetaDataInited.UnlockContext()
	}

	if err != nil {
		p.netINodeTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(uNetINode))
		return 0, err
	}

	return uNetINode, nil
}

func (p *NetINodeDriver) MustGetNetINode(netINodeID solofstypes.NetINodeID,
	size uint64, netBlockCap int, memBlockCap int) (solofstypes.NetINodeUintptr, error) {
	var (
		uObject        offheap.LKVTableObjectUPtrWithBytes64
		uNetINode      solofstypes.NetINodeUintptr
		pNetINode      *solofstypes.NetINode
		afterSetNewObj offheap.KVTableAfterSetNewObj
		err            error
	)
	uObject, afterSetNewObj = p.netINodeTable.MustGetObject(netINodeID)
	p.netINodeTablePrepareNewObjectFunc(solofstypes.NetINodeUintptr(uObject), afterSetNewObj)
	uNetINode = solofstypes.NetINodeUintptr(uObject)
	pNetINode = uNetINode.Ptr()
	if pNetINode.IsDBMetaDataInited.Load() == solodbtypes.MetaDataStateUninited {
		pNetINode.IsDBMetaDataInited.LockContext()
		if pNetINode.IsDBMetaDataInited.Load() == solodbtypes.MetaDataStateUninited {
			err = p.helper.PrepareNetINodeMetaDataWithStorDB(uNetINode, size, netBlockCap, memBlockCap)
		}
		pNetINode.IsDBMetaDataInited.UnlockContext()
	}

	if err != nil {
		p.netINodeTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(uNetINode))
		return 0, err
	}

	return uNetINode, nil
}

func (p *NetINodeDriver) ReleaseNetINode(uNetINode solofstypes.NetINodeUintptr) {
	p.netINodeTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(uNetINode))
}
