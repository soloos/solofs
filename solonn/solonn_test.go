package solonn

import (
	"soloos/common/solofsapitypes"
	"soloos/common/snet"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/solofs/memstg"
	"soloos/solofs/metastg"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	var (
		soloOSEnvForCommon     soloosbase.SoloOSEnv
		soloOSEnvForMetaStg    soloosbase.SoloOSEnv
		metaStg                metastg.MetaStg
		solonn               Solonn
		solonnSRPCPeerID     = snet.MakeSysPeerID("SolonnSRPCForTest")
		solonnSRPCListenAddr = "127.0.0.1:10300"
		solonnWebPeerID      = snet.MakeSysPeerID("SolonnWebForTest")
		solonnWebListenAddr  = "127.0.0.1:10301"
		mockServerAddr         = "127.0.0.1:10302"
		mockServer             memstg.MockServer

		soloOSEnvForClient      soloosbase.SoloOSEnv
		memBlockDriverForClient memstg.MemBlockDriver
		netBlockDriverForClient memstg.NetBlockDriver
		netINodeDriverForClient memstg.NetINodeDriver

		soloOSEnvForServer      soloosbase.SoloOSEnv
		memBlockDriverForServer memstg.MemBlockDriver
		netBlockDriverForServer memstg.NetBlockDriver
		netINodeDriverForServer memstg.NetINodeDriver

		netBlockCap int   = 1024
		memBlockCap int   = 128
		blocksLimit int32 = 4
		uNetINode   solofsapitypes.NetINodeUintptr
		i           int
		err         error
	)

	assert.NoError(t, soloOSEnvForCommon.InitWithSNet(""))
	assert.NoError(t, soloOSEnvForMetaStg.InitWithSNet(""))
	metastg.MakeMetaStgForTest(&soloOSEnvForMetaStg, &metaStg)

	assert.NoError(t, soloOSEnvForClient.InitWithSNet(""))
	assert.NoError(t, soloOSEnvForServer.InitWithSNet(""))

	memstg.MemStgMakeDriversForTest(&soloOSEnvForClient,
		solonnSRPCListenAddr,
		&memBlockDriverForClient, &netBlockDriverForClient, &netINodeDriverForClient, memBlockCap, blocksLimit)

	memstg.MemStgMakeDriversForTest(&soloOSEnvForServer,
		solonnSRPCListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer, memBlockCap, blocksLimit)
	MakeSolonnForTest(&soloOSEnvForServer, &solonn, &metaStg,
		solonnSRPCPeerID, solonnSRPCListenAddr,
		solonnWebPeerID, solonnWebListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer)
	go func() {
		assert.NoError(t, solonn.Serve())
	}()
	time.Sleep(time.Millisecond * 300)
	memstg.MakeMockServerForTest(&soloOSEnvForCommon, mockServerAddr, &mockServer)

	for i = 0; i < 6; i++ {
		var peer snettypes.Peer
		snettypes.InitTmpPeerID((*snettypes.PeerID)(&peer.ID))
		peer.SetAddress(mockServerAddr)
		peer.ServiceProtocol = solofsapitypes.DefaultSOLOFSRPCProtocol
		solonn.SolodnRegister(peer)
	}

	var netINodeID solofsapitypes.NetINodeID
	solofsapitypes.InitTmpNetINodeID(&netINodeID)
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

	assert.NoError(t, solonn.Close())
	assert.NoError(t, mockServer.Close())
}
