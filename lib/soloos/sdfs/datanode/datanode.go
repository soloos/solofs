package datanode

import (
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
	"soloos/util/offheap"
)

type DataNode struct {
	offheapDriver  *offheap.OffheapDriver
	metaStg        *metastg.MetaStg
	netBlockDriver *netstg.NetBlockDriver
	memBlockDriver *memstg.MemBlockDriver
	netINodeDriver *memstg.NetINodeDriver

	SRPCServer DataNodeSRPCServer
}

func (p *DataNode) Init(options DataNodeOptions,
	offheapDriver *offheap.OffheapDriver,
	metaStg *metastg.MetaStg,
	netBlockDriver *netstg.NetBlockDriver,
	memBlockDriver *memstg.MemBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) error {
	var err error

	p.offheapDriver = offheapDriver
	p.metaStg = metaStg

	err = p.SRPCServer.Init(p, options.SRPCServer)
	if err != nil {
		return err
	}

	return nil
}

func (p *DataNode) Serve() error {
	return p.SRPCServer.Serve()
}

func (p *DataNode) Close() error {
	return p.SRPCServer.Close()
}
