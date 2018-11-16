package netstg

import (
	"soloos/sdfs/types"
	"soloos/snet"
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
		mockMemBlockPool      MockMemBlockPool
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

	assert.NoError(t, mockServer.Init("tcp", MockServerAddr))
	go func() {
		util.AssertErrIsNil(mockServer.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	assert.NoError(t, mockMemBlockPool.Init(offheapDriver, 1024))
	assert.NoError(t, snetDriver.Init(offheapDriver))
	assert.NoError(t, snetClientDriver.Init(offheapDriver))

	var uPeer = snetDriver.NewPeer()
	uPeer.Ptr().SetAddress(MockServerAddr)
	uPeer.Ptr().SetServiceProtocol("srpc")
	assert.NoError(t, snetClientDriver.RegisterPeer(uPeer))

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
	uINode.Ptr().MemBlockSize = 1024

	uNetBlock, _ := netBlockDriver.MustGetBlock(uINode, 10)
	uNetBlock.Ptr().DataNodes.Append(uPeer)
	uMemBlock := mockMemBlockPool.AllocMemBlock()
	assert.NoError(t, netBlockDriver.PWrite(uINode, uNetBlock, uMemBlock, 0, 0, 12))
	assert.NoError(t, netBlockDriver.PWrite(uINode, uNetBlock, uMemBlock, 0, 11, 24))
	assert.NoError(t, netBlockDriver.PWrite(uINode, uNetBlock, uMemBlock, 0, 30, 64))
	assert.NoError(t, netBlockDriver.Flush(uMemBlock))

	assert.NoError(t, mockServer.Close())
}
