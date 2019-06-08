package datanode

import (
	"fmt"
	"soloos/common/sdfsapitypes"
	"soloos/common/snet"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/namenode"
	"soloos/sdfs/netstg"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	go util.PProfServe("192.168.56.100:17221")
	var (
		nameNode                  namenode.NameNode
		nameNodePeerID            = snet.MakeSysPeerID("NameNodeForTest")
		nameNodeSRPCListenAddr    = "127.0.0.1:10401"
		netDriverServerListenAddr = "127.0.0.1:10402"
		netDriverServerServeAddr  = "http://127.0.0.1:10402"
		metaStgForNameNode        metastg.MetaStg

		dataNodes               [6]DataNode
		dataNodePeerIDs         [6]snettypes.PeerID
		dataNodeSRPCListenAddrs = []string{
			"127.0.0.1:10410",
			"127.0.0.1:10411",
			"127.0.0.1:10412",
			"127.0.0.1:10413",
			"127.0.0.1:10414",
			"127.0.0.1:10415",
		}
		metaStgForDataNodes [6]metastg.MetaStg
	)

	var (
		soloOSEnvForClient      soloosbase.SoloOSEnv
		memBlockDriverForClient memstg.MemBlockDriver
		netBlockDriverForClient netstg.NetBlockDriver
		netINodeDriverForClient memstg.NetINodeDriver

		soloOSEnvForNameNode      soloosbase.SoloOSEnv
		memBlockDriverForNameNode memstg.MemBlockDriver
		netBlockDriverForNameNode netstg.NetBlockDriver
		netINodeDriverForNameNode memstg.NetINodeDriver

		soloOSEnvForDataNodes      [6]soloosbase.SoloOSEnv
		memBlockDriverForDataNodes [6]memstg.MemBlockDriver
		netBlockDriverForDataNodes [6]netstg.NetBlockDriver
		netINodeDriverForDataNodes [6]memstg.NetINodeDriver

		netBlockCap int   = 32
		memBlockCap int   = 16
		blocksLimit int32 = 4
		uNetINode   sdfsapitypes.NetINodeUintptr
		i           int
		err         error
	)

	assert.NoError(t, soloOSEnvForNameNode.Init())
	go func() {
		assert.NoError(t, soloOSEnvForNameNode.SNetDriver.StartServer(netDriverServerListenAddr,
			netDriverServerServeAddr,
			nil, nil))
	}()
	metastg.MakeMetaStgForTest(&soloOSEnvForNameNode, &metaStgForNameNode)

	assert.NoError(t, soloOSEnvForClient.Init())
	assert.NoError(t, soloOSEnvForClient.SNetDriver.StartClient(netDriverServerServeAddr))

	memstg.MakeDriversForTest(&soloOSEnvForClient,
		nameNodeSRPCListenAddr,
		&memBlockDriverForClient, &netBlockDriverForClient, &netINodeDriverForClient, memBlockCap, blocksLimit)

	memstg.MakeDriversForTest(&soloOSEnvForNameNode,
		nameNodeSRPCListenAddr,
		&memBlockDriverForNameNode, &netBlockDriverForNameNode, &netINodeDriverForNameNode, memBlockCap, blocksLimit)
	namenode.MakeNameNodeForTest(&soloOSEnvForNameNode, &nameNode, &metaStgForNameNode,
		nameNodePeerID, nameNodeSRPCListenAddr,
		&memBlockDriverForNameNode, &netBlockDriverForNameNode, &netINodeDriverForNameNode)

	go func() {
		assert.NoError(t, nameNode.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	for i = 0; i < len(dataNodeSRPCListenAddrs); i++ {
		assert.NoError(t, soloOSEnvForDataNodes[i].Init())
		assert.NoError(t, soloOSEnvForDataNodes[i].SNetDriver.StartClient(netDriverServerServeAddr))
		metastg.MakeMetaStgForTest(&soloOSEnvForDataNodes[i], &metaStgForDataNodes[i])
		dataNodePeerIDs[i] = snet.MakeSysPeerID(fmt.Sprintf("DataNodeForTest_%v", i))

		memstg.MakeDriversForTest(&soloOSEnvForDataNodes[i],
			nameNodeSRPCListenAddr,
			&memBlockDriverForDataNodes[i],
			&netBlockDriverForDataNodes[i],
			&netINodeDriverForDataNodes[i],
			memBlockCap, blocksLimit)

		MakeDataNodeForTest(&soloOSEnvForDataNodes[i],
			&dataNodes[i],
			dataNodePeerIDs[i], dataNodeSRPCListenAddrs[i],
			&metaStgForDataNodes[i],
			nameNodePeerID, nameNodeSRPCListenAddr,
			&memBlockDriverForDataNodes[i],
			&netBlockDriverForDataNodes[i],
			&netINodeDriverForDataNodes[i])
		go func(localI int) {
			assert.NoError(t, dataNodes[localI].Serve())
		}(i)
	}
	time.Sleep(time.Millisecond * 300)

	var (
		netINodeID sdfsapitypes.NetINodeID
	)
	sdfsapitypes.InitTmpNetINodeID(&netINodeID)
	uNetINode, err = netINodeDriverForClient.MustGetNetINode(netINodeID, 0, netBlockCap, memBlockCap)
	defer netINodeDriverForClient.ReleaseNetINode(uNetINode)
	assert.NoError(t, err)

	writeData := make([]byte, 73)
	writeData[3] = 1
	writeData[7] = 2
	writeData[8] = 3
	writeData[33] = 4
	writeData[60] = 5
	assert.NoError(t, netINodeDriverForClient.PWriteWithMem(uNetINode, writeData, 612))
	assert.NoError(t, netINodeDriverForClient.Sync(uNetINode))

	var maxWriteTimes int = 128
	for i = 0; i < maxWriteTimes; i++ {
		assert.NoError(t, netINodeDriverForClient.PWriteWithMem(uNetINode, writeData, uint64(netBlockCap*600+8*i)))
	}

	readData := make([]byte, 73)
	_, err = netINodeDriverForClient.PReadWithMem(uNetINode, readData, 612)
	assert.NoError(t, err)
	assert.Equal(t, writeData, readData)

	time.Sleep(time.Microsecond * 600)
	for i = 0; i < len(dataNodeSRPCListenAddrs); i++ {
		assert.NoError(t, dataNodes[i].Close())
	}
	assert.NoError(t, nameNode.Close())
}
