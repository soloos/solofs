package namenode

import "soloos/sdfs/metastg"

type NameNode struct {
	MetaStg    metastg.MetaStg
	SRPCServer NameNodeSRPCServer
}

func (p *NameNode) Init(options NameNodeOptions) error {
	var err error

	p.MetaStg.Init(options.MetaStgDBDriver, options.MetaStgDBConnect)

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
