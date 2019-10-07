package memstg

import (
	"soloos/common/snet"
	"soloos/common/solofsapi"
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetBlockDriver(t *testing.T) {
	var (
		soloosEnv         soloosbase.SoloosEnv
		mockNetINodeTable MockNetINodeTable
		mockMemBlockTable MockMemBlockTable
		mockServer        MockServer
		solonnClient      solofsapi.SolonnClient
		solodnClient      solofsapi.SolodnClient
		netBlockDriver    NetBlockDriver
	)
	assert.NoError(t, soloosEnv.InitWithSNet(""))
	mockServerAddr := "127.0.0.1:10021"
	assert.NoError(t, mockNetINodeTable.Init(&soloosEnv))
	assert.NoError(t, mockMemBlockTable.Init(&soloosEnv, 1024))
	NetStgMakeDriversWithMockServerForTest(&soloosEnv,
		mockServerAddr, &mockServer,
		&solonnClient, &solodnClient,
		&netBlockDriver)

	var peer0 snet.Peer
	peer0.ID = snet.MakeSysPeerID("Peer0ForTest")
	peer0.SetAddress(mockServerAddr)
	peer0.ServiceProtocol = solofstypes.DefaultSolofsRPCProtocol
	soloosEnv.SNetDriver.RegisterPeer(peer0)

	var peer1 snet.Peer
	peer1.ID = snet.MakeSysPeerID("Peer0ForTest")
	peer1.SetAddress(mockServerAddr)
	peer1.ServiceProtocol = solofstypes.DefaultSolofsRPCProtocol
	soloosEnv.SNetDriver.RegisterPeer(peer1)

	data := make([]byte, 8)
	for i := 0; i < len(data); i++ {
		data[i] = 1
	}

	uNetINode := mockNetINodeTable.AllocNetINode(1024, 128)
	defer mockNetINodeTable.ReleaseNetINode(uNetINode)

	netBlockIndex := int32(10)
	uNetBlock, err := netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
	defer netBlockDriver.ReleaseNetBlock(uNetBlock)
	assert.NoError(t, err)
	uNetBlock.Ptr().StorDataBackends.Append(peer0.ID)
	uNetBlock.Ptr().StorDataBackends.Append(peer1.ID)
	uMemBlock := mockMemBlockTable.AllocMemBlock()
	memBlockIndex := int32(0)
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, 0, 12))
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, 11, 24))
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, 30, 64))
	uMemBlock.Ptr().UploadJob.SyncDataSig.Wait()

	assert.NoError(t, mockServer.Close())
}
