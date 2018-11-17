package datanode

type DataNode struct {
	SRPCServer DataNodeSRPCServer
}

func (p *DataNode) Init(options DataNodeOptions) error {
	var err error

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
