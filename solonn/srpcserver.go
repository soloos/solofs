package solonn

import (
	"soloos/common/iron"
	"soloos/common/log"
	"soloos/common/solofsapitypes"
	"soloos/common/snet"
)

type SRPCServer struct {
	solonn             *Solonn
	srpcServerListenAddr string
	srpcServerServeAddr  string
	srpcServer           snet.SRPCServer
}

var _ = iron.IServer(&SRPCServer{})

func (p *SRPCServer) Init(solonn *Solonn,
	srpcServerListenAddr string,
	srpcServerServeAddr string,
) error {
	var err error

	p.solonn = solonn
	p.srpcServerListenAddr = srpcServerListenAddr
	p.srpcServerServeAddr = srpcServerServeAddr
	err = p.srpcServer.Init(solofsapitypes.DefaultSolofsRPCNetwork, p.srpcServerListenAddr)
	if err != nil {
		return err
	}

	p.srpcServer.RegisterService("/Solodn/Register", p.SolodnRegister)
	p.srpcServer.RegisterService("/NetINode/Get", p.NetINodeGet)
	p.srpcServer.RegisterService("/NetINode/MustGet", p.NetINodeMustGet)
	p.srpcServer.RegisterService("/NetINode/CommitSizeInDB", p.NetINodeCommitSizeInDB)
	p.srpcServer.RegisterService("/NetBlock/PrepareMetaData", p.NetBlockPrepareMetaData)

	return nil
}

func (p *SRPCServer) ServerName() string {
	return "Soloos.Solofs.Solonn.SRPCServer"
}

func (p *SRPCServer) Serve() error {
	log.Info("solonn srpcserver serve at:", p.srpcServerListenAddr)
	return p.srpcServer.Serve()
}

func (p *SRPCServer) Close() error {
	log.Info("solonn srpcserver stop at:", p.srpcServerListenAddr)
	return p.srpcServer.Close()
}
