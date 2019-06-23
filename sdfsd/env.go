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
}

func (p *Env) Init(options Options) {
	p.options = options
	util.AssertErrIsNil(p.SoloOSEnv.Init())

	util.AssertErrIsNil(p.MetaStg.Init(&p.SoloOSEnv,
		options.DBDriver, options.Dsn))

	p.DataNodeClient.Init(&p.SoloOSEnv)

	{
		var memBlockDriverOptions = memstg.MemBlockDriverOptions{
			[]memstg.MemBlockTableOptions{
				memstg.MemBlockTableOptions{
					p.options.DefaultMemBlockCap,
					p.options.DefaultMemBlocksLimit,
				},
			},
		}
		util.AssertErrIsNil(p.MemBlockDriver.Init(&p.SoloOSEnv, memBlockDriverOptions))
	}
}

func (p *Env) startCommon() {
	util.AssertErrIsNil(p.SoloOSEnv.SNetDriver.StartClient(p.options.SNetDriverServeAddr))

	if p.options.PProfListenAddr != "" {
		go util.PProfServe(p.options.PProfListenAddr)
	}
}

func (p *Env) startNameNode() {
	var (
		nameNode       namenode.NameNode
		nameNodePeerID snettypes.PeerID
	)
	copy(nameNodePeerID[:], []byte(p.options.NameNodePeerID))

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

	util.AssertErrIsNil(nameNode.Init(&p.SoloOSEnv,
		nameNodePeerID,
		p.options.ListenAddr, p.options.ServeAddr,
		&p.MetaStg,
		&p.MemBlockDriver,
		&p.NetBlockDriver,
		&p.NetINodeDriver,
	))

	util.AssertErrIsNil(nameNode.Serve())
	util.AssertErrIsNil(nameNode.Close())
}

func (p *Env) startDataNode() {
	var (
		dataNodePeerID  snettypes.PeerID
		dataNode        datanode.DataNode
		nameNodePeerID  snettypes.PeerID
		dataNodeOptions datanode.DataNodeOptions
	)

	copy(dataNodePeerID[:], []byte(p.options.DataNodePeerID))
	copy(nameNodePeerID[:], []byte(p.options.NameNodePeerID))

	dataNodeOptions = datanode.DataNodeOptions{
		PeerID:               dataNodePeerID,
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

	util.AssertErrIsNil(dataNode.Init(&p.SoloOSEnv, dataNodeOptions,
		&p.MemBlockDriver,
		&p.NetBlockDriver,
		&p.NetINodeDriver,
	))
	util.AssertErrIsNil(dataNode.Serve())
	util.AssertErrIsNil(dataNode.Close())
}

func (p *Env) Start() {
	if p.options.Mode == "namenode" {
		p.startCommon()
		p.startNameNode()
	}

	if p.options.Mode == "datanode" {
		p.startCommon()
		p.startDataNode()
	}
}
