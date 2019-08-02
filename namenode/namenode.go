package namenode

import (
	"soloos/common/iron"
	"soloos/common/log"
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
)

type NameNode struct {
	*soloosbase.SoloOSEnv
	srpcPeer snettypes.Peer
	webPeer  snettypes.Peer
	metaStg  *metastg.MetaStg

	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *memstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver

	heartBeatServerOptionsArr []sdfsapitypes.HeartBeatServerOptions
	serverCount               int
	srpcServer                SRPCServer
	webServer                 WebServer
}

func (p *NameNode) initSNetPeer(
	srpcPeerID snettypes.PeerID, srpcServerServeAddr string,
	webPeerID snettypes.PeerID, webServerServeAddr string,
) error {
	var err error

	p.srpcPeer.ID = srpcPeerID
	p.srpcPeer.SetAddress(srpcServerServeAddr)
	p.srpcPeer.ServiceProtocol = sdfsapitypes.DefaultSDFSRPCProtocol
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

func (p *NameNode) Init(soloOSEnv *soloosbase.SoloOSEnv,
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

	p.SoloOSEnv = soloOSEnv
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

	err = p.initSNetPeer(srpcPeerID, srpcServerServeAddr, webPeerID, webServerOptions.ServeStr)
	if err != nil {
		return err
	}

	return nil
}

func (p *NameNode) DataNodeRegister(peer snettypes.Peer) error {
	var err = p.SoloOSEnv.SNetDriver.RegisterPeer(peer)
	if err != nil {
		return err
	}

	return p.metaStg.DataNodeRegister(peer)
}

func (p *NameNode) Serve() error {
	var (
		errChan chan error
		tmpErr  error
		err     error
	)

	err = p.StartHeartBeat()
	if err != nil {
		return err
	}

	errChan = make(chan error, p.serverCount)

	p.serverCount = 2

	go func(errChan chan<- error) {
		errChan <- p.srpcServer.Serve()
	}(errChan)

	go func(errChan chan<- error) {
		errChan <- p.webServer.Serve()
	}(errChan)

	for i := 0; i < p.serverCount; i++ {
		tmpErr = <-errChan
		if tmpErr != nil {
			log.Error("serve error, err:", tmpErr)
			err = tmpErr
		}
	}

	return err
}

func (p *NameNode) Close() error {
	var (
		tmpErr error
		err    error
	)

	for i := 0; i < p.serverCount; i++ {
		tmpErr = p.srpcServer.Close()
		if err != nil {
			log.Error("server close error, err:", tmpErr)
			err = tmpErr
		}
	}

	return err
}
