package datanode

import "soloos/common/iron"

type WebServer struct {
	dataNode *DataNode
	server   iron.Server
}

func (p *WebServer) Init(dataNode *DataNode,
	webServerOptions iron.Options,
) error {
	var err error
	p.dataNode = dataNode

	err = p.server.Init(webServerOptions)
	if err != nil {
		return err
	}

	return nil
}

func (p *WebServer) Serve() error {
	return p.server.Serve()
}

func (p *WebServer) Close() error {
	return nil
}
