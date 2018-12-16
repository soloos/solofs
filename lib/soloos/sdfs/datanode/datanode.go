package datanode

import (
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
)

type DataNode struct {
	offheapDriver  *offheap.OffheapDriver
	snetDriver     *snet.SNetDriver
	metaStg        *metastg.MetaStg
	netBlockDriver *netstg.NetBlockDriver
	memBlockDriver *memstg.MemBlockDriver

	uLocalDiskPeer snettypes.PeerUintptr

	SRPCServer DataNodeSRPCServer
}

func (p *DataNode) Init(options DataNodeOptions,
	offheapDriver *offheap.OffheapDriver,
	snetDriver *snet.SNetDriver,
	metaStg *metastg.MetaStg,
	netBlockDriver *netstg.NetBlockDriver,
	memBlockDriver *memstg.MemBlockDriver,
) error {
	var err error

	p.offheapDriver = offheapDriver
	p.snetDriver = snetDriver
	p.metaStg = metaStg
	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver

	err = p.SRPCServer.Init(p, options.SRPCServer)
	if err != nil {
		return err
	}

	var peerID snettypes.PeerID
	p.uLocalDiskPeer, _ = p.snetDriver.MustGetPeer(&peerID, "/tmp/testsdfs", snettypes.ProtocolDisk)

	return nil
}

func (p *DataNode) Serve() error {
	return p.SRPCServer.Serve()
}

func (p *DataNode) Close() error {
	return p.SRPCServer.Close()
}
