package memstg

import (
	"soloos/sdbone/offheap"
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
)

type NetINodeDriverHelper struct {
	*api.NameNodeClient
	api.PrepareNetINodeMetaDataOnlyLoadDB
	api.PrepareNetINodeMetaDataWithStorDB
	api.NetINodeCommitSizeInDB
}

type NetINodeDriver struct {
	helper NetINodeDriverHelper

	offheapDriver *offheap.OffheapDriver
	netINodeTable offheap.LKVTableWithBytes64

	memBlockDriver *MemBlockDriver
	netBlockDriver *netstg.NetBlockDriver
}

func (p *NetINodeDriver) Init(offheapDriver *offheap.OffheapDriver,
	netBlockDriver *netstg.NetBlockDriver,
	memBlockDriver *MemBlockDriver,
	// for NetINodeDriverHelper
	nameNodeClient *api.NameNodeClient,
	prepareNetINodeMetaDataOnlyLoadDB api.PrepareNetINodeMetaDataOnlyLoadDB,
	prepareNetINodeMetaDataWithStorDB api.PrepareNetINodeMetaDataWithStorDB,
	netINodeCommitSizeInDB api.NetINodeCommitSizeInDB,
) error {
	var err error

	p.offheapDriver = offheapDriver
	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver

	err = p.offheapDriver.InitLKVTableWithBytes64(&p.netINodeTable, "NetINode",
		int(types.NetINodeStructSize), -1, types.DefaultKVTableSharedCount,
		p.netINodeTablePrepareNewObjectFunc, nil)
	if err != nil {
		return err
	}

	p.SetHelper(nameNodeClient,
		prepareNetINodeMetaDataOnlyLoadDB, prepareNetINodeMetaDataWithStorDB,
		netINodeCommitSizeInDB)

	return nil
}

func (p *NetINodeDriver) netINodeTablePrepareNewObjectFunc(uObject uintptr) {
	types.NetINodeUintptr(uObject).Ptr().ID = types.NetINodeUintptr(uObject).Ptr().LKVTableObjectWithBytes64.ID
}

func (p *NetINodeDriver) SetHelper(
	nameNodeClient *api.NameNodeClient,
	prepareNetINodeMetaDataOnlyLoadDB api.PrepareNetINodeMetaDataOnlyLoadDB,
	prepareNetINodeMetaDataWithStorDB api.PrepareNetINodeMetaDataWithStorDB,
	netINodeCommitSizeInDB api.NetINodeCommitSizeInDB,
) {
	p.helper.NameNodeClient = nameNodeClient
	p.helper.PrepareNetINodeMetaDataOnlyLoadDB = prepareNetINodeMetaDataOnlyLoadDB
	p.helper.PrepareNetINodeMetaDataWithStorDB = prepareNetINodeMetaDataWithStorDB
	p.helper.NetINodeCommitSizeInDB = netINodeCommitSizeInDB
}

func (p *NetINodeDriver) GetNetINodeWithReadAcquire(isForceReload bool, netINodeID types.NetINodeID) (types.NetINodeUintptr, error) {
	var (
		uObject   uintptr
		uNetINode types.NetINodeUintptr
		pNetINode *types.NetINode
		loaded    bool
		err       error
	)
	uObject, loaded = p.netINodeTable.MustGetObjectWithAcquire(netINodeID)
	uNetINode = types.NetINodeUintptr(uObject)
	pNetINode = uNetINode.Ptr()
	if isForceReload == false &&
		(loaded == false || uNetINode.Ptr().IsDBMetaDataInited.Load() == types.MetaDataStateUninited) {
		pNetINode.DBMetaDataInitMutex.Lock()
		if pNetINode.IsDBMetaDataInited.Load() == types.MetaDataStateUninited {
			err = p.helper.PrepareNetINodeMetaDataOnlyLoadDB(uNetINode)
		}
		pNetINode.DBMetaDataInitMutex.Unlock()
	}

	if err != nil {
		p.ReleaseNetINodeWithReadRelease(uNetINode)
		return 0, err
	}

	return uNetINode, nil
}

func (p *NetINodeDriver) MustGetNetINodeWithReadAcquire(netINodeID types.NetINodeID,
	size uint64, netBlockCap int, memBlockCap int) (types.NetINodeUintptr, error) {
	var (
		uObject   uintptr
		uNetINode types.NetINodeUintptr
		pNetINode *types.NetINode
		loaded    bool
		err       error
	)
	uObject, loaded = p.netINodeTable.MustGetObjectWithAcquire(netINodeID)
	uNetINode = types.NetINodeUintptr(uObject)
	pNetINode = uNetINode.Ptr()
	if loaded == false || uNetINode.Ptr().IsDBMetaDataInited.Load() == types.MetaDataStateUninited {
		pNetINode.DBMetaDataInitMutex.Lock()
		if pNetINode.IsDBMetaDataInited.Load() == types.MetaDataStateUninited {
			err = p.helper.PrepareNetINodeMetaDataWithStorDB(uNetINode, size, netBlockCap, memBlockCap)
		}
		pNetINode.DBMetaDataInitMutex.Unlock()
	}

	if err != nil {
		p.ReleaseNetINodeWithReadRelease(uNetINode)
		return 0, err
	}

	return uNetINode, nil
}

func (p *NetINodeDriver) ReleaseNetINodeWithReadRelease(uNetINode types.NetINodeUintptr) {
	p.netINodeTable.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(uNetINode))
}

func (p *NetINodeDriver) NetINodeTruncate(uNetINode types.NetINodeUintptr, size uint64) {
	p.helper.NetINodeCommitSizeInDB(uNetINode, size)
}
