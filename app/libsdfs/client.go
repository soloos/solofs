package main

import (
	"soloos/sdfs/api"
	"soloos/sdfs/memstg"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
)

type Client struct {
	offheapDriver *offheap.OffheapDriver

	snetDriver       snet.SNetDriver
	snetClientDriver snet.ClientDriver

	dataNodeClient api.DataNodeClient
	nameNodePeer   snettypes.PeerUintptr
	nameNodeClient api.NameNodeClient

	memBlockDriver memstg.MemBlockDriver
	netBlockDriver netstg.NetBlockDriver
	netINodeDriver memstg.NetINodeDriver
}

func (p *Client) Init(nameNodeSRPCServerAddr string,
	memBlockChunkSize int, memBlockChunksLimit int32) error {
	var err error
	p.offheapDriver = &offheap.DefaultOffheapDriver

	err = p.snetDriver.Init(p.offheapDriver)
	if err != nil {
		return err
	}

	err = p.snetClientDriver.Init(p.offheapDriver)
	if err != nil {
		return err
	}

	p.nameNodePeer, _ = p.snetDriver.MustGetPeer(nil, nameNodeSRPCServerAddr, types.DefaultSDFSRPCProtocol)
	err = p.nameNodeClient.Init(&p.snetClientDriver, p.nameNodePeer)
	if err != nil {
		return err
	}

	err = p.dataNodeClient.Init(&p.snetClientDriver, nil, nil)
	if err != nil {
		return err
	}

	var memBlockDriverOptions = memstg.MemBlockDriverOptions{
		MemBlockPoolOptionsList: []memstg.MemBlockPoolOptions{
			memstg.MemBlockPoolOptions{
				memBlockChunkSize,
				memBlockChunksLimit,
			},
		},
	}
	err = p.memBlockDriver.Init(memBlockDriverOptions, p.offheapDriver)
	if err != nil {
		return err
	}

	err = p.netBlockDriver.Init(p.offheapDriver, &p.snetDriver, &p.snetClientDriver,
		&p.nameNodeClient, &p.dataNodeClient,
		nil)
	if err != nil {
		return err
	}

	err = p.netINodeDriver.Init(p.offheapDriver,
		&p.netBlockDriver, &p.memBlockDriver, &p.nameNodeClient,
		nil, nil)
	if err != nil {
		return err
	}

	return nil
}
