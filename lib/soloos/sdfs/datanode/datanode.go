package datanode

import (
	"soloos/sdfs/localfs"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
)

type DataNode struct {
	offheapDriver *offheap.OffheapDriver
	snetDriver    *snet.SNetDriver
	metaStg       *metastg.MetaStg

	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *netstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver

	localFs        localfs.LocalFs
	uLocalDiskPeer snettypes.PeerUintptr

	SRPCServer DataNodeSRPCServer
}

func (p *DataNode) Init(options DataNodeOptions,
	offheapDriver *offheap.OffheapDriver,
	snetDriver *snet.SNetDriver,
	metaStg *metastg.MetaStg,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) error {
	var err error

	p.offheapDriver = offheapDriver
	p.snetDriver = snetDriver
	p.metaStg = metaStg
	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver
	p.netINodeDriver = netINodeDriver

	err = p.SRPCServer.Init(p, options.SRPCServer)
	if err != nil {
		return err
	}

	var peerID snettypes.PeerID
	err = p.localFs.Init("/tmp/testsdfs")
	if err != nil {
		return err
	}
	p.uLocalDiskPeer, _ = p.snetDriver.MustGetPeer(&peerID, "", snettypes.ProtocolDisk)
	p.netBlockDriver.SetUploadMemBlockWithDisk(p.localFs.UploadMemBlockWithDisk)

	return nil
}

func (p *DataNode) Serve() error {
	return p.SRPCServer.Serve()
}

func (p *DataNode) Close() error {
	return p.SRPCServer.Close()
}
