package solodn

import (
	"soloos/common/iron"
	"soloos/common/log"
	"soloos/common/snet"
	"soloos/common/solofsapitypes"
)

type SrpcServer struct {
	solodn               *Solodn
	srpcServerListenAddr string
	srpcServerServeAddr  string
	srpcServer           snet.SrpcServer
}

var _ = iron.IServer(&SrpcServer{})

func (p *SrpcServer) Init(solodn *Solodn,
	srpcServerListenAddr string,
	srpcServerServeAddr string,
) error {
	var err error

	p.solodn = solodn
	p.srpcServerListenAddr = srpcServerListenAddr
	err = p.srpcServer.Init(solofsapitypes.DefaultSolofsRPCNetwork, p.srpcServerListenAddr)
	if err != nil {
		return err
	}

	p.srpcServer.RegisterService("/NetINode/PWrite", p.NetINodePWrite)
	p.srpcServer.RegisterService("/NetINode/Sync", p.NetINodeSync)
	p.srpcServer.RegisterService("/NetINode/PRead", p.NetINodePRead)

	return nil
}

func (p *SrpcServer) ServerName() string {
	return "Soloos.Solofs.Solodn.SrpcServer"
}

func (p *SrpcServer) Serve() error {
	log.Info("solodn srpcserver serve at:", p.srpcServerListenAddr)
	return p.srpcServer.Serve()
}

func (p *SrpcServer) Close() error {
	log.Info("solodn srpcserver stop at:", p.srpcServerListenAddr)
	return p.srpcServer.Close()
}
