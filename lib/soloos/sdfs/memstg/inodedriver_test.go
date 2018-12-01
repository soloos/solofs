package memstg

import (
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/snet"
	"soloos/util"
	"soloos/util/offheap"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func InitMemBlockDriversForTest(t *testing.T,
	memBlockDriver *MemBlockDriver, offheapDriver *offheap.OffheapDriver, blockChunkSize int) {
	memBlockDriverOptions := MemBlockDriverOptions{
		[]MemBlockPoolOptions{
			MemBlockPoolOptions{
				blockChunkSize, 1024,
			},
		},
	}
	assert.NoError(t, memBlockDriver.Init(memBlockDriverOptions, offheapDriver))
}

func InitNetBlockDriversForTest(t *testing.T,
	netBlockDriver *netstg.NetBlockDriver,
	offheapDriver *offheap.OffheapDriver,
	snetDriver *snet.SNetDriver,
	snetClientDriver *snet.ClientDriver) {
	netBlockDriverOptions := netstg.NetBlockDriverOptions{
		netstg.NetBlockPoolOptions{
			int32(-1),
		},
	}
	assert.NoError(t, netBlockDriver.Init(netBlockDriverOptions, offheapDriver, snetDriver, snetClientDriver))
}

func InitDriversForTest(t *testing.T,
	mockServerAddr string,
	mockServer *netstg.MockServer,
	memBlockDriver *MemBlockDriver,
	inodeDriver *INodeDriver,
	blockChunkSize int) {
	var (
		offheapDriver         = &offheap.DefaultOffheapDriver
		snetDriver            snet.SNetDriver
		snetClientDriver      snet.ClientDriver
		netBlockDriverOptions = netstg.NetBlockDriverOptions{
			netstg.NetBlockPoolOptions{
				int32(-1),
			},
		}
		netBlockDriver netstg.NetBlockDriver
	)

	assert.NoError(t, mockServer.Init("tcp", mockServerAddr))
	go func() {
		util.AssertErrIsNil(mockServer.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	assert.NoError(t, snetDriver.Init(offheapDriver))
	assert.NoError(t, snetClientDriver.Init(offheapDriver))
	InitMemBlockDriversForTest(t, memBlockDriver, offheapDriver, blockChunkSize)
	InitNetBlockDriversForTest(t, &netBlockDriver, offheapDriver, &snetDriver, &snetClientDriver)

	assert.NoError(t, netBlockDriver.Init(netBlockDriverOptions, offheapDriver, &snetDriver, &snetClientDriver))
	assert.NoError(t, inodeDriver.Init(-1, offheapDriver, &netBlockDriver, memBlockDriver))
}

func InitInodeForTest(t *testing.T, inodeDriver *INodeDriver, netBlockSize, memBlockSize int) types.INodeUintptr {
	var (
		inodeID types.INodeID
	)

	util.InitUUID64(&inodeID)

	uINode, _ := inodeDriver.MustGetINode(inodeID)
	uINode.Ptr().NetBlockSize = netBlockSize
	uINode.Ptr().MemBlockSize = memBlockSize

	return uINode
}
