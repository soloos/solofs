package memstg

import (
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/snet"
	"soloos/util/offheap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func InitMemBlockDriversForTest(t *testing.T,
	memBlockDriver *MemBlockDriver, offheapDriver *offheap.OffheapDriver, blockChunkSize int) {
	memBlockDriverOptions := MemBlockDriverOptions{
		[]MemBlockPoolOptions{
			MemBlockPoolOptions{
				int32(-1),
				offheap.MakeDefaultTestChunkPoolOptions(blockChunkSize),
			},
		},
	}
	assert.NoError(t, memBlockDriver.Init(memBlockDriverOptions, offheapDriver))
}

func InitNetBlockDriversForTest(t *testing.T,
	netBlockDriver *netstg.NetBlockDriver,
	offheapDriver *offheap.OffheapDriver,
	snetClientDriver *snet.ClientDriver) {
	netBlockDriverOptions := netstg.NetBlockDriverOptions{
		netstg.NetBlockPoolOptions{
			int32(-1),
		},
	}
	assert.NoError(t, netBlockDriver.Init(netBlockDriverOptions, offheapDriver, snetClientDriver))
}

func TestINodeDriverINodeWrite(t *testing.T) {
	var (
		offheapDriver    = &offheap.DefaultOffheapDriver
		snetClientDriver snet.ClientDriver
		netBlockDriver   netstg.NetBlockDriver
		memBlockDriver   MemBlockDriver
		inodeDriver      INodeDriver
		blockChunkSize         = 4
		maxBlocks        int32 = 4
	)

	assert.NoError(t, snetClientDriver.Init(offheapDriver))
	InitMemBlockDriversForTest(t, &memBlockDriver, offheapDriver, blockChunkSize)
	InitNetBlockDriversForTest(t, &netBlockDriver, offheapDriver, &snetClientDriver)

	inodePoolOptions := INodePoolOptions{
		-1}
	assert.NoError(t, inodeDriver.Init(inodePoolOptions, offheapDriver, &netBlockDriver, &memBlockDriver))

	var inodeID types.INodeID
	uINode, _ := inodeDriver.MustGetINode(inodeID)
	uINode.Ptr().NetBlockSize = blockChunkSize
	uINode.Ptr().MemBlockSize = blockChunkSize

	var i int32
	for i = 0; i <= maxBlocks; i++ {
		// write
		writeOffset := int64(int64(i) * int64(blockChunkSize))

		err := inodeDriver.WriteAt(uINode, []byte{(byte)(i), (byte)(i * 2)}, writeOffset)
		assert.NoError(t, err)

		// check
		// netBlockIndex := int(writeOffset / int64(uINode.Ptr().NetBlockSize))
		// uNetBlock, _ := netBlockDriver.MustGetBlock(uINode, netBlockIndex)

		memBlockIndex := int(writeOffset / int64(uINode.Ptr().MemBlockSize))
		uMemBlock, _ := memBlockDriver.MustGetBlockWithReadAcquire(uINode, memBlockIndex)
		memBlockData := *uMemBlock.Ptr().BytesSlice()
		assert.Equal(t, memBlockData[0], (byte)(i))
		assert.Equal(t, memBlockData[1], (byte)(i*2))
		uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
	}
}
