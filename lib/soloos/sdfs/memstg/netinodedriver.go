package memstg

import (
	sdbapitypes "soloos/common/sdbapi/types"
	"soloos/common/sdfsapi"
	sdfsapitypes "soloos/common/sdfsapi/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/sdbone/offheap"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
)

type NetINodeDriverHelper struct {
	*sdfsapi.NameNodeClient
	sdfsapitypes.PrepareNetINodeMetaDataOnlyLoadDB
	sdfsapitypes.PrepareNetINodeMetaDataWithStorDB
	sdfsapitypes.NetINodeCommitSizeInDB
}

type NetINodeDriver struct {
	*soloosbase.SoloOSEnv
	helper NetINodeDriverHelper

	netINodeTable offheap.LKVTableWithBytes64

	memBlockDriver *MemBlockDriver
	netBlockDriver *netstg.NetBlockDriver
}

func (p *NetINodeDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	netBlockDriver *netstg.NetBlockDriver,
	memBlockDriver *MemBlockDriver,
	// for NetINodeDriverHelper
	nameNodeClient *sdfsapi.NameNodeClient,
	prepareNetINodeMetaDataOnlyLoadDB sdfsapitypes.PrepareNetINodeMetaDataOnlyLoadDB,
	prepareNetINodeMetaDataWithStorDB sdfsapitypes.PrepareNetINodeMetaDataWithStorDB,
	netINodeCommitSizeInDB sdfsapitypes.NetINodeCommitSizeInDB,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
	p.SetHelper(nameNodeClient,
		prepareNetINodeMetaDataOnlyLoadDB, prepareNetINodeMetaDataWithStorDB,
		netINodeCommitSizeInDB)

	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver

	err = p.OffheapDriver.InitLKVTableWithBytes64(&p.netINodeTable, "NetINode",
		int(types.NetINodeStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetINodeDriver) netINodeTablePrepareNewObjectFunc(uNetINode types.NetINodeUintptr,
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
	nameNodeClient *sdfsapi.NameNodeClient,
	prepareNetINodeMetaDataOnlyLoadDB sdfsapitypes.PrepareNetINodeMetaDataOnlyLoadDB,
	prepareNetINodeMetaDataWithStorDB sdfsapitypes.PrepareNetINodeMetaDataWithStorDB,
	netINodeCommitSizeInDB sdfsapitypes.NetINodeCommitSizeInDB,
) {
	p.helper.NameNodeClient = nameNodeClient
	p.helper.PrepareNetINodeMetaDataOnlyLoadDB = prepareNetINodeMetaDataOnlyLoadDB
	p.helper.PrepareNetINodeMetaDataWithStorDB = prepareNetINodeMetaDataWithStorDB
	p.helper.NetINodeCommitSizeInDB = netINodeCommitSizeInDB
}

func (p *NetINodeDriver) NetINodeTruncate(uNetINode types.NetINodeUintptr, size uint64) error {
	return p.helper.NetINodeCommitSizeInDB(uNetINode, size)
}

func (p *NetINodeDriver) GetNetINode(netINodeID types.NetINodeID) (types.NetINodeUintptr, error) {
	var (
		uObject        offheap.LKVTableObjectUPtrWithBytes64
		uNetINode      types.NetINodeUintptr
		pNetINode      *types.NetINode
		afterSetNewObj offheap.KVTableAfterSetNewObj
		err            error
	)
	uObject, afterSetNewObj = p.netINodeTable.MustGetObject(netINodeID)
	p.netINodeTablePrepareNewObjectFunc(types.NetINodeUintptr(uObject), afterSetNewObj)
	uNetINode = types.NetINodeUintptr(uObject)
	pNetINode = uNetINode.Ptr()
	if uNetINode.Ptr().IsDBMetaDataInited.Load() == sdbapitypes.MetaDataStateUninited {
		pNetINode.IsDBMetaDataInited.LockContext()
		if pNetINode.IsDBMetaDataInited.Load() == sdbapitypes.MetaDataStateUninited {
			err = p.helper.PrepareNetINodeMetaDataOnlyLoadDB(uNetINode)
		}
		pNetINode.IsDBMetaDataInited.UnlockContext()
	}

	if err != nil {
		p.netINodeTable.ForceDeleteAfterReleaseDone(offheap.LKVTableObjectUPtrWithBytes64(uNetINode))
		return 0, err
	}

	return uNetINode, nil
}

func (p *NetINodeDriver) MustGetNetINode(netINodeID types.NetINodeID,
	size uint64, netBlockCap int, memBlockCap int) (types.NetINodeUintptr, error) {
	var (
		uObject        offheap.LKVTableObjectUPtrWithBytes64
		uNetINode      types.NetINodeUintptr
		pNetINode      *types.NetINode
		afterSetNewObj offheap.KVTableAfterSetNewObj
		err            error
	)
	uObject, afterSetNewObj = p.netINodeTable.MustGetObject(netINodeID)
	p.netINodeTablePrepareNewObjectFunc(types.NetINodeUintptr(uObject), afterSetNewObj)
	uNetINode = types.NetINodeUintptr(uObject)
	pNetINode = uNetINode.Ptr()
	if uNetINode.Ptr().IsDBMetaDataInited.Load() == sdbapitypes.MetaDataStateUninited {
		pNetINode.IsDBMetaDataInited.LockContext()
		if pNetINode.IsDBMetaDataInited.Load() == sdbapitypes.MetaDataStateUninited {
			err = p.helper.PrepareNetINodeMetaDataWithStorDB(uNetINode, size, netBlockCap, memBlockCap)
		}
		pNetINode.IsDBMetaDataInited.UnlockContext()
	}

	if err != nil {
		p.netINodeTable.ForceDeleteAfterReleaseDone(offheap.LKVTableObjectUPtrWithBytes64(uNetINode))
		return 0, err
	}

	return uNetINode, nil
}

func (p *NetINodeDriver) ReleaseNetINode(uNetINode types.NetINodeUintptr) {
	p.netINodeTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(uNetINode))
}
