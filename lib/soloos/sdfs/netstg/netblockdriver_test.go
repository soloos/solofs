package netstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/snet"
	"soloos/util/offheap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetBlockDriver(t *testing.T) {
	var (
		offheapDriver    = &offheap.DefaultOffheapDriver
		mockNetINodePool types.MockNetINodePool
		mockMemBlockPool types.MockMemBlockPool
		snetDriver       snet.SNetDriver
		snetClientDriver snet.ClientDriver
		mockServer       MockServer
		nameNodeClient   api.NameNodeClient
		dataNodeClient   api.DataNodeClient
		netBlockDriver   NetBlockDriver
	)
	mockServerAddr := "127.0.0.1:10021"
	assert.NoError(t, mockNetINodePool.Init(&offheap.DefaultOffheapDriver))
	assert.NoError(t, mockMemBlockPool.Init(offheapDriver, 1024))
	MakeDriversWithMockServerForTest(t, &snetDriver, &snetClientDriver,
		mockServerAddr, &mockServer,
		&nameNodeClient, &dataNodeClient,
		&netBlockDriver)

	var uPeer0, _ = snetDriver.MustGetPeer(nil, mockServerAddr, types.DefaultSDFSRPCProtocol)
	var uPeer1, _ = snetDriver.MustGetPeer(nil, mockServerAddr, types.DefaultSDFSRPCProtocol)

	data := make([]byte, 8)
	for i := 0; i < len(data); i++ {
		data[i] = 1
	}

	uNetINode := mockNetINodePool.AllocNetINode(1024, 128)

	netBlockIndex := 10
	uNetBlock, err := netBlockDriver.MustGetBlock(uNetINode, netBlockIndex)
	assert.NoError(t, err)
	uNetBlock.Ptr().StorDataBackends.Append(uPeer0)
	uNetBlock.Ptr().StorDataBackends.Append(uPeer1)
	uMemBlock := mockMemBlockPool.AllocMemBlock()
	memBlockIndex := 0
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, 0, 12))
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, 11, 24))
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, 30, 64))
	assert.NoError(t, netBlockDriver.FlushMemBlock(uNetINode, uNetBlock, uMemBlock))

	assert.NoError(t, mockServer.Close())
}
