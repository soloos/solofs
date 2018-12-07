package namenode

import (
	"soloos/sdfs/metastg"
	"soloos/sdfs/types"
	"soloos/snet"
	"soloos/util/offheap"
)

type NameNode struct {
	offheapDriver    *offheap.OffheapDriver
	snetDriver       snet.SNetDriver
	snetClientDriver snet.ClientDriver
	netBlockPool     types.NetBlockPool
	inodePool        types.INodePool

	MetaStg    metastg.MetaStg
	SRPCServer NameNodeSRPCServer
}

func (p *NameNode) Init(options NameNodeOptions,
	offheapDriver *offheap.OffheapDriver) error {
	var err error

	p.offheapDriver = offheapDriver

	err = p.MetaStg.Init(p.offheapDriver, options.MetaStgDBDriver, options.MetaStgDBConnect)
	if err != nil {
		return err
	}

	err = p.snetDriver.Init(p.offheapDriver)
	if err != nil {
		return err
	}

	err = p.snetClientDriver.Init(p.offheapDriver)
	if err != nil {
		return err
	}

	err = p.netBlockPool.Init(-1, p.offheapDriver)
	if err != nil {
		return err
	}

	err = p.inodePool.Init(-1, p.offheapDriver)
	if err != nil {
		return err
	}

	err = p.SRPCServer.Init(p, options.SRPCServer)
	if err != nil {
		return err
	}

	return nil
}

func (p *NameNode) Serve() error {
	return p.SRPCServer.Serve()
}

func (p *NameNode) Close() error {
	return p.SRPCServer.Close()
}
