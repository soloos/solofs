package memstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
	"soloos/snet"
	"soloos/util/offheap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MakeMemBlockDriversForTest(t *testing.T,
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

func MakeDriversForTest(t *testing.T,
	snetDriver *snet.SNetDriver,
	nameNodeSRPCServerAddr string,
	memBlockDriver *MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *NetINodeDriver,
	blockChunkSize int, blockChunksLimit int32) {
	var (
		offheapDriver    = &offheap.DefaultOffheapDriver
		snetClientDriver snet.ClientDriver
		nameNodeClient   api.NameNodeClient
		dataNodeClient   api.DataNodeClient
	)

	netstg.MakeDriversForTest(t,
		snetDriver, &snetClientDriver,
		nameNodeSRPCServerAddr,
		&nameNodeClient, &dataNodeClient,
		netBlockDriver,
	)

	MakeMemBlockDriversForTest(t, memBlockDriver, offheapDriver, blockChunkSize, blockChunksLimit)

	assert.NoError(t, netINodeDriver.Init(offheapDriver, netBlockDriver, memBlockDriver, &nameNodeClient, nil, nil))
}

func MakeDriversWithMockServerForTest(t *testing.T,
	mockServerAddr string,
	mockServer *netstg.MockServer,
	snetDriver *snet.SNetDriver,
	netBlockDriver *netstg.NetBlockDriver,
	memBlockDriver *MemBlockDriver,
	netINodeDriver *NetINodeDriver,
	blockChunkSize int, blockChunksLimit int32) {
	var (
		offheapDriver    = &offheap.DefaultOffheapDriver
		snetClientDriver snet.ClientDriver
		nameNodeClient   api.NameNodeClient
		dataNodeClient   api.DataNodeClient
	)

	netstg.MakeDriversWithMockServerForTest(t,
		snetDriver, &snetClientDriver,
		mockServerAddr, mockServer,
		&nameNodeClient, &dataNodeClient,
		netBlockDriver,
	)

	MakeMemBlockDriversForTest(t, memBlockDriver, offheapDriver, blockChunkSize, blockChunksLimit)

	assert.NoError(t, netINodeDriver.Init(offheapDriver, netBlockDriver, memBlockDriver, &nameNodeClient, nil, nil))
}
