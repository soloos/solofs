package solonn

import (
	"soloos/common/snet"
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/solofs/memstg"
	"soloos/solofs/metastg"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	var (
		soloosEnvForCommon   soloosbase.SoloosEnv
		soloosEnvForMetaStg  soloosbase.SoloosEnv
		metaStg              metastg.MetaStg
		solonn               Solonn
		solonnSrpcPeerID     = snet.MakeSysPeerID("SolonnSrpcForTest")
		solonnSrpcListenAddr = "127.0.0.1:10300"
		solonnWebPeerID      = snet.MakeSysPeerID("SolonnWebForTest")
		solonnWebListenAddr  = "127.0.0.1:10301"
		mockServerAddr       = "127.0.0.1:10302"
		mockServer           memstg.MockServer

		soloosEnvForClient      soloosbase.SoloosEnv
		memBlockDriverForClient memstg.MemBlockDriver
		netBlockDriverForClient memstg.NetBlockDriver
		netINodeDriverForClient memstg.NetINodeDriver

		soloosEnvForServer      soloosbase.SoloosEnv
		memBlockDriverForServer memstg.MemBlockDriver
		netBlockDriverForServer memstg.NetBlockDriver
		netINodeDriverForServer memstg.NetINodeDriver

		netBlockCap int   = 1024
		memBlockCap int   = 128
		blocksLimit int32 = 4
		uNetINode   solofstypes.NetINodeUintptr
		i           int
		err         error
	)

	assert.NoError(t, soloosEnvForCommon.InitWithSNet(""))
	assert.NoError(t, soloosEnvForMetaStg.InitWithSNet(""))
	metastg.MakeMetaStgForTest(&soloosEnvForMetaStg, &metaStg)

	assert.NoError(t, soloosEnvForClient.InitWithSNet(""))
	assert.NoError(t, soloosEnvForServer.InitWithSNet(""))

	memstg.MemStgMakeDriversForTest(&soloosEnvForClient,
		solonnSrpcListenAddr,
		&memBlockDriverForClient, &netBlockDriverForClient, &netINodeDriverForClient, memBlockCap, blocksLimit)

	memstg.MemStgMakeDriversForTest(&soloosEnvForServer,
		solonnSrpcListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer, memBlockCap, blocksLimit)
	MakeSolonnForTest(&soloosEnvForServer, &solonn, &metaStg,
		solonnSrpcPeerID, solonnSrpcListenAddr,
		solonnWebPeerID, solonnWebListenAddr,
		&memBlockDriverForServer, &netBlockDriverForServer, &netINodeDriverForServer)
	go func() {
		assert.NoError(t, solonn.Serve())
	}()
	time.Sleep(time.Millisecond * 300)
	memstg.MakeMockServerForTest(&soloosEnvForCommon, mockServerAddr, &mockServer)

	for i = 0; i < 6; i++ {
		var peer snet.Peer
		snet.InitTmpPeerID((*snet.PeerID)(&peer.ID))
		peer.SetAddress(mockServerAddr)
		peer.ServiceProtocol = solofstypes.DefaultSolofsRPCProtocol
		solonn.SolodnRegister(peer)
	}

	var netINodeID solofstypes.NetINodeID
	solofstypes.InitTmpNetINodeID(&netINodeID)
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
