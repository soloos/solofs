package memstg

import (
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/snet"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetINodeDriverNetINodeRead(t *testing.T) {
	var (
		mockServer       netstg.MockServer
		snetDriver       snet.SNetDriver
		netBlockDriver   netstg.NetBlockDriver
		memBlockDriver   MemBlockDriver
		netINodeDriver   NetINodeDriver
		netBlockCap      int   = 128
		memBlockCap      int   = 64
		blockChunksLimit int32 = 4
		uNetINode        types.NetINodeUintptr
		err              error
	)
	MakeDriversWithMockServerForTest(t,
		"127.0.0.1:10022", &mockServer, &snetDriver,
		&netBlockDriver, &memBlockDriver, &netINodeDriver, memBlockCap, blockChunksLimit)
	uNetINode, err = netINodeDriver.AllocNetINode(0, netBlockCap, memBlockCap)
	assert.NoError(t, err)

	var (
		readData       = make([]byte, 93)
		readOff  int64 = 73
	)
	assert.NoError(t, netINodeDriver.PRead(uNetINode, readData, readOff))
	assert.NoError(t, netINodeDriver.PRead(uNetINode, readData, readOff))

	assert.NoError(t, mockServer.Close())
}
