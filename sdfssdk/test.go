package sdfssdk

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
	"soloos/sdfs/sdfstypes"
	"time"
)

func MakeClientForTest(client *Client) {
	var (
		memStg             memstg.MemStg
		metaStg            metastg.MetaStg
		soloOSEnv          soloosbase.SoloOSEnv
		netDriverSoloOSEnv soloosbase.SoloOSEnv
	)

	var (
		nameNodeSRPCPeerID        snettypes.PeerID = snet.MakeSysPeerID("NameNodeSRPCForTest")
		nameNodeSRPCListenAddr                     = "127.0.0.1:10300"
		nameNodeWebPeerID         snettypes.PeerID = snet.MakeSysPeerID("NameNodeWebForTest")
		nameNodeWebListenAddr                      = "127.0.0.1:10301"
		netDriverServerListenAddr                  = "127.0.0.1:10402"
		netDriverServerServeAddr                   = "http://127.0.0.1:10402"
		nameNode                  namenode.NameNode
		mockServerAddr            = "127.0.0.1:10302"
		mockServer                memstg.MockServer
		mockMemBlockTable         sdfstypes.MockMemBlockTable

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

	util.AssertErrIsNil(netDriverSoloOSEnv.InitWithSNet(""))
	util.AssertErrIsNil(netDriverSoloOSEnv.SNetDriver.Init(&netDriverSoloOSEnv.OffheapDriver))
	go func() {
		util.AssertErrIsNil(netDriverSoloOSEnv.SNetDriver.PrepareServer(netDriverServerListenAddr,
			netDriverServerServeAddr,
			nil, nil))
		util.AssertErrIsNil(netDriverSoloOSEnv.SNetDriver.ServerServe())
	}()

	// wait netDriverSoloOSEnv SNetDriver ServerServe
	time.Sleep(time.Millisecond * 200)

	util.AssertErrIsNil(soloOSEnv.InitWithSNet(netDriverServerServeAddr))

	memstg.MemStgMakeDriversForTest(&soloOSEnv,
		nameNodeSRPCListenAddr,
		memBlockDriverForClient, netBlockDriverForClient, netINodeDriverForClient, memBlockCap, blocksLimit)

	memstg.MemStgMakeDriversForTest(&soloOSEnv,
		nameNodeSRPCListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer, memBlockCap, blocksLimit)

	metastg.MakeMetaStgForTest(&soloOSEnv, &metaStg)
	namenode.MakeNameNodeForTest(&soloOSEnv, &nameNode, &metaStg,
		nameNodeSRPCPeerID, nameNodeSRPCListenAddr,
		nameNodeWebPeerID, nameNodeWebListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer)

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
		nameNode.DataNodeRegister(peer)
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
