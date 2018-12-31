package memstg

import (
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/snet"
	"soloos/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetINodeDriverNetINodeRead(t *testing.T) {
	var (
		mockServer       netstg.MockServer
		snetDriver       snet.NetDriver
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
	var netINodeID types.NetINodeID
	util.InitUUID64(&netINodeID)
	uNetINode, err = netINodeDriver.MustGetNetINode(netINodeID, 0, netBlockCap, memBlockCap)
	assert.NoError(t, err)

	var (
		readData       = make([]byte, 93)
		readOff  int64 = 73
	)
	assert.NoError(t, netINodeDriver.PReadWithMem(uNetINode, readData, readOff))
	assert.NoError(t, netINodeDriver.PReadWithMem(uNetINode, readData, readOff))

	assert.NoError(t, mockServer.Close())
}
