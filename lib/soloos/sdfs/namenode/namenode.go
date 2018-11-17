package namenode

type NameNode struct {
	SRPCServer NameNodeSRPCServer
}

func (p *NameNode) Init(options NameNodeOptions) error {
	var err error

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
