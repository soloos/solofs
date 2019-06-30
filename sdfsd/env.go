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
	"soloos/silicon/siliconsdk"
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

	nameNode namenode.NameNode
	dataNode datanode.DataNode
	peerID   snettypes.PeerID

	siliconClient          siliconsdk.Client
	siliconCronJobDuration time.Duration
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
	util.AssertErrIsNil(p.SoloOSEnv.Init())

	p.DataNodeClient.Init(&p.SoloOSEnv)

	util.AssertErrIsNil(p.initMetaStg())

	util.AssertErrIsNil(p.initMemStg())

	util.AssertErrIsNil(p.initSilicon())
}

func (p *Env) startCommon() {
	util.AssertErrIsNil(p.SoloOSEnv.SNetDriver.StartClient(p.options.SNetDriverServeAddr))

	if p.options.PProfListenAddr != "" {
		go util.PProfServe(p.options.PProfListenAddr)
	}
}

func (p *Env) startNameNode() {
	copy(p.peerID[:], []byte(p.options.NameNodePeerID))

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
		p.peerID,
		p.options.ListenAddr, p.options.ServeAddr,
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
		nameNodePeerID  snettypes.PeerID
		dataNodeOptions datanode.DataNodeOptions
	)

	copy(p.peerID[:], []byte(p.options.DataNodePeerID))
	copy(nameNodePeerID[:], []byte(p.options.NameNodePeerID))

	dataNodeOptions = datanode.DataNodeOptions{
		PeerID:               p.peerID,
		SrpcServerListenAddr: p.options.ListenAddr,
		SrpcServerServeAddr:  p.options.ServeAddr,
		LocalFSRoot:          p.options.DataNodeLocalFSRoot,
		NameNodePeerID:       nameNodePeerID,
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
		util.AssertErrIsNil(p.cronSiliconJob())
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
