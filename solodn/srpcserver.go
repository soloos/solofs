package solodn

import (
	"soloos/common/iron"
	"soloos/common/log"
	"soloos/common/solofsapitypes"
	"soloos/common/snet"
)

type SRPCServer struct {
	solodn             *Solodn
	srpcServerListenAddr string
	srpcServerServeAddr  string
	srpcServer           snet.SRPCServer
}

var _ = iron.IServer(&SRPCServer{})

func (p *SRPCServer) Init(solodn *Solodn,
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

func (p *SRPCServer) ServerName() string {
	return "Soloos.Solofs.Solodn.SRPCServer"
}

func (p *SRPCServer) Serve() error {
	log.Info("solodn srpcserver serve at:", p.srpcServerListenAddr)
	return p.srpcServer.Serve()
}

func (p *SRPCServer) Close() error {
	log.Info("solodn srpcserver stop at:", p.srpcServerListenAddr)
	return p.srpcServer.Close()
}
