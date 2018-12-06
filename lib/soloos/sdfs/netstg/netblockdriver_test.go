package netstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/snet"
	"soloos/util/offheap"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestNetBlockDriver(t *testing.T) {
	var (
		offheapDriver         = &offheap.DefaultOffheapDriver
		mockMemBlockPool      types.MockMemBlockPool
		netBlockDriverOptions = NetBlockDriverOptions{
			int32(-1),
		}
		snetDriver       snet.SNetDriver
		snetClientDriver snet.ClientDriver
		mockServer       MockServer
		nameNodeClient   api.NameNodeClient
		dataNodeClient   api.DataNodeClient
		netBlockDriver   NetBlockDriver
	)
	mockServerAddr := "127.0.0.1:10021"
	mockMemBlockPool.Init(&offheap.DefaultOffheapDriver, 1024)
	InitDriversForTest(t, &snetDriver, &snetClientDriver,
		mockServerAddr, &mockServer,
		&nameNodeClient, &dataNodeClient,
		&netBlockDriver)

	var uPeer0 = snetDriver.MustGetPeer(nil, mockServerAddr, types.DefaultSDFSRPCProtocol)
	var uPeer1 = snetDriver.MustGetPeer(nil, mockServerAddr, types.DefaultSDFSRPCProtocol)

	assert.NoError(t, netBlockDriver.Init(netBlockDriverOptions, offheapDriver,
		&snetDriver, &snetClientDriver,
		&nameNodeClient, &dataNodeClient,
	))

	data := make([]byte, 8)
	for i := 0; i < len(data); i++ {
		data[i] = 1
	}

	var (
		inode  types.INode
		uINode types.INodeUintptr = types.INodeUintptr((unsafe.Pointer(&inode)))
	)
	uINode.Ptr().NetBlockCap = 1024
	uINode.Ptr().MemBlockCap = 128

	uNetBlock := netBlockDriver.MustGetBlock(uINode, 10)
	uNetBlock.Ptr().DataNodes.Append(uPeer0)
	uNetBlock.Ptr().DataNodes.Append(uPeer1)
	uMemBlock := mockMemBlockPool.AllocMemBlock()
	memBlockIndex := 0
	assert.NoError(t, netBlockDriver.PWrite(uINode, uNetBlock, uMemBlock, memBlockIndex, 0, 12))
	assert.NoError(t, netBlockDriver.PWrite(uINode, uNetBlock, uMemBlock, memBlockIndex, 11, 24))
	assert.NoError(t, netBlockDriver.PWrite(uINode, uNetBlock, uMemBlock, memBlockIndex, 30, 64))
	assert.NoError(t, netBlockDriver.FlushMemBlock(uMemBlock))

	assert.NoError(t, mockServer.Close())
}
