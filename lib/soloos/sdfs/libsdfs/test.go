package libsdfs

import (
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/namenode"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"
	"time"
)

func MakeMetaStgForTest(metaStg *metastg.MetaStg) {
	var (
		offheapDriver          = &offheap.DefaultOffheapDriver
		nameNodeSRPCListenAddr = "127.0.0.1:10300"
		nameNode               namenode.NameNode
		mockServerAddr         = "127.0.0.1:10301"
		mockServer             netstg.MockServer
		mockMemBlockPool       types.MockMemBlockPool
		snetDriver             snet.NetDriver
		snetClientDriver       snet.ClientDriver

		memBlockDriverClient memstg.MemBlockDriver
		netBlockDriverClient netstg.NetBlockDriver
		netINodeDriverClient memstg.NetINodeDriver

		memBlockDriverServer memstg.MemBlockDriver
		netBlockDriverServer netstg.NetBlockDriver
		netINodeDriverServer memstg.NetINodeDriver

		memBlockCap      int   = 128
		blockChunksLimit int32 = 4
		peerID           snettypes.PeerID
		i                int
	)

	memstg.MakeDriversForTest(&snetDriver, &snetClientDriver,
		nameNodeSRPCListenAddr,
		&memBlockDriverClient, &netBlockDriverClient, &netINodeDriverClient, memBlockCap, blockChunksLimit)
	memstg.MakeDriversForTest(&snetDriver, &snetClientDriver,
		nameNodeSRPCListenAddr,
		&memBlockDriverServer, &netBlockDriverServer, &netINodeDriverServer, memBlockCap, blockChunksLimit)
	metastg.MakeMetaStgForTest(offheapDriver, metaStg)
	namenode.MakeNameNodeForTest(&nameNode, metaStg, nameNodeSRPCListenAddr,
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

	metaStg.DirTreeDriver.Init(offheapDriver,
		metaStg.GetDBConn(),
		metaStg.FetchAndUpdateMaxID,
		netINodeDriverClient.GetNetINode,
		netINodeDriverClient.MustGetNetINode,
	)
}

func MakeClientForTest(client *Client) {
	client.offheapDriver = &offheap.DefaultOffheapDriver
	MakeMetaStgForTest(&client.MetaStg)
}
