package namenode

import "soloos/common/iron"

type WebServer struct {
	nameNode *NameNode
	server   iron.Server
}

func (p *WebServer) Init(nameNode *NameNode,
	webServerListenAddr string,
	webServerServeAddr string,
) error {
	var err error
	p.nameNode = nameNode

	var webOptions iron.Options
	webOptions.ListenStr = webServerListenAddr
	webOptions.ServeStr = webServerServeAddr
	err = p.server.Init(webOptions)
	if err != nil {
		return err
	}

	p.server.Router("/DataNode/HeartBeat", p.ctrDataNodeHeartBeat)

	return nil
}

func (p *WebServer) Serve() error {
	return p.server.Serve()
}

func (p *WebServer) Close() error {
	return nil
}
