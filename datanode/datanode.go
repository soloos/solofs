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
	srpcPeer snettypes.Peer

	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *memstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver
	nameNodeClient sdfsapi.NameNodeClient

	localFS         localfs.LocalFS
	localFsSNetPeer snettypes.Peer

	srpcServer SRPCServer
	webServer  WebServer
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
	p.srpcPeer.ID = options.SRPCPeerID
	p.srpcPeer.SetAddress(options.SRPCServerServeAddr)
	p.srpcPeer.ServiceProtocol = sdfsapitypes.DefaultSDFSRPCProtocol
	err = p.SNetDriver.RegisterPeer(p.srpcPeer)
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

	err = p.srpcServer.Init(p, options.SRPCServerListenAddr, options.SRPCServerServeAddr)
	if err != nil {
		log.Warn("DataNode Init SRPCServer.Init failed, err:", err)
		return err
	}

	err = p.webServer.Init(p, options.WebServerListenAddr, options.WebServerServeAddr)
	if err != nil {
		log.Warn("DataNode Init WebServer.Init failed, err:", err)
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

	err = p.nameNodeClient.Init(p.SoloOSEnv, options.NameNodeSRPCPeerID)
	if err != nil {
		log.Warn("DataNode Init nameNodeClient.Init failed, err:", err)
		return err
	}

	return nil
}

func (p *DataNode) Serve() error {
	var err error
	err = p.nameNodeClient.DataNodeRegister(p.srpcPeer.ID, p.srpcPeer.AddressStr(), p.srpcPeer.ServiceProtocol)
	if err != nil {
		return err
	}

	return p.srpcServer.Serve()
}

func (p *DataNode) Close() error {
	return p.srpcServer.Close()
}
