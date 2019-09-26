package solonn

import (
	"soloos/common/iron"
	"soloos/common/solofsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/solofs/memstg"
	"soloos/solofs/metastg"
)

type Solonn struct {
	*soloosbase.SoloosEnv
	srpcPeer snettypes.Peer
	webPeer  snettypes.Peer
	metaStg  *metastg.MetaStg

	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *memstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver

	heartBeatServerOptionsArr []snettypes.HeartBeatServerOptions
	serverCount               int
	srpcServer                SRPCServer
	webServer                 WebServer
	serverDriver              iron.ServerDriver
}

func (p *Solonn) initSNetPeer(
	srpcPeerID snettypes.PeerID, srpcServerServeAddr string,
	webPeerID snettypes.PeerID, webServerServeAddr string,
) error {
	var err error

	p.srpcPeer.ID = srpcPeerID
	p.srpcPeer.SetAddress(srpcServerServeAddr)
	p.srpcPeer.ServiceProtocol = solofsapitypes.DefaultSolofsRPCProtocol
	err = p.SNetDriver.RegisterPeer(p.srpcPeer)
	if err != nil {
		return err
	}

	p.webPeer.ID = webPeerID
	p.webPeer.SetAddress(webServerServeAddr)
	p.webPeer.ServiceProtocol = snettypes.ProtocolWeb
	err = p.SNetDriver.RegisterPeer(p.webPeer)
	if err != nil {
		return err
	}

	return nil
}

func (p *Solonn) Init(soloosEnv *soloosbase.SoloosEnv,
	srpcPeerID snettypes.PeerID,
	srpcServerListenAddr string,
	srpcServerServeAddr string,
	webPeerID snettypes.PeerID,
	webServerOptions iron.Options,
	metaStg *metastg.MetaStg,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *memstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) error {
	var err error

	p.SoloosEnv = soloosEnv
	p.metaStg = metaStg
	p.memBlockDriver = memBlockDriver
	p.netBlockDriver = netBlockDriver
	p.netINodeDriver = netINodeDriver

	err = p.srpcServer.Init(p, srpcServerListenAddr, srpcServerServeAddr)
	if err != nil {
		return err
	}

	err = p.webServer.Init(p, webServerOptions)
	if err != nil {
		return err
	}

	err = p.serverDriver.Init(&p.srpcServer, &p.webServer)
	if err != nil {
		return err
	}

	err = p.initSNetPeer(srpcPeerID, srpcServerServeAddr, webPeerID, webServerOptions.ServeStr)
	if err != nil {
		return err
	}

	return nil
}

func (p *Solonn) SolodnRegister(peer snettypes.Peer) error {
	var err = p.SoloosEnv.SNetDriver.RegisterPeer(peer)
	if err != nil {
		return err
	}

	return p.metaStg.SolodnRegister(peer)
}

func (p *Solonn) Serve() error {
	var err error

	err = p.StartHeartBeat()
	if err != nil {
		return err
	}

	err = p.serverDriver.Serve()
	if err != nil {
		return err
	}

	return nil
}

func (p *Solonn) Close() error {
	var err error
	err = p.serverDriver.Close()
	if err != nil {
		return err
	}

	return nil
}
