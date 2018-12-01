package memstg

import (
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestINodeDriverINodeRead(t *testing.T) {
	var (
		mockServer     netstg.MockServer
		memBlockDriver MemBlockDriver
		inodeDriver    INodeDriver
		netBlockSize   int = 128
		memBlockSize   int = 64
		uINode         types.INodeUintptr
	)
	InitDriversForTest(t, "127.0.0.1:10022", &mockServer, &memBlockDriver, &inodeDriver, memBlockSize)
	uINode = InitInodeForTest(t, &inodeDriver, netBlockSize, memBlockSize)

	var (
		readData       = make([]byte, 93)
		readOff  int64 = 73
	)
	assert.NoError(t, inodeDriver.PRead(uINode, readData, readOff))
	assert.NoError(t, inodeDriver.PRead(uINode, readData, readOff))

	assert.NoError(t, mockServer.Close())
}
