package libsdfs

import (
	"soloos/common/sdbapi"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/namenode"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"time"
)

func MakeClientForTest(client *Client) {
	var (
		memStg    memstg.MemStg
		metaStg   metastg.MetaStg
		soloOSEnv soloosbase.SoloOSEnv
	)

	util.AssertErrIsNil(soloOSEnv.Init())

	var (
		nameNodeSRPCListenAddr = "127.0.0.1:10300"
		nameNode               namenode.NameNode
		mockServerAddr         = "127.0.0.1:10301"
		mockServer             netstg.MockServer
		mockMemBlockTable      types.MockMemBlockTable

		memBlockDriverForClient *memstg.MemBlockDriver = &memStg.MemBlockDriver
		netBlockDriverForClient *netstg.NetBlockDriver = &memStg.NetBlockDriver
		netINodeDriverForClient *memstg.NetINodeDriver = &memStg.NetINodeDriver

		memBlockDriverForServer memstg.MemBlockDriver
		netBlockDriverForServer netstg.NetBlockDriver
		netINodeDriverForServer memstg.NetINodeDriver

		netBlockCap int   = 1280
		memBlockCap int   = 128
		blocksLimit int32 = 4
		peerID      snettypes.PeerID
		i           int
	)

	memstg.MakeDriversForTest(&soloOSEnv,
		nameNodeSRPCListenAddr,
		memBlockDriverForClient, netBlockDriverForClient, netINodeDriverForClient, memBlockCap, blocksLimit)

	memstg.MakeDriversForTest(&soloOSEnv,
		nameNodeSRPCListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer, memBlockCap, blocksLimit)

	metastg.MakeMetaStgForTest(&soloOSEnv, &metaStg)
	namenode.MakeNameNodeForTest(&soloOSEnv, &nameNode, &metaStg, nameNodeSRPCListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer)

	go func() {
		util.AssertErrIsNil(nameNode.Serve())
	}()
	time.Sleep(time.Millisecond * 300)
	netstg.MakeMockServerForTest(&soloOSEnv, mockServerAddr, &mockServer)
	mockMemBlockTable.Init(&soloOSEnv, 1024)

	for i = 0; i < 6; i++ {
		snettypes.InitTmpPeerID(&peerID)
		nameNode.RegisterDataNode(peerID, mockServerAddr)
	}

	var (
		dbConn sdbapi.Connection
		err    error
	)
	err = dbConn.Init(metastg.TestMetaStgDBDriver, metastg.TestMetaStgDBConnect)
	util.AssertErrIsNil(err)
	util.AssertErrIsNil(client.Init(&soloOSEnv, &memStg, &dbConn, netBlockCap, memBlockCap))
}
