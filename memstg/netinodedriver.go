package memstg

import (
	"soloos/common/solodbapitypes"
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"soloos/solodb/offheap"
)

type NetINodeDriverHelper struct {
	solofsapitypes.PrepareNetINodeMetaDataOnlyLoadDB
	solofsapitypes.PrepareNetINodeMetaDataWithStorDB
	solofsapitypes.NetINodeCommitSizeInDB
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
	prepareNetINodeMetaDataOnlyLoadDB solofsapitypes.PrepareNetINodeMetaDataOnlyLoadDB,
	prepareNetINodeMetaDataWithStorDB solofsapitypes.PrepareNetINodeMetaDataWithStorDB,
	netINodeCommitSizeInDB solofsapitypes.NetINodeCommitSizeInDB,
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
		int(solofsapitypes.NetINodeStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetINodeDriver) netINodeTablePrepareNewObjectFunc(uNetINode solofsapitypes.NetINodeUintptr,
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
	prepareNetINodeMetaDataOnlyLoadDB solofsapitypes.PrepareNetINodeMetaDataOnlyLoadDB,
	prepareNetINodeMetaDataWithStorDB solofsapitypes.PrepareNetINodeMetaDataWithStorDB,
	netINodeCommitSizeInDB solofsapitypes.NetINodeCommitSizeInDB,
) {
	p.helper.PrepareNetINodeMetaDataOnlyLoadDB = prepareNetINodeMetaDataOnlyLoadDB
	p.helper.PrepareNetINodeMetaDataWithStorDB = prepareNetINodeMetaDataWithStorDB
	p.helper.NetINodeCommitSizeInDB = netINodeCommitSizeInDB
}

func (p *NetINodeDriver) SetSolonnClient(solonnClient *solofsapi.SolonnClient) {
	p.solonnClient = solonnClient
}

func (p *NetINodeDriver) NetINodeTruncate(uNetINode solofsapitypes.NetINodeUintptr, size uint64) error {
	return p.helper.NetINodeCommitSizeInDB(uNetINode, size)
}

func (p *NetINodeDriver) GetNetINode(netINodeID solofsapitypes.NetINodeID) (solofsapitypes.NetINodeUintptr, error) {
	var (
		uObject        offheap.LKVTableObjectUPtrWithBytes64
		uNetINode      solofsapitypes.NetINodeUintptr
		pNetINode      *solofsapitypes.NetINode
		afterSetNewObj offheap.KVTableAfterSetNewObj
		err            error
	)
	uObject, afterSetNewObj = p.netINodeTable.MustGetObject(netINodeID)
	p.netINodeTablePrepareNewObjectFunc(solofsapitypes.NetINodeUintptr(uObject), afterSetNewObj)
	uNetINode = solofsapitypes.NetINodeUintptr(uObject)
	pNetINode = uNetINode.Ptr()
	if pNetINode.IsDBMetaDataInited.Load() == solodbapitypes.MetaDataStateUninited {
		pNetINode.IsDBMetaDataInited.LockContext()
		if pNetINode.IsDBMetaDataInited.Load() == solodbapitypes.MetaDataStateUninited {
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

func (p *NetINodeDriver) MustGetNetINode(netINodeID solofsapitypes.NetINodeID,
	size uint64, netBlockCap int, memBlockCap int) (solofsapitypes.NetINodeUintptr, error) {
	var (
		uObject        offheap.LKVTableObjectUPtrWithBytes64
		uNetINode      solofsapitypes.NetINodeUintptr
		pNetINode      *solofsapitypes.NetINode
		afterSetNewObj offheap.KVTableAfterSetNewObj
		err            error
	)
	uObject, afterSetNewObj = p.netINodeTable.MustGetObject(netINodeID)
	p.netINodeTablePrepareNewObjectFunc(solofsapitypes.NetINodeUintptr(uObject), afterSetNewObj)
	uNetINode = solofsapitypes.NetINodeUintptr(uObject)
	pNetINode = uNetINode.Ptr()
	if pNetINode.IsDBMetaDataInited.Load() == solodbapitypes.MetaDataStateUninited {
		pNetINode.IsDBMetaDataInited.LockContext()
		if pNetINode.IsDBMetaDataInited.Load() == solodbapitypes.MetaDataStateUninited {
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

func (p *NetINodeDriver) ReleaseNetINode(uNetINode solofsapitypes.NetINodeUintptr) {
	p.netINodeTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(uNetINode))
}
