package main

import (
	"soloos/common/snet"
	snettypes "soloos/common/snet/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/common/util"
	"soloos/sdbone/offheap"
	"soloos/sdfs/api"
	"soloos/sdfs/datanode"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/namenode"
	"soloos/sdfs/netstg"
)

type Env struct {
	options          Options
	SoloOSEnv        soloosbase.SoloOSEnv
	offheapDriver    *offheap.OffheapDriver
	SNetDriver       snet.NetDriver
	SNetClientDriver snet.ClientDriver
	MetaStg          metastg.MetaStg
	DataNodeClient   api.DataNodeClient
	MemBlockDriver   memstg.MemBlockDriver
	NetBlockDriver   netstg.NetBlockDriver
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
	if p.options.PProfListenAddr != "" {
		go util.PProfServe(p.options.PProfListenAddr)
	}
}

func (p *Env) startNameNode() {
	var (
		nameNode       namenode.NameNode
		nameNodePeerID snettypes.PeerID
	)
	copy(nameNodePeerID[:], []byte(p.options.NameNodePeerIDStr))

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
		p.options.ListenAddr,
		nameNodePeerID,
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

	copy(dataNodePeerID[:], []byte(p.options.DataNodePeerIDStr))
	copy(nameNodePeerID[:], []byte(p.options.NameNodePeerIDStr))

	dataNodeOptions = datanode.DataNodeOptions{
		PeerID:               dataNodePeerID,
		SrpcServerListenAddr: p.options.ListenAddr,
		SrpcServerServeAddr:  p.options.ListenAddr,
		LocalFsRoot:          p.options.DataNodeLocalFsRoot,
		NameNodePeerID:       nameNodePeerID,
		NameNodeSRPCServer:   p.options.NameNodeAddr,
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
		&p.MetaStg,
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
