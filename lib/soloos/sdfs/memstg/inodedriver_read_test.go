package memstg

import (
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestINodeDriverINodeRead(t *testing.T) {
	var (
		mockServer       netstg.MockServer
		memBlockDriver   MemBlockDriver
		inodeDriver      INodeDriver
		netBlockCap      int   = 128
		memBlockCap      int   = 64
		blockChunksLimit int32 = 4
		uINode           types.INodeUintptr
	)
	InitDriversForTest(t, "127.0.0.1:10022", &mockServer, &memBlockDriver, &inodeDriver, memBlockCap, blockChunksLimit)
	uINode = inodeDriver.InitINode(netBlockCap, memBlockCap)

	var (
		readData       = make([]byte, 93)
		readOff  int64 = 73
	)
	assert.NoError(t, inodeDriver.PRead(uINode, readData, readOff))
	assert.NoError(t, inodeDriver.PRead(uINode, readData, readOff))

	assert.NoError(t, mockServer.Close())
}
