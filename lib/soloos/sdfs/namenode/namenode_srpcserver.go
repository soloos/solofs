package namenode

import "soloos/snet/srpc"

type NameNodeSRPCServer struct {
	nameNode   *NameNode
	srpcServer srpc.Server
}

func (p *NameNodeSRPCServer) Init(nameNode *NameNode, options NameNodeSRPCServerOptions) error {
	var err error

	p.nameNode = nameNode
	err = p.srpcServer.Init(options.Network, options.ListenAddr)
	if err != nil {
		return err
	}

	p.srpcServer.RegisterService("/NetBlock/PrepareMetadata", p.NetBlockPrepareMetadata)

	return nil
}

func (p *NameNodeSRPCServer) Serve() error {
	return p.srpcServer.Serve()
}

func (p *NameNodeSRPCServer) Close() error {
	return p.srpcServer.Close()
}
