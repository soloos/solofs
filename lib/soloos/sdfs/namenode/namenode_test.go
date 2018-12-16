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
	)
	metastg.MakeMetaStgForTest(offheapDriver, &metaStg)
	MakeNameNodeForTest(&nameNode, &metaStg, nameNodeSRPCListenAddr)
	go func() {
		assert.NoError(t, nameNode.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	var (
		mockMemBlockPool types.MockMemBlockPool
		snetDriver       snet.SNetDriver
		memBlockDriver   memstg.MemBlockDriver
		netBlockDriver   netstg.NetBlockDriver
		netINodeDriver   memstg.NetINodeDriver
		netBlockCap      int   = 1024
		memBlockCap      int   = 128
		blockChunksLimit int32 = 4
		uNetINode        types.NetINodeUintptr
		peerID           snettypes.PeerID
		i                int
		err              error
	)
	memstg.MakeDriversForTest(t,
		&snetDriver,
		nameNodeSRPCListenAddr,
		&memBlockDriver, &netBlockDriver, &netINodeDriver, memBlockCap, blockChunksLimit)
	netstg.MakeMockServerForTest(t, &snetDriver, mockServerAddr, &mockServer)
	mockMemBlockPool.Init(offheapDriver, 1024)

	for i = 0; i < 6; i++ {
		util.InitUUID64(&peerID)
		nameNode.metaStg.RegisterDataNode(&peerID, mockServerAddr)
	}

	var netINodeID types.NetINodeID
	util.InitUUID64(&netINodeID)
	uNetINode, err = netINodeDriver.MustGetNetINode(netINodeID, 0, netBlockCap, memBlockCap)
	assert.NoError(t, err)

	var (
		testData = make([]byte, 93)
	)

	assert.NoError(t, netINodeDriver.PWrite(uNetINode, testData[0:12], 0))
	assert.NoError(t, netINodeDriver.PWrite(uNetINode, testData[11:24], 24))
	assert.NoError(t, netINodeDriver.PWrite(uNetINode, testData[30:64], 64))
	assert.NoError(t, netINodeDriver.Flush(uNetINode))
	assert.NoError(t, netINodeDriver.PRead(uNetINode, testData, 73))

	assert.NoError(t, nameNode.Close())
	assert.NoError(t, mockServer.Close())
}
