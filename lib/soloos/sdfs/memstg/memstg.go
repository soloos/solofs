package memstg

import (
	"soloos/common/snet"
	snettypes "soloos/common/snet/types"
	"soloos/sdbone/offheap"
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
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

	err = p.SnetDriver.Init(p.offheapDriver, "MemStgNetDriver")
	if err != nil {
		return err
	}

	nameNodePeer = p.SnetDriver.AllocPeer(nameNodeSRPCServerAddr, types.DefaultSDFSRPCProtocol)
	err = p.NameNodeClient.Init(&p.SnetClientDriver, nameNodePeer)
	if err != nil {
		return err
	}

	err = p.SnetClientDriver.Init(p.offheapDriver)
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
		p.NetBlockDriver.PrepareNetBlockMetaDataWithFanout,
	)
	if err != nil {
		return err
	}

	err = p.MemBlockDriver.Init(p.offheapDriver, memBlockDriverOptions)
	if err != nil {
		return err
	}

	err = p.NetINodeDriver.Init(p.offheapDriver, &p.NetBlockDriver, &p.MemBlockDriver,
		&p.NameNodeClient,
		p.NetINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		p.NetINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		p.NetINodeDriver.NetINodeCommitSizeInDB,
	)
	if err != nil {
		return err
	}

	return nil
}
