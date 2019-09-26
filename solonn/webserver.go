package solonn

import "soloos/common/iron"

type WebServer struct {
	solonn *Solonn
	server   iron.Server
}

var _ = iron.IServer(&WebServer{})

func (p *WebServer) Init(solonn *Solonn,
	webServerOptions iron.Options,
) error {
	var err error
	p.solonn = solonn

	err = p.server.Init(webServerOptions)
	if err != nil {
		return err
	}

	p.server.Router("/Solodn/HeartBeat", p.ctrSolodnHeartBeat)

	return nil
}

func (p *WebServer) ServerName() string {
	return "Soloos.Solofs.Solonn.WebServer"
}

func (p *WebServer) Serve() error {
	return p.server.Serve()
}

func (p *WebServer) Close() error {
	return nil
}
