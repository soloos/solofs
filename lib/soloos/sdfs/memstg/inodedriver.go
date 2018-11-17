package memstg

import (
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
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
