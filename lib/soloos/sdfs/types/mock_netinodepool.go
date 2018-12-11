package types

import (
	"soloos/util"
	"soloos/util/offheap"
)

type MockINodePool struct {
	offheapDriver *offheap.OffheapDriver
	pool          offheap.RawObjectPool
}

func (p *MockINodePool) Init(offheapDriver *offheap.OffheapDriver) error {
	var err error
	p.offheapDriver = offheapDriver

	err = p.offheapDriver.InitRawObjectPool(&p.pool,
		int(INodeStructSize), -1, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *MockINodePool) MustGetINode(inodeID INodeID) (INodeUintptr, bool) {
	u, loaded := p.pool.MustGetRawObject(inodeID)
	uINode := (INodeUintptr)(u)
	return uINode, loaded
}

func (p *MockINodePool) InitINode(netBlockCap, memBlockCap int) INodeUintptr {
	var inodeID INodeID
	util.InitUUID64(&inodeID)
	uINode, _ := p.MustGetINode(inodeID)
	uINode.Ptr().ID = inodeID
	uINode.Ptr().NetBlockCap = netBlockCap
	uINode.Ptr().MemBlockCap = memBlockCap
	return uINode
}
