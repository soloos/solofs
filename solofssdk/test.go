package solofssdk

import (
	"soloos/common/iron"
	"soloos/common/snet"
	"soloos/common/snettypes"
	"soloos/common/solodbapi"
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/solofs/memstg"
	"soloos/solofs/metastg"
	"soloos/solofs/solofstypes"
	"soloos/solofs/solonn"
	"time"
)

func MakeClientForTest(client *Client) {
	var (
		memStg             memstg.MemStg
		metaStg            metastg.MetaStg
		soloosEnv          soloosbase.SoloosEnv
		netDriverSoloosEnv soloosbase.SoloosEnv
	)

	var (
		solonnSrpcPeerID          snettypes.PeerID = snet.MakeSysPeerID("SolonnSrpcForTest")
		solonnSrpcListenAddr                       = "127.0.0.1:10300"
		solonnWebPeerID           snettypes.PeerID = snet.MakeSysPeerID("SolonnWebForTest")
		solonnWebListenAddr                        = "127.0.0.1:10301"
		netDriverWebServer        iron.Server
		netDriverServerListenAddr = "127.0.0.1:10402"
		netDriverServerServeAddr  = "http://127.0.0.1:10402"
		solonnIns                 solonn.Solonn
		mockServerAddr            = "127.0.0.1:10302"
		mockServer                memstg.MockServer
		mockMemBlockTable         solofstypes.MockMemBlockTable

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

	util.AssertErrIsNil(netDriverSoloosEnv.InitWithSNet(""))
	util.AssertErrIsNil(netDriverSoloosEnv.SNetDriver.Init(&netDriverSoloosEnv.OffheapDriver))
	{
		var webServerOptions = iron.Options{
			ListenStr: netDriverServerListenAddr,
			ServeStr:  netDriverServerServeAddr,
		}
		util.AssertErrIsNil(netDriverWebServer.Init(webServerOptions))
	}
	go func() {
		util.AssertErrIsNil(netDriverSoloosEnv.SNetDriver.PrepareServer("",
			&netDriverWebServer,
			nil, nil))
		util.AssertErrIsNil(netDriverSoloosEnv.SNetDriver.ServerServe())
	}()

	// wait netDriverSoloosEnv SNetDriver ServerServe
	time.Sleep(time.Millisecond * 200)

	util.AssertErrIsNil(soloosEnv.InitWithSNet(netDriverServerServeAddr))

	memstg.MemStgMakeDriversForTest(&soloosEnv,
		solonnSrpcListenAddr,
		memBlockDriverForClient, netBlockDriverForClient, netINodeDriverForClient, memBlockCap, blocksLimit)

	memstg.MemStgMakeDriversForTest(&soloosEnv,
		solonnSrpcListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer, memBlockCap, blocksLimit)

	metastg.MakeMetaStgForTest(&soloosEnv, &metaStg)
	solonn.MakeSolonnForTest(&soloosEnv, &solonnIns, &metaStg,
		solonnSrpcPeerID, solonnSrpcListenAddr,
		solonnWebPeerID, solonnWebListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer)

	go func() {
		util.AssertErrIsNil(solonnIns.Serve())
	}()

	time.Sleep(time.Millisecond * 600)

	memstg.MakeMockServerForTest(&soloosEnv, mockServerAddr, &mockServer)
	mockMemBlockTable.Init(&soloosEnv, 1024)

	for i = 0; i < 6; i++ {
		snettypes.InitTmpPeerID((*snettypes.PeerID)(&peer.ID))
		peer.SetAddress(mockServerAddr)
		peer.ServiceProtocol = solofsapitypes.DefaultSolofsRPCProtocol
		solonnIns.SolodnRegister(peer)
	}

	var (
		dbConn solodbapi.Connection
		err    error
	)
	err = dbConn.Init(metastg.TestMetaStgDBDriver, metastg.TestMetaStgDBConnect)
	util.AssertErrIsNil(err)
	util.AssertErrIsNil(client.Init(&soloosEnv, solofsapitypes.DefaultNameSpaceID,
		&memStg, &dbConn, netBlockCap, memBlockCap))
}
