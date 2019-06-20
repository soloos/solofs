package sdfscli

import (
	"soloos/common/sdbapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/snet"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/namenode"
	"soloos/sdfs/types"
	"time"
)

func MakeClientForTest(client *Client) {
	var (
		memStg             memstg.MemStg
		metaStg            metastg.MetaStg
		soloOSEnv          soloosbase.SoloOSEnv
		netDriverSoloOSEnv soloosbase.SoloOSEnv
	)

	util.AssertErrIsNil(soloOSEnv.Init())

	var (
		nameNodePeerID            snettypes.PeerID = snet.MakeSysPeerID("NameNodeForTest")
		nameNodeSRPCListenAddr                     = "127.0.0.1:10300"
		netDriverServerListenAddr                  = "127.0.0.1:10402"
		netDriverServerServeAddr                   = "http://127.0.0.1:10402"
		nameNode                  namenode.NameNode
		mockServerAddr            = "127.0.0.1:10301"
		mockServer                memstg.MockServer
		mockMemBlockTable         types.MockMemBlockTable

		memBlockDriverForClient *memstg.MemBlockDriver = &memStg.MemBlockDriver
		netBlockDriverForClient *memstg.NetBlockDriver = &memStg.NetBlockDriver
		netINodeDriverForClient *memstg.NetINodeDriver = &memStg.NetINodeDriver

		memBlockDriverForServer memstg.MemBlockDriver
		netBlockDriverForServer memstg.NetBlockDriver
		netINodeDriverForServer memstg.NetINodeDriver

		netBlockCap int   = 1280
		memBlockCap int   = 128
		blocksLimit int32 = 4
		peer        snettypes.Peer
		i           int
	)

	memstg.MemStgMakeDriversForTest(&soloOSEnv,
		nameNodeSRPCListenAddr,
		memBlockDriverForClient, netBlockDriverForClient, netINodeDriverForClient, memBlockCap, blocksLimit)

	memstg.MemStgMakeDriversForTest(&soloOSEnv,
		nameNodeSRPCListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer, memBlockCap, blocksLimit)

	metastg.MakeMetaStgForTest(&soloOSEnv, &metaStg)
	namenode.MakeNameNodeForTest(&soloOSEnv, &nameNode, &metaStg,
		nameNodePeerID, nameNodeSRPCListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer)

	util.AssertErrIsNil(netDriverSoloOSEnv.Init())
	util.AssertErrIsNil(netDriverSoloOSEnv.SNetDriver.Init(&netDriverSoloOSEnv.OffheapDriver))
	go func() {
		util.AssertErrIsNil(netDriverSoloOSEnv.SNetDriver.StartServer(netDriverServerListenAddr,
			netDriverServerServeAddr,
			nil, nil))
	}()

	util.AssertErrIsNil(soloOSEnv.SNetDriver.StartClient(netDriverServerServeAddr))

	go func() {
		util.AssertErrIsNil(nameNode.Serve())
	}()

	time.Sleep(time.Millisecond * 600)

	memstg.MakeMockServerForTest(&soloOSEnv, mockServerAddr, &mockServer)
	mockMemBlockTable.Init(&soloOSEnv, 1024)

	for i = 0; i < 6; i++ {
		snettypes.InitTmpPeerID((*snettypes.PeerID)(&peer.ID))
		peer.SetAddress(mockServerAddr)
		peer.ServiceProtocol = sdfsapitypes.DefaultSDFSRPCProtocol
		nameNode.RegisterDataNode(peer)
	}

	var (
		dbConn sdbapi.Connection
		err    error
	)
	err = dbConn.Init(metastg.TestMetaStgDBDriver, metastg.TestMetaStgDBConnect)
	util.AssertErrIsNil(err)
	util.AssertErrIsNil(client.Init(&soloOSEnv, sdfsapitypes.DefaultNameSpaceID,
		&memStg, &dbConn, netBlockCap, memBlockCap))
}
