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
	netINodePool  types.NetINodePool

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
	p.offheapDriver = offheapDriver
	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver
	p.netINodePool.Init(p.offheapDriver)

	p.SetHelper(nameNodeClient,
		prepareNetINodeMetaDataOnlyLoadDB, prepareNetINodeMetaDataWithStorDB,
		netINodeCommitSizeInDB)

	return nil
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
		uNetINode types.NetINodeUintptr
		pNetINode *types.NetINode
		isLoaded  bool
		err       error
	)
	uNetINode, isLoaded = p.netINodePool.MustGetNetINodeWithReadAcquire(netINodeID)
	pNetINode = uNetINode.Ptr()
	if isForceReload == false &&
		(isLoaded == false || uNetINode.Ptr().IsDBMetaDataInited.Load() == types.MetaDataStateUninited) {
		pNetINode.DBMetaDataInitMutex.Lock()
		if pNetINode.IsDBMetaDataInited.Load() == types.MetaDataStateUninited {
			err = p.helper.PrepareNetINodeMetaDataOnlyLoadDB(uNetINode)
		}
		pNetINode.DBMetaDataInitMutex.Unlock()
	}

	if err != nil {
		pNetINode.SharedPointer.SetReleasable()
		p.ReleaseNetINodeWithReadRelease(uNetINode)
		return 0, err
	}

	return uNetINode, nil
}

func (p *NetINodeDriver) MustGetNetINodeWithReadAcquire(netINodeID types.NetINodeID,
	size uint64, netBlockCap int, memBlockCap int) (types.NetINodeUintptr, error) {
	var (
		uNetINode types.NetINodeUintptr
		pNetINode *types.NetINode
		isLoaded  bool
		err       error
	)
	uNetINode, isLoaded = p.netINodePool.MustGetNetINodeWithReadAcquire(netINodeID)
	pNetINode = uNetINode.Ptr()
	if isLoaded == false || uNetINode.Ptr().IsDBMetaDataInited.Load() == types.MetaDataStateUninited {
		pNetINode.DBMetaDataInitMutex.Lock()
		if pNetINode.IsDBMetaDataInited.Load() == types.MetaDataStateUninited {
			err = p.helper.PrepareNetINodeMetaDataWithStorDB(uNetINode, size, netBlockCap, memBlockCap)
		}
		pNetINode.DBMetaDataInitMutex.Unlock()
	}

	if err != nil {
		pNetINode.SharedPointer.SetReleasable()
		p.ReleaseNetINodeWithReadRelease(uNetINode)
		return 0, err
	}

	return uNetINode, nil
}

func (p *NetINodeDriver) ReleaseNetINodeWithReadRelease(uNetINode types.NetINodeUintptr) {
	p.netINodePool.ReleaseNetINodeWithReadRelease(uNetINode)
}

func (p *NetINodeDriver) NetINodeTruncate(uNetINode types.NetINodeUintptr, size uint64) {
	p.helper.NetINodeCommitSizeInDB(uNetINode, size)
}
