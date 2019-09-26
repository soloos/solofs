package solodn

import (
	"fmt"
	"soloos/common/iron"
	"soloos/common/log"
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/snet"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/solofs/localfs"
	"soloos/solofs/memstg"
)

type Solodn struct {
	*soloosbase.SoloOSEnv
	srpcPeer snettypes.Peer
	webPeer  snettypes.Peer

	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *memstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver
	solonnClient solofsapi.SolonnClient

	localFS         localfs.LocalFS
	localFsSNetPeer snettypes.Peer

	heartBeatServerOptionsArr []snettypes.HeartBeatServerOptions
	srpcServer                SRPCServer
	webServer                 WebServer
	serverDriver              iron.ServerDriver
}

func (p *Solodn) initLocalFs(options SolodnOptions) error {
	var err error
	err = p.localFS.Init(options.LocalFSRoot)
	if err != nil {
		return err
	}

	p.localFsSNetPeer.ID = snet.MakeSysPeerID(fmt.Sprintf("SOLODN_LOCAL_FS"))
	p.localFsSNetPeer.SetAddress("LocalFs")
	p.localFsSNetPeer.ServiceProtocol = snettypes.ProtocolLocalFS
	err = p.SNetDriver.RegisterPeer(p.localFsSNetPeer)
	if err != nil {
		return err
	}

	return nil
}

func (p *Solodn) initSNetPeer(options SolodnOptions) error {
	var err error

	p.srpcPeer.ID = options.SRPCPeerID
	p.srpcPeer.SetAddress(options.SRPCServerServeAddr)
	p.srpcPeer.ServiceProtocol = solofsapitypes.DefaultSOLOFSRPCProtocol
	err = p.SNetDriver.RegisterPeer(p.srpcPeer)
	if err != nil {
		return err
	}

	p.webPeer.ID = options.WebPeerID
	p.webPeer.SetAddress(options.WebServer.ServeStr)
	p.webPeer.ServiceProtocol = snettypes.ProtocolWeb
	err = p.SNetDriver.RegisterPeer(p.webPeer)
	if err != nil {
		return err
	}

	return nil
}

func (p *Solodn) initNetBlockDriver() error {
	p.netBlockDriver.SetPReadMemBlockWithDisk(p.localFS.PReadMemBlockWithDisk)
	p.netBlockDriver.SetUploadMemBlockWithDisk(p.localFS.UploadMemBlockWithDisk)
	p.netBlockDriver.SetHelper(&p.solonnClient, p.netBlockDriver.PrepareNetBlockMetaData)
	return nil
}

func (p *Solodn) initNetINodeDriver() error {
	p.netINodeDriver.SetHelper(&p.solonnClient,
		p.netINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		p.netINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		p.netINodeDriver.NetINodeCommitSizeInDB,
	)
	return nil
}

func (p *Solodn) Init(soloOSEnv *soloosbase.SoloOSEnv,
	options SolodnOptions,
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
		log.Warn("Solodn Init initSNetPeer failed, err:", err)
		return err
	}

	err = p.srpcServer.Init(p, options.SRPCServerListenAddr, options.SRPCServerServeAddr)
	if err != nil {
		log.Warn("Solodn Init SRPCServer.Init failed, err:", err)
		return err
	}

	err = p.webServer.Init(p, options.WebServer)
	if err != nil {
		log.Warn("Solodn Init WebServer.Init failed, err:", err)
		return err
	}

	err = p.serverDriver.Init(&p.srpcServer, &p.webServer)
	if err != nil {
		return err
	}

	err = p.initLocalFs(options)
	if err != nil {
		log.Warn("Solodn Init initLocalFs failed, err:", err)
		return err
	}

	err = p.initNetBlockDriver()
	if err != nil {
		log.Warn("Solodn Init initNetBlockDriver failed, err:", err)
		return err
	}

	err = p.initNetINodeDriver()
	if err != nil {
		log.Warn("Solodn Init initNetINodeDriver failed, err:", err)
		return err
	}

	err = p.solonnClient.Init(p.SoloOSEnv, options.SolonnSRPCPeerID)
	if err != nil {
		log.Warn("Solodn Init solonnClient.Init failed, err:", err)
		return err
	}

	return nil
}

func (p *Solodn) Serve() error {
	var err error

	err = p.solonnClient.SolodnRegister(p.srpcPeer.ID, p.srpcPeer.AddressStr(), p.srpcPeer.ServiceProtocol)
	if err != nil {
		return err
	}

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

func (p *Solodn) Close() error {
	var err error
	err = p.serverDriver.Close()
	if err != nil {
		return err
	}

	return nil
}
