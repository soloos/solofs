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
	peer    snettypes.Peer
	metaStg *metastg.MetaStg

	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *memstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver

	SRPCServer NameNodeSRPCServer
}

func (p *NameNode) initSNetPeer(peerID snettypes.PeerID, srpcServerServeAddr string) error {
	var err error
	p.peer.ID = peerID
	p.peer.SetAddress(srpcServerServeAddr)
	p.peer.ServiceProtocol = sdfsapitypes.DefaultSDFSRPCProtocol
	err = p.SNetDriver.RegisterPeer(p.peer)
	if err != nil {
		return err
	}
	return nil
}

func (p *NameNode) Init(soloOSEnv *soloosbase.SoloOSEnv,
	peerID snettypes.PeerID,
	srpcServerListenAddr string,
	srpcServerServeAddr string,
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

	err = p.SRPCServer.Init(p, srpcServerListenAddr, srpcServerServeAddr)
	if err != nil {
		return err
	}

	err = p.initSNetPeer(peerID, srpcServerServeAddr)
	if err != nil {
		return err
	}

	return nil
}

func (p *NameNode) DataNodeHeartBeat(peer snettypes.Peer) error {
	var err = p.SoloOSEnv.SNetDriver.RegisterPeer(peer)
	if err != nil {
		return err
	}

	return p.metaStg.DataNodeHeartBeat(peer)
}

func (p *NameNode) Serve() error {
	return p.SRPCServer.Serve()
}

func (p *NameNode) Close() error {
	return p.SRPCServer.Close()
}
