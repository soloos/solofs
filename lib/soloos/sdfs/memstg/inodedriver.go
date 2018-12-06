package memstg

import (
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/util"
	"soloos/util/offheap"
)

type INodeDriver struct {
	offheapDriver  *offheap.OffheapDriver
	netBlockDriver *netstg.NetBlockDriver
	memBlockDriver *MemBlockDriver
	inodePool      INodePool
}

func (p *INodeDriver) Init(rawChunksLimit int32,
	offheapDriver *offheap.OffheapDriver,
	netBlockDriver *netstg.NetBlockDriver,
	memBlockDriver *MemBlockDriver) error {
	p.offheapDriver = offheapDriver
	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver
	p.inodePool.Init(rawChunksLimit, p)
	return nil
}

// MustGetINode get or init a inodeblock
func (p *INodeDriver) MustGetINode(inodeID types.INodeID) (types.INodeUintptr, bool) {
	return p.inodePool.MustGetINode(inodeID)
}

func (p *INodeDriver) InitINode(netBlockCap, memBlockCap int) types.INodeUintptr {
	var inodeID types.INodeID
	util.InitUUID64(&inodeID)
	uINode, _ := p.MustGetINode(inodeID)
	uINode.Ptr().ID = inodeID
	uINode.Ptr().NetBlockCap = netBlockCap
	uINode.Ptr().MemBlockCap = memBlockCap
	return uINode
}
