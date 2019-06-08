package datanode

import (
	"soloos/common/log"
	"soloos/common/sdfsapitypes"
	"soloos/common/snet"
)

type DataNodeSRPCServer struct {
	dataNode             *DataNode
	srpcServerListenAddr string
	srpcServerServeAddr  string
	srpcServer           snet.SRPCServer
}

func (p *DataNodeSRPCServer) Init(dataNode *DataNode,
	srpcServerListenAddr string,
	srpcServerServeAddr string,
) error {
	var err error

	p.dataNode = dataNode
	p.srpcServerListenAddr = srpcServerListenAddr
	err = p.srpcServer.Init(sdfsapitypes.DefaultSDFSRPCNetwork, p.srpcServerListenAddr)
	if err != nil {
		return err
	}

	p.srpcServer.RegisterService("/NetINode/PWrite", p.NetINodePWrite)
	p.srpcServer.RegisterService("/NetINode/Sync", p.NetINodeSync)
	p.srpcServer.RegisterService("/NetINode/PRead", p.NetINodePRead)

	return nil
}

func (p *DataNodeSRPCServer) Serve() error {
	log.Info("datanode srpcserver serve at:", p.srpcServerListenAddr)
	return p.srpcServer.Serve()
}

func (p *DataNodeSRPCServer) Close() error {
	log.Info("datanode srpcserver stop at:", p.srpcServerListenAddr)
	return p.srpcServer.Close()
}
