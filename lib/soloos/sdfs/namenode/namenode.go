package namenode

import (
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
	snettypes "soloos/common/snet/types"
	"soloos/common/util/offheap"
)

type NameNode struct {
	offheapDriver *offheap.OffheapDriver
	metaStg       *metastg.MetaStg

	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *netstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver

	SRPCServer NameNodeSRPCServer
}

func (p *NameNode) Init(offheapDriver *offheap.OffheapDriver,
	srpcServerListenAddr string,
	metaStg *metastg.MetaStg,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) error {
	var err error

	p.offheapDriver = offheapDriver
	p.metaStg = metaStg
	p.memBlockDriver = memBlockDriver
	p.netBlockDriver = netBlockDriver
	p.netINodeDriver = netINodeDriver

	err = p.SRPCServer.Init(p, srpcServerListenAddr)
	if err != nil {
		return err
	}

	return nil
}

func (p *NameNode) RegisterDataNode(peerID *snettypes.PeerID, serveAddr string) (snettypes.PeerUintptr, error) {
	return p.metaStg.MustGetDataNode(peerID, serveAddr)
}

func (p *NameNode) Serve() error {
	return p.SRPCServer.Serve()
}

func (p *NameNode) Close() error {
	return p.SRPCServer.Close()
}
