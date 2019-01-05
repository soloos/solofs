package memstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/util/offheap"
)

type NetINodeDriver struct {
	Helper api.NetINodeDriverHelper

	offheapDriver *offheap.OffheapDriver
	netINodePool  types.NetINodePool

	memBlockDriver *MemBlockDriver
	netBlockDriver *netstg.NetBlockDriver
	netINodeDriver *NetINodeDriver
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
	p.netINodePool.Init(-1, p.offheapDriver)

	p.SetHelper(nameNodeClient, prepareNetINodeMetaDataOnlyLoadDB, prepareNetINodeMetaDataWithStorDB)

	return nil
}

func (p *NetINodeDriver) SetHelper(
	nameNodeClient *api.NameNodeClient,
	prepareNetINodeMetaDataOnlyLoadDB api.PrepareNetINodeMetaDataOnlyLoadDB,
	prepareNetINodeMetaDataWithStorDB api.PrepareNetINodeMetaDataWithStorDB,
) {
	p.Helper.NameNodeClient = nameNodeClient
	if prepareNetINodeMetaDataOnlyLoadDB != nil {
		p.Helper.PrepareNetINodeMetaDataOnlyLoadDB = prepareNetINodeMetaDataOnlyLoadDB
	} else {
		p.Helper.PrepareNetINodeMetaDataOnlyLoadDB = p.prepareNetINodeMetaDataOnlyLoadDB
	}
	if prepareNetINodeMetaDataWithStorDB != nil {
		p.Helper.PrepareNetINodeMetaDataWithStorDB = prepareNetINodeMetaDataWithStorDB
	} else {
		p.Helper.PrepareNetINodeMetaDataWithStorDB = p.prepareNetINodeMetaDataWithStorDB
	}
}

func (p *NetINodeDriver) GetNetINode(netINodeID types.NetINodeID) (types.NetINodeUintptr, error) {
	var (
		uNetINode types.NetINodeUintptr
		pNetINode *types.NetINode
		isLoaded  bool
		err       error
	)
	uNetINode, isLoaded = p.netINodePool.MustGetNetINode(netINodeID)
	pNetINode = uNetINode.Ptr()
	if isLoaded == false || uNetINode.Ptr().IsDBMetaDataInited == false {
		pNetINode.DBMetaDataInitMutex.Lock()
		if pNetINode.IsDBMetaDataInited == false {
			err = p.Helper.PrepareNetINodeMetaDataOnlyLoadDB(uNetINode)
			if err == nil {
				pNetINode.IsDBMetaDataInited = true
			}
		}
		pNetINode.DBMetaDataInitMutex.Unlock()
	}

	if err != nil {
		// TODO: clean uNetINode
		return 0, err
	}

	return uNetINode, nil
}

func (p *NetINodeDriver) MustGetNetINode(netINodeID types.NetINodeID,
	size int64, netBlockCap int, memBlockCap int) (types.NetINodeUintptr, error) {
	var (
		uNetINode types.NetINodeUintptr
		pNetINode *types.NetINode
		isLoaded  bool
		err       error
	)
	uNetINode, isLoaded = p.netINodePool.MustGetNetINode(netINodeID)
	pNetINode = uNetINode.Ptr()
	if isLoaded == false || uNetINode.Ptr().IsDBMetaDataInited == false {
		pNetINode.DBMetaDataInitMutex.Lock()
		if pNetINode.IsDBMetaDataInited == false {
			err = p.Helper.PrepareNetINodeMetaDataWithStorDB(uNetINode, size, netBlockCap, memBlockCap)
			if err == nil {
				pNetINode.IsDBMetaDataInited = true
			}
		}
		pNetINode.DBMetaDataInitMutex.Unlock()
	}

	if err != nil {
		// TODO: clean uNetINode
		return 0, err
	}

	return uNetINode, nil
}

func (p *NetINodeDriver) prepareNetINodeMetaDataOnlyLoadDB(uNetINode types.NetINodeUintptr) error {
	var err error

	err = p.Helper.NameNodeClient.GetNetINodeMetaData(uNetINode, -1, -1, -1)
	if err != nil {
		return err
	}

	return nil
}

func (p *NetINodeDriver) prepareNetINodeMetaDataWithStorDB(uNetINode types.NetINodeUintptr,
	size int64, netBlockCap int, memBlockCap int) error {
	var err error

	err = p.Helper.NameNodeClient.MustGetNetINodeMetaData(uNetINode, size, netBlockCap, memBlockCap)
	if err != nil {
		return err
	}

	return nil
}
