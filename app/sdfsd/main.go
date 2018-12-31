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
	offheapDriver    *offheap.OffheapDriver
	SNetDriver       snet.NetDriver
	SNetClientDriver snet.ClientDriver
	MetaStg          metastg.MetaStg
	DataNodeClient   api.DataNodeClient
	MemBlockDriver   memstg.MemBlockDriver
	NetBlockDriver   netstg.NetBlockDriver
	NetINodeDriver   memstg.NetINodeDriver
}

func (p *Env) Init() {
	var (
		dbDriver = "sqlite3"
		dsn      = "/tmp/sdfs.db"
	)

	p.offheapDriver = &offheap.DefaultOffheapDriver

	util.AssertErrIsNil(p.SNetDriver.Init(p.offheapDriver))
	util.AssertErrIsNil(p.SNetClientDriver.Init(p.offheapDriver))

	util.AssertErrIsNil(p.MetaStg.Init(p.offheapDriver,
		dbDriver, dsn, nil))

	p.DataNodeClient.Init(&p.SNetClientDriver)

	{
		var options = memstg.MemBlockDriverOptions{
			[]memstg.MemBlockPoolOptions{
				memstg.MemBlockPoolOptions{
					1024 * 1024 * 2,
					256,
				},
			},
		}
		util.AssertErrIsNil(p.MemBlockDriver.Init(p.offheapDriver, options))
	}

	util.AssertErrIsNil(p.NetBlockDriver.Init(p.offheapDriver,
		&p.SNetDriver, &p.SNetClientDriver,
		nil, &p.DataNodeClient, nil))

	util.AssertErrIsNil(p.NetINodeDriver.Init(p.offheapDriver,
		&p.NetBlockDriver,
		&p.MemBlockDriver,
		nil, nil, nil,
	))
}

func (p *Env) startNameNode() {
	var (
		listenAddr = os.Args[2]
		nameNode   namenode.NameNode
	)

	p.NetBlockDriver.SetHelper(nil, p.MetaStg.PrepareNetBlockMetaData)
	p.NetINodeDriver.SetHelper(nil,
		p.MetaStg.PrepareNetINodeMetaDataOnlyLoadDB, p.MetaStg.PrepareNetINodeMetaDataWithStorDB)
	util.AssertErrIsNil(nameNode.Init(p.offheapDriver,
		listenAddr,
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
		dataNodePeerIDStr   = os.Args[2]
		listenAddr          = os.Args[3]
		dataNodeLocalFsRoot = os.Args[4]
		nameNodePeerIDStr   = os.Args[5]
		nameNodeAddr        = os.Args[6]
		dataNodePeerID      snettypes.PeerID
		dataNode            datanode.DataNode
		nameNodePeerID      snettypes.PeerID
		options             datanode.DataNodeOptions
	)

	copy(dataNodePeerID[:], []byte(dataNodePeerIDStr))
	copy(nameNodePeerID[:], []byte(nameNodePeerIDStr))

	options = datanode.DataNodeOptions{
		PeerID:               dataNodePeerID,
		SrpcServerListenAddr: listenAddr,
		SrpcServerServeAddr:  listenAddr,
		LocalFsRoot:          dataNodeLocalFsRoot,
		NameNodePeerID:       nameNodePeerID,
		NameNodeSRPCServer:   nameNodeAddr,
	}

	p.NetBlockDriver.SetHelper(nil, p.MetaStg.PrepareNetBlockMetaData)
	p.NetINodeDriver.SetHelper(nil,
		p.MetaStg.PrepareNetINodeMetaDataOnlyLoadDB, p.MetaStg.PrepareNetINodeMetaDataWithStorDB)
	util.AssertErrIsNil(dataNode.Init(p.offheapDriver, options,
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
	var env Env
	env.Init()

	mode := os.Args[1]

	if mode == "namenode" {
		env.startNameNode()
	}

	if mode == "datanode" {
		env.startDataNode()
	}
}
