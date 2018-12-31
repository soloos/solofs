package main

import (
	"fmt"
	"soloos/tinyiron"
)

func (p *Env) initWebServer(webServerPort int) error {
	var (
		err              error
		webServerOptions tinyiron.Options
	)
	webServerOptions.ListenStr = fmt.Sprintf("0.0.0.0:%d", webServerPort)
	err = p.WebServer.Init()
	if err != nil {
		return err
	}

	err = p.WebServer.LoadOptions(webServerOptions)
	if err != nil {
		return err
	}

	p.WebServer.Router("/NetBlock/Write", p.CtrNetBlockWrite)

	return nil
}
