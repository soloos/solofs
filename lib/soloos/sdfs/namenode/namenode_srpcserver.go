package namenode

import (
	"soloos/log"
	"soloos/sdfs/types"
	"soloos/snet/srpc"
)

type NameNodeSRPCServer struct {
	nameNode             *NameNode
	srpcServerListenAddr string
	srpcServer           srpc.Server
}

func (p *NameNodeSRPCServer) Init(nameNode *NameNode, srpcServerListenAddr string) error {
	var err error

	p.nameNode = nameNode
	p.srpcServerListenAddr = srpcServerListenAddr
	err = p.srpcServer.Init(types.DefaultSDFSRPCNetwork, p.srpcServerListenAddr)
	if err != nil {
		return err
	}

	p.srpcServer.RegisterService("/DataNode/Register", p.DataNodeRegister)
	p.srpcServer.RegisterService("/NetINode/Get", p.NetINodeGet)
	p.srpcServer.RegisterService("/NetINode/MustGet", p.NetINodeMustGet)
	p.srpcServer.RegisterService("/NetBlock/PrepareMetaData", p.NetBlockPrepareMetaData)

	return nil
}

func (p *NameNodeSRPCServer) Serve() error {
	log.Info("namenode srpcserver serve at:", p.srpcServerListenAddr)
	return p.srpcServer.Serve()
}

func (p *NameNodeSRPCServer) Close() error {
	log.Info("namenode srpcserver stop at:", p.srpcServerListenAddr)
	return p.srpcServer.Close()
}
