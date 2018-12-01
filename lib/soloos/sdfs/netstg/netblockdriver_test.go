package netstg

import (
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestNetBlockDriver(t *testing.T) {
	var (
		offheapDriver         = &offheap.DefaultOffheapDriver
		mockMemBlockPool      types.MockMemBlockPool
		netBlockDriverOptions = NetBlockDriverOptions{
			NetBlockPoolOptions{
				int32(-1),
			},
		}
		mockServer       MockServer
		netBlockDriver   NetBlockDriver
		snetDriver       snet.SNetDriver
		snetClientDriver snet.ClientDriver
	)
	mockServerAddr := "127.0.0.1:10021"

	assert.NoError(t, mockServer.Init("tcp", mockServerAddr))
	go func() {
		util.AssertErrIsNil(mockServer.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	assert.NoError(t, mockMemBlockPool.Init(offheapDriver, 1024))
	assert.NoError(t, snetDriver.Init(offheapDriver))
	assert.NoError(t, snetClientDriver.Init(offheapDriver))

	var uPeer0 = snetDriver.NewPeer()
	util.InitUUID64(&uPeer0.Ptr().ID)
	uPeer0.Ptr().SetAddress(mockServerAddr)
	uPeer0.Ptr().ServiceProtocol = snettypes.ProtocolSRPC

	var uPeer1 = snetDriver.NewPeer()
	util.InitUUID64(&uPeer1.Ptr().ID)
	uPeer1.Ptr().SetAddress(mockServerAddr)
	uPeer1.Ptr().ServiceProtocol = snettypes.ProtocolSRPC

	assert.NoError(t, netBlockDriver.Init(netBlockDriverOptions, offheapDriver, &snetDriver, &snetClientDriver))

	data := make([]byte, 8)
	for i := 0; i < len(data); i++ {
		data[i] = 1
	}

	var (
		inode  types.INode
		uINode types.INodeUintptr = types.INodeUintptr((unsafe.Pointer(&inode)))
	)
	uINode.Ptr().NetBlockSize = 1024
	uINode.Ptr().MemBlockSize = 128

	uNetBlock, _ := netBlockDriver.MustGetBlock(uINode, 10)
	uNetBlock.Ptr().DataNodes.Append(uPeer0)
	uNetBlock.Ptr().DataNodes.Append(uPeer1)
	uMemBlock := mockMemBlockPool.AllocMemBlock()
	memBlockIndex := 0
	assert.NoError(t, netBlockDriver.PWrite(uINode, uNetBlock, uMemBlock, memBlockIndex, 0, 12))
	assert.NoError(t, netBlockDriver.PWrite(uINode, uNetBlock, uMemBlock, memBlockIndex, 11, 24))
	assert.NoError(t, netBlockDriver.PWrite(uINode, uNetBlock, uMemBlock, memBlockIndex, 30, 64))
	assert.NoError(t, netBlockDriver.Flush(uMemBlock))

	assert.NoError(t, mockServer.Close())
}
