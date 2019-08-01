package namenode

import (
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

	srpcServer SRPCServer
	webServer  WebServer
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
	webServerListenAddr string,
	webServerServeAddr string,
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

	err = p.webServer.Init(p, webServerListenAddr, webServerServeAddr)
	if err != nil {
		return err
	}

	err = p.initSNetPeer(srpcPeerID, srpcServerServeAddr, webPeerID, webServerServeAddr)
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
	return p.srpcServer.Serve()
}

func (p *NameNode) Close() error {
	return p.srpcServer.Close()
}
