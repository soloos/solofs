package datanode

import "soloos/common/iron"

type WebServer struct {
	dataNode *DataNode
	server   iron.Server
}

func (p *WebServer) Init(dataNode *DataNode,
	webServerListenAddr string,
	webServerServeAddr string,
) error {
	var err error
	p.dataNode = dataNode

	var webOptions iron.Options
	webOptions.ListenStr = webServerListenAddr
	webOptions.ServeStr = webServerServeAddr
	err = p.server.Init(webOptions)
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
