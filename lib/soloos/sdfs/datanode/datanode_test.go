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
		offheapDriver           = &offheap.DefaultOffheapDriver
		metaStg                 metastg.MetaStg
		nameNodeSRPCListenAddr  = "127.0.0.1:10401"
		nameNode                namenode.NameNode
		dataNodeSRPCListenAddrs = []string{
			"127.0.0.1:10410",
			"127.0.0.1:10411",
			"127.0.0.1:10412",
			"127.0.0.1:10413",
			"127.0.0.1:10414",
			"127.0.0.1:10415",
		}
		dataNodes [6]DataNode
	)

	var (
		mockMemBlockPool types.MockMemBlockPool
		snetDriver       snet.SNetDriver

		memBlockDriverClient memstg.MemBlockDriver
		netBlockDriverClient netstg.NetBlockDriver
		netINodeDriverClient memstg.NetINodeDriver

		memBlockDriverNameNode memstg.MemBlockDriver
		netBlockDriverNameNode netstg.NetBlockDriver
		netINodeDriverNameNode memstg.NetINodeDriver

		memBlockDriverDataNodes [6]memstg.MemBlockDriver
		netBlockDriverDataNodes [6]netstg.NetBlockDriver
		netINodeDriverDataNodes [6]memstg.NetINodeDriver

		netBlockCap      int   = 32
		memBlockCap      int   = 16
		blockChunksLimit int32 = 4
		uNetINode        types.NetINodeUintptr
		i                int
		err              error
	)
	metastg.MakeMetaStgForTest(offheapDriver, &metaStg)
	memstg.MakeDriversForTest(t,
		&snetDriver,
		nameNodeSRPCListenAddr,
		&memBlockDriverClient, &netBlockDriverClient, &netINodeDriverClient, memBlockCap, blockChunksLimit)

	memstg.MakeDriversForTest(t,
		&snetDriver,
		nameNodeSRPCListenAddr,
		&memBlockDriverNameNode, &netBlockDriverNameNode, &netINodeDriverNameNode, memBlockCap, blockChunksLimit)
	namenode.MakeNameNodeForTest(&nameNode, &metaStg, nameNodeSRPCListenAddr,
		&memBlockDriverNameNode, &netBlockDriverNameNode, &netINodeDriverNameNode)

	mockMemBlockPool.Init(offheapDriver, 1024)
	go func() {
		assert.NoError(t, nameNode.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	for i = 0; i < len(dataNodeSRPCListenAddrs); i++ {
		memstg.MakeDriversForTest(t,
			&snetDriver,
			nameNodeSRPCListenAddr,
			&memBlockDriverDataNodes[i],
			&netBlockDriverDataNodes[i],
			&netINodeDriverDataNodes[i],
			memBlockCap, blockChunksLimit)

		var peerID snettypes.PeerID
		util.InitUUID64(&peerID)
		nameNode.RegisterDataNode(&peerID, dataNodeSRPCListenAddrs[i])

		MakeDataNodeForTest(&snetDriver,
			&dataNodes[i], dataNodeSRPCListenAddrs[i], &metaStg,
			&memBlockDriverDataNodes[i],
			&netBlockDriverDataNodes[i],
			&netINodeDriverDataNodes[i])
		go func(localI int) {
			assert.NoError(t, dataNodes[localI].Serve())
		}(i)
	}
	time.Sleep(time.Millisecond * 300)

	var netINodeID types.NetINodeID
	util.InitUUID64(&netINodeID)
	uNetINode, err = netINodeDriverClient.MustGetNetINode(netINodeID, 0, netBlockCap, memBlockCap)
	assert.NoError(t, err)

	writeData := make([]byte, 73)
	writeData[3] = 12
	writeData[7] = 12
	writeData[8] = 12
	writeData[33] = 12
	writeData[60] = 12
	assert.NoError(t, netINodeDriverClient.PWriteWithMem(uNetINode, writeData, 612))
	assert.NoError(t, netINodeDriverClient.Flush(uNetINode))
	readData := make([]byte, 73)
	assert.NoError(t, netINodeDriverClient.PReadWithMem(uNetINode, readData, 612))
	assert.Equal(t, writeData, readData)

	time.Sleep(time.Microsecond * 800)
	for i = 0; i < len(dataNodeSRPCListenAddrs); i++ {
		assert.NoError(t, dataNodes[i].Close())
	}
	assert.NoError(t, nameNode.Close())
}
