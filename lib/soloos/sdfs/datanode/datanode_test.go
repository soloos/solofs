package datanode

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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	var (
		offheapDriver          = &offheap.DefaultOffheapDriver
		metaStg                metastg.MetaStg
		nameNodeSRPCListenAddr = "127.0.0.1:10401"
		nameNode               namenode.NameNode
		dataNodeSRPCListenAddr = "127.0.0.1:10400"
		dataNode               DataNode
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
		i                int
		err              error
	)
	metastg.MakeMetaStgForTest(offheapDriver, &metaStg)
	namenode.MakeNameNodeForTest(&nameNode, &metaStg, nameNodeSRPCListenAddr)
	memstg.MakeDriversForTest(t,
		&snetDriver,
		nameNodeSRPCListenAddr,
		&memBlockDriver, &netBlockDriver, &netINodeDriver, memBlockCap, blockChunksLimit)
	mockMemBlockPool.Init(offheapDriver, 1024)
	go func() {
		assert.NoError(t, nameNode.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	for i = 0; i < 6; i++ {
		var peerID snettypes.PeerID
		util.InitUUID64(&peerID)
		nameNode.RegisterDataNode(&peerID, dataNodeSRPCListenAddr)
	}

	MakeDataNodeForTest(&snetDriver,
		&dataNode, &metaStg, &netBlockDriver, &memBlockDriver, &netINodeDriver, dataNodeSRPCListenAddr)
	go func() {
		assert.NoError(t, dataNode.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	var netINodeID types.NetINodeID
	util.InitUUID64(&netINodeID)
	uNetINode, err = netINodeDriver.MustGetNetINode(netINodeID, 0, netBlockCap, memBlockCap)
	assert.NoError(t, err)

	writeData := make([]byte, 6)
	assert.NoError(t, netINodeDriver.PWrite(uNetINode, writeData, 12))
	assert.NoError(t, netINodeDriver.Flush(uNetINode))

	assert.NoError(t, dataNode.Close())
	assert.NoError(t, nameNode.Close())
	// time.Sleep(time.Second * 1)
	// assert.Equal(t, true, false)
}
