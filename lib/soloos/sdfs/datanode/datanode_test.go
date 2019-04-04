package datanode

import (
	"soloos/common/snet"
	"soloos/common/util"
	"soloos/sdbone/offheap"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/namenode"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
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

		snetDriverClient       snet.NetDriver
		snetClientDriverClient snet.ClientDriver
		memBlockDriverClient   memstg.MemBlockDriver
		netBlockDriverClient   netstg.NetBlockDriver
		netINodeDriverClient   memstg.NetINodeDriver

		snetDriverNameNode       snet.NetDriver
		snetClientDriverNameNode snet.ClientDriver
		memBlockDriverNameNode   memstg.MemBlockDriver
		netBlockDriverNameNode   netstg.NetBlockDriver
		netINodeDriverNameNode   memstg.NetINodeDriver

		snetDriverDataNodes       [6]snet.NetDriver
		snetClientDriverDataNodes [6]snet.ClientDriver
		memBlockDriverDataNodes   [6]memstg.MemBlockDriver
		netBlockDriverDataNodes   [6]netstg.NetBlockDriver
		netINodeDriverDataNodes   [6]memstg.NetINodeDriver

		netBlockCap      int   = 32
		memBlockCap      int   = 16
		blockChunksLimit int32 = 4
		uNetINode        types.NetINodeUintptr
		i                int
		err              error
	)
	metastg.MakeMetaStgForTest(offheapDriver, &metaStg)
	memstg.MakeDriversForTest(&snetDriverClient, &snetClientDriverClient,
		nameNodeSRPCListenAddr,
		&memBlockDriverClient, &netBlockDriverClient, &netINodeDriverClient, memBlockCap, blockChunksLimit)

	memstg.MakeDriversForTest(&snetDriverNameNode, &snetClientDriverNameNode,
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
		memstg.MakeDriversForTest(&snetDriverDataNodes[i], &snetClientDriverDataNodes[i],
			nameNodeSRPCListenAddr,
			&memBlockDriverDataNodes[i],
			&netBlockDriverDataNodes[i],
			&netINodeDriverDataNodes[i],
			memBlockCap, blockChunksLimit)

		MakeDataNodeForTest(&snetDriverDataNodes[i], &snetClientDriverDataNodes[i],
			&dataNodes[i], dataNodeSRPCListenAddrs[i], &metaStg,
			nameNodeSRPCListenAddr,
			&memBlockDriverDataNodes[i],
			&netBlockDriverDataNodes[i],
			&netINodeDriverDataNodes[i])
		go func(localI int) {
			assert.NoError(t, dataNodes[localI].Serve())
		}(i)
	}
	time.Sleep(time.Millisecond * 300)

	var (
		netINodeID types.NetINodeID
	)
	util.InitUUID64(&netINodeID)
	uNetINode, err = netINodeDriverClient.MustGetNetINodeWithReadAcquire(netINodeID, 0, netBlockCap, memBlockCap)
	assert.NoError(t, err)

	writeData := make([]byte, 73)
	writeData[3] = 1
	writeData[7] = 2
	writeData[8] = 3
	writeData[33] = 4
	writeData[60] = 5
	assert.NoError(t, netINodeDriverClient.PWriteWithMem(uNetINode, writeData, 612))
	assert.NoError(t, netINodeDriverClient.Flush(uNetINode))

	// var maxWriteTimes int = 128
	// for i = 0; i < maxWriteTimes; i++ {
	// assert.NoError(t, netINodeDriverClient.PWriteWithMem(uNetINode, writeData, uint64(netBlockCap*600+8*i)))
	// }

	readData := make([]byte, 73)
	_, err = netINodeDriverClient.PReadWithMem(uNetINode, readData, 612)
	assert.NoError(t, err)
	assert.Equal(t, writeData, readData)

	time.Sleep(time.Microsecond * 600)
	for i = 0; i < len(dataNodeSRPCListenAddrs); i++ {
		assert.NoError(t, dataNodes[i].Close())
	}
	assert.NoError(t, nameNode.Close())
}
