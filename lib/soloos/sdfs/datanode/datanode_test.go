package datanode

import (
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	var (
		offheapDriver          = &offheap.DefaultOffheapDriver
		metaStg                metastg.MetaStg
		dataNodeSRPCListenAddr = "127.0.0.1:10400"
		dataNode               DataNode
		uDataNode              snettypes.PeerUintptr
		mockServerAddr         = "127.0.0.1:10401"
		mockServer             netstg.MockServer
	)

	var (
		mockMemBlockPool types.MockMemBlockPool
		snetDriver       snet.SNetDriver
		memBlockDriver   memstg.MemBlockDriver
		netBlockDriver   netstg.NetBlockDriver
		netINodeDriver   memstg.NetINodeDriver
		netBlockCap      int   = 32
		memBlockCap      int   = 16
		blockChunksLimit int32 = 4
		uNetINode        types.NetINodeUintptr
		err              error
	)
	memstg.MakeDriversWithMockServerForTest(t,
		mockServerAddr, &mockServer, &snetDriver,
		&netBlockDriver, &memBlockDriver, &netINodeDriver, memBlockCap, blockChunksLimit)
	mockMemBlockPool.Init(offheapDriver, 1024)

	uDataNode, _ = snetDriver.MustGetPeer(nil, dataNodeSRPCListenAddr, types.DefaultSDFSRPCProtocol)
	mockServer.SetDataNodePeers([]snettypes.PeerUintptr{uDataNode, uDataNode})

	metastg.MakeMetaStgForTest(offheapDriver, &metaStg)
	MakeDataNodeForTest(&dataNode, &metaStg, &netBlockDriver, &memBlockDriver, &netINodeDriver, dataNodeSRPCListenAddr)
	go func() {
		assert.NoError(t, dataNode.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	uNetINode, err = netINodeDriver.AllocNetINode(0, netBlockCap, memBlockCap)
	assert.NoError(t, err)

	writeData := make([]byte, 6)
	assert.NoError(t, netINodeDriver.PWrite(uNetINode, writeData, 12))
	assert.NoError(t, netINodeDriver.Flush(uNetINode))

	assert.NoError(t, dataNode.Close())
	assert.NoError(t, mockServer.Close())
	// assert.Equal(t, true, false)
}
