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

	uNetINode, err = netINodeDriver.AllocNetINode(0, netBlockCap, memBlockCap)
	assert.NoError(t, err)

	var (
		readData       = make([]byte, 93)
		readOff  int64 = 73
	)

	uNetBlock, err := netBlockDriver.MustGetBlock(uNetINode, 10)
	uMemBlock := mockMemBlockPool.AllocMemBlock()
	memBlockIndex := 0
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, uMemBlock, memBlockIndex, 0, 12))
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, uMemBlock, memBlockIndex, 11, 24))
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, uMemBlock, memBlockIndex, 30, 64))
	assert.NoError(t, netBlockDriver.FlushMemBlock(uMemBlock))
	assert.NoError(t, netINodeDriver.PRead(uNetINode, readData, readOff))
	util.Ignore(uNetINode)

	assert.NoError(t, nameNode.Close())
	assert.NoError(t, mockServer.Close())
}
