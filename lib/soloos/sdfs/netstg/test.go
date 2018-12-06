package netstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func InitNetBlockDriversForTest(t *testing.T,
	netBlockDriver *NetBlockDriver,
	offheapDriver *offheap.OffheapDriver,
	snetDriver *snet.SNetDriver,
	snetClientDriver *snet.ClientDriver,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
) {
	netBlockDriverOptions := NetBlockDriverOptions{
		int32(-1),
	}
	assert.NoError(t, netBlockDriver.Init(netBlockDriverOptions, offheapDriver,
		snetDriver, snetClientDriver,
		nameNodeClient, dataNodeClient,
	))
}

func InitDriversForTest(t *testing.T,
	snetDriver *snet.SNetDriver,
	snetClientDriver *snet.ClientDriver,
	mockServerAddr string,
	mockServer *MockServer,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
	netBlockDriver *NetBlockDriver,
) {
	var (
		offheapDriver = &offheap.DefaultOffheapDriver
		nameNodePeer  snettypes.PeerUintptr
	)

	assert.NoError(t, snetDriver.Init(offheapDriver))
	assert.NoError(t, snetClientDriver.Init(offheapDriver))

	assert.NoError(t, mockServer.Init(snetDriver, "tcp", mockServerAddr))
	go func() {
		util.AssertErrIsNil(mockServer.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	nameNodePeer = snetDriver.MustGetPeer(nil, mockServerAddr, types.DefaultSDFSRPCProtocol)
	assert.NoError(t, nameNodeClient.Init(snetDriver, snetClientDriver, nameNodePeer))
	assert.NoError(t, dataNodeClient.Init(snetDriver, snetClientDriver))
	InitNetBlockDriversForTest(t, netBlockDriver, offheapDriver,
		snetDriver, snetClientDriver,
		nameNodeClient, dataNodeClient,
	)
}
