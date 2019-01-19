package namenode

import (
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	var (
		offheapDriver          = &offheap.DefaultOffheapDriver
		metaStg                metastg.MetaStg
		nameNodeSRPCListenAddr = "127.0.0.1:10300"
		nameNode               NameNode
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

		netBlockCap      int   = 1024
		memBlockCap      int   = 128
		blockChunksLimit int32 = 4
		uNetINode        types.NetINodeUintptr
		peerID           snettypes.PeerID
		i                int
		err              error
	)
	memstg.MakeDriversForTest(&snetDriver, &snetClientDriver,
		nameNodeSRPCListenAddr,
		&memBlockDriverClient, &netBlockDriverClient, &netINodeDriverClient, memBlockCap, blockChunksLimit)
	memstg.MakeDriversForTest(&snetDriver, &snetClientDriver,
		nameNodeSRPCListenAddr,
		&memBlockDriverServer, &netBlockDriverServer, &netINodeDriverServer, memBlockCap, blockChunksLimit)
	metastg.MakeMetaStgForTest(offheapDriver, &metaStg)
	MakeNameNodeForTest(&nameNode, &metaStg, nameNodeSRPCListenAddr,
		&memBlockDriverServer, &netBlockDriverServer, &netINodeDriverServer)
	go func() {
		assert.NoError(t, nameNode.Serve())
	}()
	time.Sleep(time.Millisecond * 300)
	netstg.MakeMockServerForTest(&snetDriver, mockServerAddr, &mockServer)
	mockMemBlockPool.Init(offheapDriver, 1024)

	for i = 0; i < 6; i++ {
		util.InitUUID64(&peerID)
		nameNode.RegisterDataNode(&peerID, mockServerAddr)
	}

	var netINodeID types.NetINodeID
	util.InitUUID64(&netINodeID)
	uNetINode, err = netINodeDriverClient.MustGetNetINodeWithReadAcquire(netINodeID, 0, netBlockCap, memBlockCap)
	assert.NoError(t, err)

	var (
		testData = make([]byte, 93)
	)

	assert.NoError(t, netINodeDriverClient.PWriteWithMem(uNetINode, testData[0:12], 0))
	assert.NoError(t, netINodeDriverClient.PWriteWithMem(uNetINode, testData[11:24], 24))
	assert.NoError(t, netINodeDriverClient.PWriteWithMem(uNetINode, testData[30:64], 64))
	assert.NoError(t, netINodeDriverClient.Flush(uNetINode))
	_, err = netINodeDriverClient.PReadWithMem(uNetINode, testData, 73)
	assert.NoError(t, err)

	assert.NoError(t, nameNode.Close())
	assert.NoError(t, mockServer.Close())
}
