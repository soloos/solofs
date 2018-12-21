package datanode

import "soloos/snet/srpc"

type DataNodeSRPCServer struct {
	dataNode   *DataNode
	srpcServer srpc.Server
}

func (p *DataNodeSRPCServer) Init(dataNode *DataNode, options DataNodeSRPCServerOptions) error {
	var err error

	p.dataNode = dataNode
	err = p.srpcServer.Init(options.Network, options.ListenAddr)
	if err != nil {
		return err
	}

	p.srpcServer.RegisterService("/NetINode/PWrite", p.NetINodePWrite)
	p.srpcServer.RegisterService("/NetINode/PRead", p.NetINodePRead)

	return nil
}

func (p *DataNodeSRPCServer) Serve() error {
	return p.srpcServer.Serve()
}

func (p *DataNodeSRPCServer) Close() error {
	return p.srpcServer.Close()
}
