package types

import (
	"soloos/common/util"
	"soloos/sdbone/offheap"
)

type MockNetINodeTable struct {
	offheapDriver *offheap.OffheapDriver
	table         offheap.LKVTableWithBytes64
}

func (p *MockNetINodeTable) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error
	p.offheapDriver = offheapDriver

	err = p.offheapDriver.InitLKVTableWithBytes64(&p.table, "MockNetINode",
		int(NetINodeStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *MockNetINodeTable) MustGetNetINode(netINodeID NetINodeID) (NetINodeUintptr, bool) {
	u, afterSetNewObj := p.table.MustGetObjectWithAcquire(netINodeID)
	var loaded = afterSetNewObj == nil
	if afterSetNewObj != nil {
		afterSetNewObj()
	}
	uNetINode := (NetINodeUintptr)(u)
	return uNetINode, loaded
}

func (p *MockNetINodeTable) AllocNetINode(netBlockCap, memBlockCap int) NetINodeUintptr {
	var netINodeID NetINodeID
	util.InitUUID64(&netINodeID)
	uNetINode, _ := p.MustGetNetINode(netINodeID)
	uNetINode.Ptr().ID = netINodeID
	uNetINode.Ptr().NetBlockCap = netBlockCap
	uNetINode.Ptr().MemBlockCap = memBlockCap
	return uNetINode
}
