package memstg

import (
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetINodeDriverNetINodeRead(t *testing.T) {
	var (
		soloosEnv      soloosbase.SoloosEnv
		mockServer     MockServer
		netBlockDriver NetBlockDriver
		memBlockDriver MemBlockDriver
		netINodeDriver NetINodeDriver
		netBlockCap    int   = 128
		memBlockCap    int   = 64
		blocksLimit    int32 = 4
		uNetINode      solofstypes.NetINodeUintptr
		err            error
	)
	util.AssertErrIsNil(soloosEnv.InitWithSNet(""))
	MemStgMakeDriversWithMockServerForTest(&soloosEnv, "127.0.0.1:10022", &mockServer,
		&netBlockDriver, &memBlockDriver, &netINodeDriver, memBlockCap, blocksLimit)
	var netINodeID solofstypes.NetINodeID
	solofstypes.InitTmpNetINodeID(&netINodeID)
	uNetINode, err = netINodeDriver.MustGetNetINode(netINodeID, 0, netBlockCap, memBlockCap)
	defer netINodeDriver.ReleaseNetINode(uNetINode)
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
