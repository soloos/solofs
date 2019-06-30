package datanode

import (
	"fmt"
	"soloos/common/log"
	"soloos/common/sdfsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/snet"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/sdfs/localfs"
	"soloos/sdfs/memstg"
)

type DataNode struct {
	*soloosbase.SoloOSEnv
	peer snettypes.Peer

	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *memstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver
	nameNodeClient sdfsapi.NameNodeClient

	localFS         localfs.LocalFS
	localFsSNetPeer snettypes.Peer

	SRPCServer DataNodeSRPCServer
}

func (p *DataNode) initLocalFs(options DataNodeOptions) error {
	var err error
	err = p.localFS.Init(options.LocalFSRoot)
	if err != nil {
		return err
	}

	p.localFsSNetPeer.ID = snet.MakeSysPeerID(fmt.Sprintf("DATANODE_LOCAL_FS"))
	p.localFsSNetPeer.SetAddress("LocalFs")
	p.localFsSNetPeer.ServiceProtocol = snettypes.ProtocolDisk
	err = p.SNetDriver.RegisterPeer(p.localFsSNetPeer)
	if err != nil {
		return err
	}

	return nil
}

func (p *DataNode) initSNetPeer(options DataNodeOptions) error {
	var err error
	p.peer.ID = options.PeerID
	p.peer.SetAddress(options.SrpcServerServeAddr)
	p.peer.ServiceProtocol = sdfsapitypes.DefaultSDFSRPCProtocol
	err = p.SNetDriver.RegisterPeer(p.peer)
	if err != nil {
		return err
	}
	return nil
}

func (p *DataNode) initNetBlockDriver() error {
	p.netBlockDriver.SetPReadMemBlockWithDisk(p.localFS.PReadMemBlockWithDisk)
	p.netBlockDriver.SetUploadMemBlockWithDisk(p.localFS.UploadMemBlockWithDisk)
	p.netBlockDriver.SetHelper(&p.nameNodeClient, p.netBlockDriver.PrepareNetBlockMetaData)
	return nil
}

func (p *DataNode) initNetINodeDriver() error {
	p.netINodeDriver.SetHelper(&p.nameNodeClient,
		p.netINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		p.netINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		p.netINodeDriver.NetINodeCommitSizeInDB,
	)
	return nil
}

func (p *DataNode) Init(soloOSEnv *soloosbase.SoloOSEnv,
	options DataNodeOptions,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *memstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
	p.netBlockDriver = netBlockDriver
	p.memBlockDriver = memBlockDriver
	p.netINodeDriver = netINodeDriver

	err = p.initSNetPeer(options)
	if err != nil {
		log.Warn("DataNode Init initSNetPeer failed, err:", err)
		return err
	}

	err = p.SRPCServer.Init(p, options.SrpcServerListenAddr, options.SrpcServerServeAddr)
	if err != nil {
		log.Warn("DataNode Init SRPCServer.Init failed, err:", err)
		return err
	}

	err = p.initLocalFs(options)
	if err != nil {
		log.Warn("DataNode Init initLocalFs failed, err:", err)
		return err
	}

	err = p.initNetBlockDriver()
	if err != nil {
		log.Warn("DataNode Init initNetBlockDriver failed, err:", err)
		return err
	}

	err = p.initNetINodeDriver()
	if err != nil {
		log.Warn("DataNode Init initNetINodeDriver failed, err:", err)
		return err
	}

	err = p.nameNodeClient.Init(p.SoloOSEnv, options.NameNodePeerID)
	if err != nil {
		log.Warn("DataNode Init nameNodeClient.Init failed, err:", err)
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
