package types

import (
	"soloos/util"
	"soloos/util/offheap"
)

type MockNetINodePool struct {
	offheapDriver *offheap.OffheapDriver
	pool          offheap.RawObjectPool
}

func (p *MockNetINodePool) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error
	p.offheapDriver = offheapDriver

	err = p.offheapDriver.InitRawObjectPool(&p.pool,
		int(NetINodeStructSize), -1, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *MockNetINodePool) MustGetNetINode(netINodeID NetINodeID) (NetINodeUintptr, bool) {
	u, loaded := p.pool.MustGetRawObject(netINodeID)
	uNetINode := (NetINodeUintptr)(u)
	return uNetINode, loaded
}

func (p *MockNetINodePool) AllocNetINode(netBlockCap, memBlockCap int) NetINodeUintptr {
	var netINodeID NetINodeID
	util.InitUUID64(&netINodeID)
	uNetINode, _ := p.MustGetNetINode(netINodeID)
	uNetINode.Ptr().ID = netINodeID
	uNetINode.Ptr().NetBlockCap = netBlockCap
	uNetINode.Ptr().MemBlockCap = memBlockCap
	return uNetINode
}
