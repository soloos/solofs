package memstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/util/offheap"
)

type NetINodeDriver struct {
	offheapDriver  *offheap.OffheapDriver
	netBlockDriver *netstg.NetBlockDriver
	memBlockDriver *MemBlockDriver
	netINodePool   types.NetINodePool

	nameNodeClient *api.NameNodeClient
}

func (p *NetINodeDriver) Init(offheapDriver *offheap.OffheapDriver,
	netBlockDriver *netstg.NetBlockDriver,
	memBlockDriver *MemBlockDriver,
	nameNodeClient *api.NameNodeClient) error {
	p.offheapDriver = offheapDriver
	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver
	p.netINodePool.Init(-1, p.offheapDriver)
	p.nameNodeClient = nameNodeClient
	return nil
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
	if isLoaded == false || uNetINode.Ptr().IsMetaDataInited == false {
		pNetINode.MetaDataInitMutex.Lock()
		if pNetINode.IsMetaDataInited == false {
			err = p.PrepareNetINodeMetaDataWithStorDB(uNetINode, size, netBlockCap, memBlockCap)
			if err == nil {
				pNetINode.IsMetaDataInited = true
			}
		}
		pNetINode.MetaDataInitMutex.Unlock()
	}

	if err != nil {
		// TODO: clean uNetINode
		return 0, err
	}

	return uNetINode, nil
}

func (p *NetINodeDriver) PrepareNetINodeMetaDataWithStorDB(uNetINode types.NetINodeUintptr,
	size int64, netBlockCap int, memBlockCap int) error {
	var err error

	// do alloc
	err = p.nameNodeClient.AllocNetINodeMetaData(uNetINode, size, netBlockCap, memBlockCap)
	if err != nil {
		return err
	}

	return nil
}
