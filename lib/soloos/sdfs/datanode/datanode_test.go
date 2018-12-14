package datanode

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
		netBlockCap      int   = 1024
		memBlockCap      int   = 128
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

	var (
		readData            = make([]byte, 93)
		readOff       int64 = 73
		uNetBlock     types.NetBlockUintptr
		uMemBlock     types.MemBlockUintptr
		memBlockIndex int
	)

	uNetBlock, err = netBlockDriver.MustGetBlock(uNetINode, 10)
	uMemBlock = mockMemBlockPool.AllocMemBlock()
	memBlockIndex = 0
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, uMemBlock, memBlockIndex, 0, 12))
	// assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, uMemBlock, memBlockIndex, 11, 24))
	// assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, uMemBlock, memBlockIndex, 30, 64))
	assert.NoError(t, netBlockDriver.FlushMemBlock(uMemBlock))
	// assert.NoError(t, netINodeDriver.PRead(uNetINode, readData, readOff))
	util.Ignore(readOff)
	util.Ignore(readData)
	util.Ignore(memBlockIndex)
	util.Ignore(uNetBlock)
	util.Ignore(uMemBlock)
	util.Ignore(uNetINode)

	assert.NoError(t, dataNode.Close())
	assert.NoError(t, mockServer.Close())
}
