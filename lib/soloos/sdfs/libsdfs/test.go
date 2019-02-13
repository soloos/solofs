package libsdfs

import (
	"soloos/common/sdbapi"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/namenode"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/common/snet"
	snettypes "soloos/common/snet/types"
	"soloos/common/util"
	"soloos/common/util/offheap"
	"time"
)

func MakeClientForTest(client *Client) {
	client.offheapDriver = &offheap.DefaultOffheapDriver
	var (
		memStg  memstg.MemStg
		metaStg metastg.MetaStg
	)

	var (
		offheapDriver          = &offheap.DefaultOffheapDriver
		nameNodeSRPCListenAddr = "127.0.0.1:10300"
		nameNode               namenode.NameNode
		mockServerAddr         = "127.0.0.1:10301"
		mockServer             netstg.MockServer
		mockMemBlockPool       types.MockMemBlockPool
		snetDriver             snet.NetDriver
		snetClientDriver       snet.ClientDriver

		memBlockDriverClient *memstg.MemBlockDriver = &memStg.MemBlockDriver
		netBlockDriverClient *netstg.NetBlockDriver = &memStg.NetBlockDriver
		netINodeDriverClient *memstg.NetINodeDriver = &memStg.NetINodeDriver

		memBlockDriverServer memstg.MemBlockDriver
		netBlockDriverServer netstg.NetBlockDriver
		netINodeDriverServer memstg.NetINodeDriver

		netBlockCap      int   = 1280
		memBlockCap      int   = 128
		blockChunksLimit int32 = 4
		peerID           snettypes.PeerID
		i                int
	)

	memstg.MakeDriversForTest(&snetDriver, &snetClientDriver,
		nameNodeSRPCListenAddr,
		memBlockDriverClient, netBlockDriverClient, netINodeDriverClient, memBlockCap, blockChunksLimit)
	memstg.MakeDriversForTest(&snetDriver, &snetClientDriver,
		nameNodeSRPCListenAddr,
		&memBlockDriverServer, &netBlockDriverServer, &netINodeDriverServer, memBlockCap, blockChunksLimit)
	metastg.MakeMetaStgForTest(offheapDriver, &metaStg)
	namenode.MakeNameNodeForTest(&nameNode, &metaStg, nameNodeSRPCListenAddr,
		&memBlockDriverServer, &netBlockDriverServer, &netINodeDriverServer)
	go func() {
		util.AssertErrIsNil(nameNode.Serve())
	}()
	time.Sleep(time.Millisecond * 300)
	netstg.MakeMockServerForTest(&snetDriver, mockServerAddr, &mockServer)
	mockMemBlockPool.Init(offheapDriver, 1024)

	for i = 0; i < 6; i++ {
		util.InitUUID64(&peerID)
		nameNode.RegisterDataNode(&peerID, mockServerAddr)
	}

	var (
		dbConn sdbapi.Connection
		err    error
	)
	err = dbConn.Init(metastg.TestMetaStgDBDriver, metastg.TestMetaStgDBConnect)
	util.AssertErrIsNil(err)
	util.AssertErrIsNil(client.Init(&memStg, &dbConn, netBlockCap, memBlockCap))
}
