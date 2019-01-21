package main

import (
	"os"
	"soloos/sdfs/api"
	"soloos/sdfs/datanode"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/namenode"
	"soloos/sdfs/netstg"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"
)

type Env struct {
	options          Options
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
	p.offheapDriver = &offheap.DefaultOffheapDriver

	util.AssertErrIsNil(p.SNetDriver.Init(p.offheapDriver))
	util.AssertErrIsNil(p.SNetClientDriver.Init(p.offheapDriver))

	util.AssertErrIsNil(p.MetaStg.Init(p.offheapDriver,
		options.DBDriver, options.Dsn))

	p.DataNodeClient.Init(&p.SNetClientDriver)

	{
		var memBlockDriverOptions = memstg.MemBlockDriverOptions{
			[]memstg.MemBlockPoolOptions{
				memstg.MemBlockPoolOptions{
					p.options.MemBlockChunkSize,
					p.options.MemBlockChunksLimit,
				},
			},
		}
		util.AssertErrIsNil(p.MemBlockDriver.Init(p.offheapDriver, memBlockDriverOptions))
	}
}

func (p *Env) startCommon(options Options) {
	if options.PProfListenAddr != "" {
		go util.PProfServe(options.PProfListenAddr)
	}
}

func (p *Env) startNameNode(options Options) {
	var (
		nameNode namenode.NameNode
	)

	util.AssertErrIsNil(p.NetBlockDriver.Init(p.offheapDriver,
		&p.SNetDriver, &p.SNetClientDriver,
		nil, &p.DataNodeClient, p.MetaStg.PrepareNetBlockMetaData))

	util.AssertErrIsNil(p.NetINodeDriver.Init(p.offheapDriver,
		&p.NetBlockDriver,
		&p.MemBlockDriver,
		nil,
		p.MetaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		p.MetaStg.PrepareNetINodeMetaDataWithStorDB,
	))

	util.AssertErrIsNil(nameNode.Init(p.offheapDriver,
		options.ListenAddr,
		&p.MetaStg,
		&p.MemBlockDriver,
		&p.NetBlockDriver,
		&p.NetINodeDriver,
	))

	util.AssertErrIsNil(nameNode.Serve())
	util.AssertErrIsNil(nameNode.Close())
}

func (p *Env) startDataNode(options Options) {
	var (
		dataNodePeerID  snettypes.PeerID
		dataNode        datanode.DataNode
		nameNodePeerID  snettypes.PeerID
		dataNodeOptions datanode.DataNodeOptions
	)

	copy(dataNodePeerID[:], []byte(options.DataNodePeerIDStr))
	copy(nameNodePeerID[:], []byte(options.NameNodePeerIDStr))

	dataNodeOptions = datanode.DataNodeOptions{
		PeerID:               dataNodePeerID,
		SrpcServerListenAddr: options.ListenAddr,
		SrpcServerServeAddr:  options.ListenAddr,
		LocalFsRoot:          options.DataNodeLocalFsRoot,
		NameNodePeerID:       nameNodePeerID,
		NameNodeSRPCServer:   options.NameNodeAddr,
	}

	util.AssertErrIsNil(p.NetBlockDriver.Init(p.offheapDriver,
		&p.SNetDriver, &p.SNetClientDriver,
		nil, &p.DataNodeClient, p.MetaStg.PrepareNetBlockMetaData))

	util.AssertErrIsNil(p.NetINodeDriver.Init(p.offheapDriver,
		&p.NetBlockDriver,
		&p.MemBlockDriver,
		nil,
		p.MetaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		p.MetaStg.PrepareNetINodeMetaDataWithStorDB,
	))

	util.AssertErrIsNil(dataNode.Init(p.offheapDriver, dataNodeOptions,
		&p.SNetDriver, &p.SNetClientDriver,
		&p.MetaStg,
		&p.MemBlockDriver,
		&p.NetBlockDriver,
		&p.NetINodeDriver,
	))
	util.AssertErrIsNil(dataNode.Serve())
	util.AssertErrIsNil(dataNode.Close())
}

func main() {
	var (
		env     Env
		options Options
		err     error
	)

	optionsFile := os.Args[1]

	options, err = LoadOptionsFile(optionsFile)
	util.AssertErrIsNil(err)

	env.Init(options)

	if options.Mode == "namenode" {
		env.startCommon(options)
		env.startNameNode(options)
	}

	if options.Mode == "datanode" {
		env.startCommon(options)
		env.startDataNode(options)
	}
}
