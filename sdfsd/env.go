package main

import (
	"soloos/common/sdfsapi"
	"soloos/common/snet"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/sdbone/offheap"
	"soloos/sdfs/datanode"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/namenode"
	"soloos/soloboat/soloboatsdk"
	"time"
)

type Env struct {
	options          Options
	SoloOSEnv        soloosbase.SoloOSEnv
	offheapDriver    *offheap.OffheapDriver
	SNetDriver       snet.NetDriver
	SNetClientDriver snet.SRPCClientDriver
	MetaStg          metastg.MetaStg
	DataNodeClient   sdfsapi.DataNodeClient
	MemBlockDriver   memstg.MemBlockDriver
	NetBlockDriver   memstg.NetBlockDriver
	NetINodeDriver   memstg.NetINodeDriver

	nameNode   namenode.NameNode
	dataNode   datanode.DataNode
	srpcPeerID snettypes.PeerID
	webPeerID  snettypes.PeerID

	soloboatClient          soloboatsdk.Client
	soloboatCronJobDuration time.Duration
}

func (p *Env) initMetaStg() error {
	return p.MetaStg.Init(&p.SoloOSEnv,
		p.options.DBDriver, p.options.Dsn)
}

func (p *Env) initMemStg() error {
	var memBlockDriverOptions = memstg.MemBlockDriverOptions{
		[]memstg.MemBlockTableOptions{
			memstg.MemBlockTableOptions{
				p.options.DefaultMemBlockCap,
				p.options.DefaultMemBlocksLimit,
			},
		},
	}
	return (p.MemBlockDriver.Init(&p.SoloOSEnv, memBlockDriverOptions))
}

func (p *Env) Init(options Options) {
	p.options = options
	util.AssertErrIsNil(p.SoloOSEnv.InitWithSNet(p.options.SNetDriverServeAddr))

	p.DataNodeClient.Init(&p.SoloOSEnv)

	util.AssertErrIsNil(p.initMetaStg())

	util.AssertErrIsNil(p.initMemStg())

	util.AssertErrIsNil(p.initSoloBoat())
}

func (p *Env) startCommon() {
	if p.options.PProfListenAddr != "" {
		go util.PProfServe(p.options.PProfListenAddr)
	}
}

func (p *Env) startNameNode() {
	copy(p.srpcPeerID[:], []byte(p.options.NameNodeSRPCPeerID))
	copy(p.webPeerID[:], []byte(p.options.NameNodeWebPeerID))

	util.AssertErrIsNil(p.NetBlockDriver.Init(&p.SoloOSEnv,
		nil, &p.DataNodeClient, p.MetaStg.PrepareNetBlockMetaData))

	util.AssertErrIsNil(p.NetINodeDriver.Init(&p.SoloOSEnv,
		&p.NetBlockDriver,
		&p.MemBlockDriver,
		nil,
		p.MetaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		p.MetaStg.PrepareNetINodeMetaDataWithStorDB,
		p.MetaStg.NetINodeCommitSizeInDB,
	))

	util.AssertErrIsNil(p.nameNode.Init(&p.SoloOSEnv,
		p.srpcPeerID, p.options.SRPCListenAddr, p.options.SRPCServeAddr,
		p.webPeerID, p.options.WebListenAddr, p.options.WebServeAddr,
		&p.MetaStg,
		&p.MemBlockDriver,
		&p.NetBlockDriver,
		&p.NetINodeDriver,
	))

	util.AssertErrIsNil(p.nameNode.Serve())
	util.AssertErrIsNil(p.nameNode.Close())
}

func (p *Env) startDataNode() {
	var (
		nameNodeSRPCPeerID snettypes.PeerID
		dataNodeOptions    datanode.DataNodeOptions
	)

	copy(p.srpcPeerID[:], []byte(p.options.DataNodeSRPCPeerID))
	copy(nameNodeSRPCPeerID[:], []byte(p.options.NameNodeSRPCPeerID))

	dataNodeOptions = datanode.DataNodeOptions{
		SRPCPeerID:           p.srpcPeerID,
		SRPCServerListenAddr: p.options.SRPCListenAddr,
		SRPCServerServeAddr:  p.options.SRPCServeAddr,
		WebServerListenAddr:  p.options.WebListenAddr,
		WebServerServeAddr:   p.options.WebServeAddr,
		LocalFSRoot:          p.options.DataNodeLocalFSRoot,
		NameNodeSRPCPeerID:   nameNodeSRPCPeerID,
	}

	util.AssertErrIsNil(p.NetBlockDriver.Init(&p.SoloOSEnv,
		nil, &p.DataNodeClient, p.MetaStg.PrepareNetBlockMetaData))

	util.AssertErrIsNil(p.NetINodeDriver.Init(&p.SoloOSEnv,
		&p.NetBlockDriver,
		&p.MemBlockDriver,
		nil,
		p.MetaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		p.MetaStg.PrepareNetINodeMetaDataWithStorDB,
		p.MetaStg.NetINodeCommitSizeInDB,
	))

	util.AssertErrIsNil(p.dataNode.Init(&p.SoloOSEnv, dataNodeOptions,
		&p.MemBlockDriver,
		&p.NetBlockDriver,
		&p.NetINodeDriver,
	))
	util.AssertErrIsNil(p.dataNode.Serve())
	util.AssertErrIsNil(p.dataNode.Close())
}

func (p *Env) Start() {
	go func() {
		util.AssertErrIsNil(p.cronSoloBoatJob())
	}()

	if p.options.Mode == "namenode" {
		p.startCommon()
		p.startNameNode()
	}

	if p.options.Mode == "datanode" {
		p.startCommon()
		p.startDataNode()
	}
}
