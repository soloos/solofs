package datanode

import (
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/snet"
	"soloos/util/offheap"
)

type DataNode struct {
	offheapDriver    *offheap.OffheapDriver
	snetDriver       *snet.SNetDriver
	snetClientDriver *snet.ClientDriver
	netINodeDriver   *memstg.NetINodeDriver
	metaStg          *metastg.MetaStg

	SRPCServer DataNodeSRPCServer
}

func (p *DataNode) Init(options DataNodeOptions, offheapDriver *offheap.OffheapDriver) error {
	var err error

	p.offheapDriver = offheapDriver

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
