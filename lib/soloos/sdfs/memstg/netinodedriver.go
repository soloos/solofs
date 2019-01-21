package memstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/util/offheap"
)

type NetINodeDriverHelper struct {
	*api.NameNodeClient
	api.PrepareNetINodeMetaDataOnlyLoadDB
	api.PrepareNetINodeMetaDataWithStorDB
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
) error {
	p.offheapDriver = offheapDriver
	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver
	p.netINodePool.Init(p.offheapDriver)

	p.SetHelper(nameNodeClient, prepareNetINodeMetaDataOnlyLoadDB, prepareNetINodeMetaDataWithStorDB)

	return nil
}

func (p *NetINodeDriver) SetHelper(
	nameNodeClient *api.NameNodeClient,
	prepareNetINodeMetaDataOnlyLoadDB api.PrepareNetINodeMetaDataOnlyLoadDB,
	prepareNetINodeMetaDataWithStorDB api.PrepareNetINodeMetaDataWithStorDB,
) {
	p.helper.NameNodeClient = nameNodeClient
	p.helper.PrepareNetINodeMetaDataOnlyLoadDB = prepareNetINodeMetaDataOnlyLoadDB
	p.helper.PrepareNetINodeMetaDataWithStorDB = prepareNetINodeMetaDataWithStorDB
}

func (p *NetINodeDriver) GetNetINodeWithReadAcquire(netINodeID types.NetINodeID) (types.NetINodeUintptr, error) {
	var (
		uNetINode types.NetINodeUintptr
		pNetINode *types.NetINode
		isLoaded  bool
		err       error
	)
	uNetINode, isLoaded = p.netINodePool.MustGetNetINodeWithReadAcquire(netINodeID)
	pNetINode = uNetINode.Ptr()
	if isLoaded == false || uNetINode.Ptr().IsDBMetaDataInited == false {
		pNetINode.DBMetaDataInitMutex.Lock()
		if pNetINode.IsDBMetaDataInited == false {
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
	if isLoaded == false || uNetINode.Ptr().IsDBMetaDataInited == false {
		pNetINode.DBMetaDataInitMutex.Lock()
		if pNetINode.IsDBMetaDataInited == false {
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
