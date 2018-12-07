package memstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
	"soloos/snet"
	"soloos/util/offheap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func InitMemBlockDriversForTest(t *testing.T,
	memBlockDriver *MemBlockDriver, offheapDriver *offheap.OffheapDriver,
	blockChunkSize int, blockChunksLimit int32) {
	memBlockDriverOptions := MemBlockDriverOptions{
		[]MemBlockPoolOptions{
			MemBlockPoolOptions{
				blockChunkSize, blockChunksLimit,
			},
		},
	}
	assert.NoError(t, memBlockDriver.Init(memBlockDriverOptions, offheapDriver))
}

func InitDriversForTest(t *testing.T,
	nameNodeSRPCServerAddr string,
	memBlockDriver *MemBlockDriver,
	inodeDriver *INodeDriver,
	blockChunkSize int, blockChunksLimit int32) {
	var (
		offheapDriver    = &offheap.DefaultOffheapDriver
		snetDriver       snet.SNetDriver
		snetClientDriver snet.ClientDriver
		nameNodeClient   api.NameNodeClient
		dataNodeClient   api.DataNodeClient
		netBlockDriver   netstg.NetBlockDriver
	)

	netstg.InitDriversForTest(t,
		&snetDriver, &snetClientDriver,
		nameNodeSRPCServerAddr,
		&nameNodeClient, &dataNodeClient,
		&netBlockDriver,
	)

	InitMemBlockDriversForTest(t, memBlockDriver, offheapDriver, blockChunkSize, blockChunksLimit)

	assert.NoError(t, inodeDriver.Init(offheapDriver, &netBlockDriver, memBlockDriver))
}

func InitDriversWithMockServerForTest(t *testing.T,
	mockServerAddr string,
	mockServer *netstg.MockServer,
	memBlockDriver *MemBlockDriver,
	inodeDriver *INodeDriver,
	blockChunkSize int, blockChunksLimit int32) {
	var (
		offheapDriver    = &offheap.DefaultOffheapDriver
		snetDriver       snet.SNetDriver
		snetClientDriver snet.ClientDriver
		nameNodeClient   api.NameNodeClient
		dataNodeClient   api.DataNodeClient
		netBlockDriver   netstg.NetBlockDriver
	)

	netstg.InitDriversWithMockServerForTest(t,
		&snetDriver, &snetClientDriver,
		mockServerAddr, mockServer,
		&nameNodeClient, &dataNodeClient,
		&netBlockDriver,
	)

	InitMemBlockDriversForTest(t, memBlockDriver, offheapDriver, blockChunkSize, blockChunksLimit)

	assert.NoError(t, inodeDriver.Init(offheapDriver, &netBlockDriver, memBlockDriver))
}
