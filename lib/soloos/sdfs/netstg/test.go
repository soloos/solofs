package netstg

import (
	snettypes "soloos/common/snet/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/common/util"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"time"
)

func MakeNetBlockDriversForTest(soloOSEnv *soloosbase.SoloOSEnv,
	netBlockDriver *NetBlockDriver,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
) {
	util.AssertErrIsNil(netBlockDriver.Init(soloOSEnv,
		nameNodeClient, dataNodeClient,
		netBlockDriver.PrepareNetBlockMetaDataWithTransfer,
	))
}

func MakeDriversForTest(soloOSEnv *soloosbase.SoloOSEnv,
	nameNodeSRPCServerAddr string,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
	netBlockDriver *NetBlockDriver,
) {
	var (
		nameNodePeer snettypes.PeerUintptr
	)

	nameNodePeer = soloOSEnv.SNetDriver.AllocPeer(nameNodeSRPCServerAddr, types.DefaultSDFSRPCProtocol)
	util.AssertErrIsNil(nameNodeClient.Init(soloOSEnv, nameNodePeer))
	util.AssertErrIsNil(dataNodeClient.Init(soloOSEnv))
	MakeNetBlockDriversForTest(soloOSEnv,
		netBlockDriver, nameNodeClient, dataNodeClient,
	)
}

func MakeMockServerForTest(soloOSEnv *soloosbase.SoloOSEnv,
	mockServerAddr string, mockServer *MockServer) {
	util.AssertErrIsNil(mockServer.Init(soloOSEnv, "tcp", mockServerAddr))
	go func() {
		util.AssertErrIsNil(mockServer.Serve())
	}()
	time.Sleep(time.Millisecond * 300)
}

func MakeDriversWithMockServerForTest(soloOSEnv *soloosbase.SoloOSEnv,
	mockServerAddr string,
	mockServer *MockServer,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
	netBlockDriver *NetBlockDriver,
) {
	MakeDriversForTest(soloOSEnv, mockServerAddr, nameNodeClient, dataNodeClient, netBlockDriver)
	MakeMockServerForTest(soloOSEnv, mockServerAddr, mockServer)
}
