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
		soloOSEnv          soloosbase.SoloOSEnv
		netDriverSoloOSEnv soloosbase.SoloOSEnv
	)

	var (
		solonnSRPCPeerID          snettypes.PeerID = snet.MakeSysPeerID("SolonnSRPCForTest")
		solonnSRPCListenAddr                       = "127.0.0.1:10300"
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

	util.AssertErrIsNil(netDriverSoloOSEnv.InitWithSNet(""))
	util.AssertErrIsNil(netDriverSoloOSEnv.SNetDriver.Init(&netDriverSoloOSEnv.OffheapDriver))
	{
		var webServerOptions = iron.Options{
			ListenStr: netDriverServerListenAddr,
			ServeStr:  netDriverServerServeAddr,
		}
		util.AssertErrIsNil(netDriverWebServer.Init(webServerOptions))
	}
	go func() {
		util.AssertErrIsNil(netDriverSoloOSEnv.SNetDriver.PrepareServer("",
			&netDriverWebServer,
			nil, nil))
		util.AssertErrIsNil(netDriverSoloOSEnv.SNetDriver.ServerServe())
	}()

	// wait netDriverSoloOSEnv SNetDriver ServerServe
	time.Sleep(time.Millisecond * 200)

	util.AssertErrIsNil(soloOSEnv.InitWithSNet(netDriverServerServeAddr))

	memstg.MemStgMakeDriversForTest(&soloOSEnv,
		solonnSRPCListenAddr,
		memBlockDriverForClient, netBlockDriverForClient, netINodeDriverForClient, memBlockCap, blocksLimit)

	memstg.MemStgMakeDriversForTest(&soloOSEnv,
		solonnSRPCListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer, memBlockCap, blocksLimit)

	metastg.MakeMetaStgForTest(&soloOSEnv, &metaStg)
	solonn.MakeSolonnForTest(&soloOSEnv, &solonnIns, &metaStg,
		solonnSRPCPeerID, solonnSRPCListenAddr,
		solonnWebPeerID, solonnWebListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer)

	go func() {
		util.AssertErrIsNil(solonnIns.Serve())
	}()

	time.Sleep(time.Millisecond * 600)

	memstg.MakeMockServerForTest(&soloOSEnv, mockServerAddr, &mockServer)
	mockMemBlockTable.Init(&soloOSEnv, 1024)

	for i = 0; i < 6; i++ {
		snettypes.InitTmpPeerID((*snettypes.PeerID)(&peer.ID))
		peer.SetAddress(mockServerAddr)
		peer.ServiceProtocol = solofsapitypes.DefaultSOLOFSRPCProtocol
		solonnIns.SolodnRegister(peer)
	}

	var (
		dbConn solodbapi.Connection
		err    error
	)
	err = dbConn.Init(metastg.TestMetaStgDBDriver, metastg.TestMetaStgDBConnect)
	util.AssertErrIsNil(err)
	util.AssertErrIsNil(client.Init(&soloOSEnv, solofsapitypes.DefaultNameSpaceID,
		&memStg, &dbConn, netBlockCap, memBlockCap))
}
