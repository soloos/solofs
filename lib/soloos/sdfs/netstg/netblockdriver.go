package netstg

import (
	"soloos/sdfs/types"
	"soloos/snet"
	"soloos/util/offheap"
)

type NetBlockDriver struct {
	offheapDriver *offheap.OffheapDriver
	netBlockPool  NetBlockPool

	netBlockDriverUploader netBlockDriverUploader
}

func (p *NetBlockDriver) Init(options NetBlockDriverOptions,
	offheapDriver *offheap.OffheapDriver,
	snetDriver *snet.SNetDriver,
	snetClientDriver *snet.ClientDriver) error {
	var err error
	p.offheapDriver = offheapDriver
	err = p.netBlockPool.Init(options.NetBlockPoolOptions, p)
	if err != nil {
		return err
	}

	err = p.netBlockDriverUploader.Init(p, snetDriver, snetClientDriver)
	if err != nil {
		return err
	}

	return nil
}

// MustGetNetBlock get or init a netBlockblock
func (p *NetBlockDriver) MustGetBlock(uINode types.INodeUintptr,
	netBlockIndex int) (types.NetBlockUintptr, bool) {
	var netBlockID types.PtrBindIndex
	types.EncodePtrBindIndex(&netBlockID, uintptr(uINode), netBlockIndex)
	return p.netBlockPool.MustGetBlock(netBlockID)
}
