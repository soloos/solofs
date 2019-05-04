package namenode

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdbone/offheap"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
)

type NameNode struct {
	offheapDriver *offheap.OffheapDriver
	peerID        snettypes.PeerID
	metaStg       *metastg.MetaStg

	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *netstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver

	SRPCServer NameNodeSRPCServer
}

func (p *NameNode) Init(offheapDriver *offheap.OffheapDriver,
	srpcServerListenAddr string,
	peerID snettypes.PeerID,
	metaStg *metastg.MetaStg,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) error {
	var err error

	p.offheapDriver = offheapDriver
	p.peerID = peerID

	err = p.SRPCServer.Init(p, srpcServerListenAddr)
	if err != nil {
		return err
	}

	p.metaStg = metaStg
	p.memBlockDriver = memBlockDriver
	p.netBlockDriver = netBlockDriver
	p.netINodeDriver = netINodeDriver

	return nil
}

func (p *NameNode) RegisterDataNode(peerID snettypes.PeerID, serveAddr string) error {
	return p.metaStg.RegisterDataNode(peerID, serveAddr)
}

func (p *NameNode) Serve() error {
	return p.SRPCServer.Serve()
}

func (p *NameNode) Close() error {
	return p.SRPCServer.Close()
}
