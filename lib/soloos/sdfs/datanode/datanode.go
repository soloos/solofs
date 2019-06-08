package datanode

import (
	"fmt"
	"soloos/common/sdfsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/snet"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/sdfs/localfs"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
)

type DataNode struct {
	*soloosbase.SoloOSEnv
	peer    snettypes.Peer
	metaStg *metastg.MetaStg

	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *netstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver
	nameNodeClient sdfsapi.NameNodeClient

	localFS         localfs.LocalFS
	localFsSNetPeer snettypes.Peer

	SRPCServer DataNodeSRPCServer
}

func (p *DataNode) Init(soloOSEnv *soloosbase.SoloOSEnv,
	options DataNodeOptions,
	metaStg *metastg.MetaStg,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) error {
	var (
		err error
	)

	p.SoloOSEnv = soloOSEnv
	p.peer.ID = options.PeerID
	p.peer.SetAddress(options.SrpcServerServeAddr)
	p.peer.ServiceProtocol = sdfsapitypes.DefaultSDFSRPCProtocol

	p.metaStg = metaStg
	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver
	p.netINodeDriver = netINodeDriver

	err = p.SRPCServer.Init(p, options.SrpcServerListenAddr, options.SrpcServerServeAddr)
	if err != nil {
		return err
	}

	err = p.localFS.Init(options.LocalFSRoot)
	if err != nil {
		return err
	}

	p.localFsSNetPeer.ID = snet.MakeSysPeerID(fmt.Sprintf("DATANODE_"))
	p.localFsSNetPeer.SetAddress("LocalFs")
	p.localFsSNetPeer.ServiceProtocol = snettypes.ProtocolDisk
	err = p.SNetDriver.RegisterPeer(p.localFsSNetPeer)
	if err != nil {
		return err
	}

	p.netBlockDriver.SetPReadMemBlockWithDisk(p.localFS.PReadMemBlockWithDisk)
	p.netBlockDriver.SetUploadMemBlockWithDisk(p.localFS.UploadMemBlockWithDisk)
	p.netBlockDriver.SetHelper(&p.nameNodeClient, p.netBlockDriver.PrepareNetBlockMetaData)

	p.netINodeDriver.SetHelper(nil,
		p.metaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		p.metaStg.PrepareNetINodeMetaDataWithStorDB,
		p.metaStg.NetINodeCommitSizeInDB,
	)

	err = p.SNetDriver.RegisterPeer(p.peer)
	if err != nil {
		return err
	}

	err = p.nameNodeClient.Init(p.SoloOSEnv, options.NameNodePeerID)
	if err != nil {
		return err
	}

	return nil
}

func (p *DataNode) Serve() error {
	var err error
	err = p.nameNodeClient.RegisterDataNode(p.peer.ID, p.peer.AddressStr(), p.peer.ServiceProtocol)
	if err != nil {
		return err
	}

	return p.SRPCServer.Serve()
}

func (p *DataNode) Close() error {
	return p.SRPCServer.Close()
}
