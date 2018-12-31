package main

import (
	"soloos/snet/srpc"
	"soloos/tinyiron"
)

type Env struct {
	WebServer  tinyiron.Server
	SRPCServer srpc.Server
}

func (p *Env) Init(srpcServerPort, webServerPort int) error {
	var (
		err error
	)

	err = p.initWebServer(webServerPort)
	if err != nil {
		return err
	}

	err = p.initSRPCServer(srpcServerPort)
	if err != nil {
		return err
	}

	return nil
}
