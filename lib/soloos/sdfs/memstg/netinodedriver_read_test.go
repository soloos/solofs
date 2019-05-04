package memstg

import (
	sdfsapitypes "soloos/common/sdfsapi/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/common/util"
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetINodeDriverNetINodeRead(t *testing.T) {
	var (
		soloOSEnv      soloosbase.SoloOSEnv
		mockServer     netstg.MockServer
		netBlockDriver netstg.NetBlockDriver
		memBlockDriver MemBlockDriver
		netINodeDriver NetINodeDriver
		netBlockCap    int   = 128
		memBlockCap    int   = 64
		blocksLimit    int32 = 4
		uNetINode      types.NetINodeUintptr
		err            error
	)
	util.AssertErrIsNil(soloOSEnv.Init())
	MakeDriversWithMockServerForTest(&soloOSEnv, "127.0.0.1:10022", &mockServer,
		&netBlockDriver, &memBlockDriver, &netINodeDriver, memBlockCap, blocksLimit)
	var netINodeID types.NetINodeID
	sdfsapitypes.InitTmpNetINodeID(&netINodeID)
	uNetINode, err = netINodeDriver.MustGetNetINodeWithReadAcquire(netINodeID, 0, netBlockCap, memBlockCap)
	assert.NoError(t, err)

	var (
		readData        = make([]byte, 93)
		readOff  uint64 = 73
	)

	err = netINodeDriver.PWriteWithMem(uNetINode, readData, readOff)
	assert.NoError(t, err)

	_, err = netINodeDriver.PReadWithMem(uNetINode, readData, readOff)
	assert.NoError(t, err)

	assert.NoError(t, mockServer.Close())
}
