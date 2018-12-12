package memstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/util"
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

// MustGetNetINode get or init a netINodeblock
func (p *NetINodeDriver) MustGetNetINode(netINodeID types.NetINodeID) (types.NetINodeUintptr, bool) {
	return p.netINodePool.MustGetNetINode(netINodeID)
}

func (p *NetINodeDriver) AllocNetINode(size int64, netBlockCap int, memBlockCap int) (types.NetINodeUintptr, error) {
	var (
		netINodeID types.NetINodeID
		uNetINode  types.NetINodeUintptr
		exists     bool
		err        error
	)

	util.InitUUID64(&netINodeID)
	uNetINode, exists = p.MustGetNetINode(netINodeID)
	if exists {
		panic("netINode should not exists")
	}

	err = p.allocNetINodeMetadata(uNetINode, size, netBlockCap, memBlockCap)
	if err != nil {
		return 0, err
	}

	return uNetINode, nil
}

func (p *NetINodeDriver) allocNetINodeMetadata(uNetINode types.NetINodeUintptr,
	size int64, netBlockCap int, memBlockCap int) error {
	var (
		pNetINode = uNetINode.Ptr()
		err       error
	)

	pNetINode.MetaDataMutex.Lock()
	if pNetINode.IsMetaDataInited {
		goto PREPARE_DONE
	}

	err = p.nameNodeClient.AllocNetINodeMetadata(uNetINode, size, netBlockCap, memBlockCap)
	if err != nil {
		goto PREPARE_DONE
	}

	pNetINode.IsMetaDataInited = true

PREPARE_DONE:
	pNetINode.MetaDataMutex.Unlock()
	return err
}
