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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	go util.PProfServe("192.168.56.100:17221")
	var (
		nameNode                  namenode.NameNode
		nameNodeSRPCPeerID        = snet.MakeSysPeerID("NameNodeSRPCForTest")
		nameNodeWebPeerID         = snet.MakeSysPeerID("NameNodeWebForTest")
		nameNodeSRPCListenAddr    = "127.0.0.1:10401"
		nameNodeWebListenAddr     = "127.0.0.1:10402"
		netDriverServerListenAddr = "127.0.0.1:10403"
		netDriverServerServeAddr  = "http://127.0.0.1:10403"
		metaStgForNameNode        metastg.MetaStg

		dataNodes               [6]DataNode
		dataNodeSRPCPeerIDs     [6]snettypes.PeerID
		dataNodeSRPCListenAddrs = []string{
			"127.0.0.1:10410",
			"127.0.0.1:10411",
			"127.0.0.1:10412",
			"127.0.0.1:10413",
			"127.0.0.1:10414",
			"127.0.0.1:10415",
		}
	)

	var (
		soloOSEnvForClient      soloosbase.SoloOSEnv
		memBlockDriverForClient memstg.MemBlockDriver
		netBlockDriverForClient memstg.NetBlockDriver
		netINodeDriverForClient memstg.NetINodeDriver

		soloOSEnvForNameNode      soloosbase.SoloOSEnv
		memBlockDriverForNameNode memstg.MemBlockDriver
		netBlockDriverForNameNode memstg.NetBlockDriver
		netINodeDriverForNameNode memstg.NetINodeDriver

		soloOSEnvForDataNodes      [6]soloosbase.SoloOSEnv
		memBlockDriverForDataNodes [6]memstg.MemBlockDriver
		netBlockDriverForDataNodes [6]memstg.NetBlockDriver
		netINodeDriverForDataNodes [6]memstg.NetINodeDriver

		netBlockCap int   = 32
		memBlockCap int   = 16
		blocksLimit int32 = 4
		uNetINode   sdfsapitypes.NetINodeUintptr
		i           int
		err         error
	)

	assert.NoError(t, soloOSEnvForNameNode.InitWithSNet(""))
	go func() {
		assert.NoError(t, soloOSEnvForNameNode.SNetDriver.PrepareServer(netDriverServerListenAddr,
			netDriverServerServeAddr,
			nil, nil))
		assert.NoError(t, soloOSEnvForNameNode.SNetDriver.ServerServe())
	}()
	time.Sleep(100 * time.Millisecond)
	metastg.MakeMetaStgForTest(&soloOSEnvForNameNode, &metaStgForNameNode)

	assert.NoError(t, soloOSEnvForClient.InitWithSNet(netDriverServerServeAddr))

	memstg.MemStgMakeDriversForTest(&soloOSEnvForClient,
		nameNodeSRPCListenAddr,
		&memBlockDriverForClient, &netBlockDriverForClient, &netINodeDriverForClient, memBlockCap, blocksLimit)

	memstg.MemStgMakeDriversForTest(&soloOSEnvForNameNode,
		nameNodeSRPCListenAddr,
		&memBlockDriverForNameNode, &netBlockDriverForNameNode, &netINodeDriverForNameNode, memBlockCap, blocksLimit)
	namenode.MakeNameNodeForTest(&soloOSEnvForNameNode, &nameNode, &metaStgForNameNode,
		nameNodeSRPCPeerID, nameNodeSRPCListenAddr,
		nameNodeWebPeerID, nameNodeWebListenAddr,
		&memBlockDriverForNameNode, &netBlockDriverForNameNode, &netINodeDriverForNameNode)

	go func() {
		assert.NoError(t, nameNode.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	for i = 0; i < len(dataNodeSRPCListenAddrs); i++ {
		assert.NoError(t, soloOSEnvForDataNodes[i].InitWithSNet(netDriverServerServeAddr))
		dataNodeSRPCPeerIDs[i] = snet.MakeSysPeerID(fmt.Sprintf("DataNodeForTest_%v", i))

		memstg.MemStgMakeDriversForTest(&soloOSEnvForDataNodes[i],
			nameNodeSRPCListenAddr,
			&memBlockDriverForDataNodes[i],
			&netBlockDriverForDataNodes[i],
			&netINodeDriverForDataNodes[i],
			memBlockCap, blocksLimit)

		MakeDataNodeForTest(&soloOSEnvForDataNodes[i],
			&dataNodes[i],
			dataNodeSRPCPeerIDs[i], dataNodeSRPCListenAddrs[i],
			nameNodeSRPCPeerID, nameNodeSRPCListenAddr,
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
