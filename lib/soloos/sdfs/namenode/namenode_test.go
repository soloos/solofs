package namenode

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/snet"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	var (
		soloOSEnvForCommon     soloosbase.SoloOSEnv
		soloOSEnvForMetaStg    soloosbase.SoloOSEnv
		metaStg                metastg.MetaStg
		nameNode               NameNode
		nameNodePeerID         = snet.MakeSysPeerID("NameNodeForTest")
		nameNodeSRPCListenAddr = "127.0.0.1:10300"
		mockServerAddr         = "127.0.0.1:10301"
		mockServer             netstg.MockServer

		soloOSEnvForClient      soloosbase.SoloOSEnv
		memBlockDriverForClient memstg.MemBlockDriver
		netBlockDriverForClient netstg.NetBlockDriver
		netINodeDriverForClient memstg.NetINodeDriver

		soloOSEnvForServer      soloosbase.SoloOSEnv
		memBlockDriverForServer memstg.MemBlockDriver
		netBlockDriverForServer netstg.NetBlockDriver
		netINodeDriverForServer memstg.NetINodeDriver

		netBlockCap int   = 1024
		memBlockCap int   = 128
		blocksLimit int32 = 4
		uNetINode   sdfsapitypes.NetINodeUintptr
		i           int
		err         error
	)

	assert.NoError(t, soloOSEnvForCommon.Init())
	assert.NoError(t, soloOSEnvForMetaStg.Init())
	metastg.MakeMetaStgForTest(&soloOSEnvForMetaStg, &metaStg)

	assert.NoError(t, soloOSEnvForClient.Init())
	assert.NoError(t, soloOSEnvForServer.Init())

	memstg.MakeDriversForTest(&soloOSEnvForClient,
		nameNodeSRPCListenAddr,
		&memBlockDriverForClient, &netBlockDriverForClient, &netINodeDriverForClient, memBlockCap, blocksLimit)

	memstg.MakeDriversForTest(&soloOSEnvForServer,
		nameNodeSRPCListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer, memBlockCap, blocksLimit)
	MakeNameNodeForTest(&soloOSEnvForServer, &nameNode, &metaStg,
		nameNodePeerID, nameNodeSRPCListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer)
	go func() {
		assert.NoError(t, nameNode.Serve())
	}()
	time.Sleep(time.Millisecond * 300)
	netstg.MakeMockServerForTest(&soloOSEnvForCommon, mockServerAddr, &mockServer)

	for i = 0; i < 6; i++ {
		var peer snettypes.Peer
		snettypes.InitTmpPeerID((*snettypes.PeerID)(&peer.ID))
		peer.SetAddress(mockServerAddr)
		peer.ServiceProtocol = sdfsapitypes.DefaultSDFSRPCProtocol
		nameNode.RegisterDataNode(peer)
	}

	var netINodeID sdfsapitypes.NetINodeID
	sdfsapitypes.InitTmpNetINodeID(&netINodeID)
	uNetINode, err = netINodeDriverForClient.MustGetNetINode(netINodeID, 0, netBlockCap, memBlockCap)
	defer netINodeDriverForClient.ReleaseNetINode(uNetINode)
	assert.NoError(t, err)

	var (
		testData = make([]byte, 93)
	)

	assert.NoError(t, netINodeDriverForClient.PWriteWithMem(uNetINode, testData[0:12], 0))
	assert.NoError(t, netINodeDriverForClient.PWriteWithMem(uNetINode, testData[11:24], 24))
	assert.NoError(t, netINodeDriverForClient.PWriteWithMem(uNetINode, testData[30:64], 64))
	assert.NoError(t, netINodeDriverForClient.Sync(uNetINode))
	_, err = netINodeDriverForClient.PReadWithMem(uNetINode, testData, 73)
	assert.NoError(t, err)

	assert.NoError(t, nameNode.Close())
	assert.NoError(t, mockServer.Close())
}
