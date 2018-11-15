package netstg

import (
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
)

type NetBlockDriver struct {
	offheapDriver    *offheap.OffheapDriver
	netBlockPool     NetBlockPool
	snetClientDriver *snet.ClientDriver
	nameNodePeer     snettypes.PeerUintptr
}

func (p *NetBlockDriver) Init(options NetBlockDriverOptions,
	offheapDriver *offheap.OffheapDriver,
	snetClientDriver *snet.ClientDriver) error {
	p.offheapDriver = offheapDriver
	p.netBlockPool.Init(options.NetBlockPoolOptions, p)
	p.snetClientDriver = snetClientDriver
	return nil
}

// MustGetNetBlock get or init a netBlockblock
func (p *NetBlockDriver) MustGetBlock(uINode types.INodeUintptr,
	netBlockIndex int) (types.NetBlockUintptr, bool) {
	var netBlockID types.PtrBindIndex
	types.EncodePtrBindIndex(&netBlockID, uintptr(uINode), netBlockIndex)
	return p.netBlockPool.MustGetBlock(netBlockID)
}
