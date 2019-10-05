package solodn

import "soloos/common/iron"

type WebServer struct {
	solodn *Solodn
	server iron.Server
}

var _ = iron.IServer(&WebServer{})

func (p *WebServer) Init(solodn *Solodn,
	webServerOptions iron.Options,
) error {
	var err error
	p.solodn = solodn

	err = p.server.Init(webServerOptions)
	if err != nil {
		return err
	}

	return nil
}

func (p *WebServer) ServerName() string {
	return "Soloos.Solofs.Solodn.WebServer"
}

func (p *WebServer) Serve() error {
	return p.server.Serve()
}

func (p *WebServer) Close() error {
	return nil
}
