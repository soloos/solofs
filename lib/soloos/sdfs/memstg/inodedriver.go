package memstg

import (
	"soloos/sdfs/netstg"
	"soloos/util/offheap"
	"soloos/sdfs/types"
)

type INodeDriver struct {
	offheapDriver  *offheap.OffheapDriver
	netBlockDriver *netstg.NetBlockDriver
	memBlockDriver *MemBlockDriver
	inodePool       INodePool
}

func (p *INodeDriver) Init(options INodePoolOptions,
	offheapDriver *offheap.OffheapDriver,
	netBlockDriver *netstg.NetBlockDriver,
	memBlockDriver *MemBlockDriver) error {
	p.offheapDriver = offheapDriver
	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver
	p.inodePool.Init(options, p)
	return nil
}

// MustGetINode get or init a inodeblock
func (p *INodeDriver) MustGetINode(inodeID types.INodeID) (types.INodeUintptr, bool) {
	return p.inodePool.MustGetINode(inodeID)
}
