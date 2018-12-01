package memstg

import (
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestINodeDriverINodeWrite(t *testing.T) {
	var (
		mockServer     netstg.MockServer
		memBlockDriver MemBlockDriver
		inodeDriver    INodeDriver
		maxBlocks      int32 = 4
		i              int32
		netBlockSize   int = 4
		memBlockSize   int = 4
		uINode         types.INodeUintptr
	)
	InitDriversForTest(t, "127.0.0.1:10023", &mockServer, &memBlockDriver, &inodeDriver, memBlockSize)
	uINode = InitInodeForTest(t, &inodeDriver, netBlockSize, memBlockSize)

	for i = 0; i <= maxBlocks; i++ {
		writeOffset := int64(int64(i) * int64(memBlockSize))

		assert.NoError(t, inodeDriver.PWrite(uINode, []byte{(byte)(i), (byte)(i * 2)}, writeOffset))

		memBlockIndex := int(writeOffset / int64(uINode.Ptr().MemBlockSize))
		uMemBlock, _ := memBlockDriver.MustGetBlockWithReadAcquire(uINode, memBlockIndex)
		memBlockData := *uMemBlock.Ptr().BytesSlice()
		assert.Equal(t, memBlockData[0], (byte)(i))
		assert.Equal(t, memBlockData[1], (byte)(i*2))
		uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
	}

	for i = 0; i <= maxBlocks; i++ {
		writeOffset := int64(int64(i) * int64(memBlockSize))
		memBlockIndex := int(writeOffset / int64(uINode.Ptr().MemBlockSize))
		uMemBlock, _ := memBlockDriver.MustGetBlockWithReadAcquire(uINode, memBlockIndex)
		assert.NoError(t, inodeDriver.FlushMemBlock(uINode, uMemBlock))
		uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
	}

	assert.NoError(t, mockServer.Close())
}
