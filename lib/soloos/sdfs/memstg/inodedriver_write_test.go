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

func TestINodeDriverINodeWrite(t *testing.T) {
	var (
		offheapDriver         = &offheap.DefaultOffheapDriver
		snetDriver            snet.SNetDriver
		snetClientDriver      snet.ClientDriver
		netBlockDriverOptions = netstg.NetBlockDriverOptions{
			netstg.NetBlockPoolOptions{
				int32(-1),
			},
		}
		mockServer     netstg.MockServer
		netBlockDriver netstg.NetBlockDriver
		memBlockDriver MemBlockDriver
		inodeDriver    INodeDriver
		blockChunkSize       = 4
		maxBlocks      int32 = 4
	)

	assert.NoError(t, mockServer.Init("tcp", netstg.MockServerAddr))
	go func() {
		util.AssertErrIsNil(mockServer.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	assert.NoError(t, snetDriver.Init(offheapDriver))
	assert.NoError(t, snetClientDriver.Init(offheapDriver))
	InitMemBlockDriversForTest(t, &memBlockDriver, offheapDriver, blockChunkSize)
	InitNetBlockDriversForTest(t, &netBlockDriver, offheapDriver, &snetDriver, &snetClientDriver)

	assert.NoError(t, netBlockDriver.Init(netBlockDriverOptions, offheapDriver, &snetDriver, &snetClientDriver))
	assert.NoError(t, inodeDriver.Init(-1, offheapDriver, &netBlockDriver, &memBlockDriver))

	var (
		inodeID types.INodeID
		i       int32
	)

	uINode, _ := inodeDriver.MustGetINode(inodeID)
	uINode.Ptr().NetBlockSize = blockChunkSize
	uINode.Ptr().MemBlockSize = blockChunkSize

	for i = 0; i <= maxBlocks; i++ {
		writeOffset := int64(int64(i) * int64(blockChunkSize))

		assert.NoError(t, inodeDriver.PWrite(uINode, []byte{(byte)(i), (byte)(i * 2)}, writeOffset))

		memBlockIndex := int(writeOffset / int64(uINode.Ptr().MemBlockSize))
		uMemBlock, _ := memBlockDriver.MustGetBlockWithReadAcquire(uINode, memBlockIndex)
		memBlockData := *uMemBlock.Ptr().BytesSlice()
		assert.Equal(t, memBlockData[0], (byte)(i))
		assert.Equal(t, memBlockData[1], (byte)(i*2))
		uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
	}

	for i = 0; i <= maxBlocks; i++ {
		writeOffset := int64(int64(i) * int64(blockChunkSize))
		memBlockIndex := int(writeOffset / int64(uINode.Ptr().MemBlockSize))
		uMemBlock, _ := memBlockDriver.MustGetBlockWithReadAcquire(uINode, memBlockIndex)
		assert.NoError(t, inodeDriver.FlushMemBlock(uINode, uMemBlock))
		uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
	}

	assert.NoError(t, mockServer.Close())
}
