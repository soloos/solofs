package memstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
)

type MemStg struct {
	offheapDriver    *offheap.OffheapDriver
	SnetDriver       snet.NetDriver
	SnetClientDriver snet.ClientDriver
	NameNodeClient   api.NameNodeClient
	DataNodeClient   api.DataNodeClient
	netstg.NetBlockDriver
	MemBlockDriver
	NetINodeDriver
}

func (p *MemStg) Init(offheapDriver *offheap.OffheapDriver,
	nameNodeSRPCServerAddr string,
	memBlockDriverOptions MemBlockDriverOptions,
) error {
	var (
		nameNodePeer snettypes.PeerUintptr
		err          error
	)

	p.offheapDriver = offheapDriver

	err = p.SnetDriver.Init(p.offheapDriver)
	if err != nil {
		return err
	}

	err = p.SnetClientDriver.Init(p.offheapDriver)
	if err != nil {
		return err
	}

	nameNodePeer, _ = p.SnetDriver.MustGetPeer(nil, nameNodeSRPCServerAddr, types.DefaultSDFSRPCProtocol)
	err = p.NameNodeClient.Init(&p.SnetClientDriver, nameNodePeer)
	if err != nil {
		return err
	}

	err = p.DataNodeClient.Init(&p.SnetClientDriver)
	if err != nil {
		return err
	}

	err = p.NetBlockDriver.Init(p.offheapDriver,
		&p.SnetDriver, &p.SnetClientDriver,
		&p.NameNodeClient, &p.DataNodeClient,
		nil,
	)
	if err != nil {
		return err
	}

	err = p.MemBlockDriver.Init(p.offheapDriver, memBlockDriverOptions)
	if err != nil {
		return err
	}

	err = p.NetINodeDriver.Init(p.offheapDriver, &p.NetBlockDriver, &p.MemBlockDriver,
		&p.NameNodeClient, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
