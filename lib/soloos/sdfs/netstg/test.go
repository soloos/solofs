package netstg

import (
	"soloos/common/sdfsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"time"
)

func MakeNetBlockDriversForTest(soloOSEnv *soloosbase.SoloOSEnv,
	netBlockDriver *NetBlockDriver,
	nameNodeClient *sdfsapi.NameNodeClient,
	dataNodeClient *sdfsapi.DataNodeClient,
) {
	util.AssertErrIsNil(netBlockDriver.Init(soloOSEnv,
		nameNodeClient, dataNodeClient,
		// TODO test netBlockDriver.PrepareNetBlockMetaDataWithTransfer,
		netBlockDriver.PrepareNetBlockMetaData,
	))
}

func MakeDriversForTest(soloOSEnv *soloosbase.SoloOSEnv,
	nameNodeSRPCServerAddr string,
	nameNodeClient *sdfsapi.NameNodeClient,
	dataNodeClient *sdfsapi.DataNodeClient,
	netBlockDriver *NetBlockDriver,
) {
	var (
		nameNodePeer snettypes.PeerUintptr
	)

	nameNodePeer, _ = soloOSEnv.SNetDriver.MustGetPeer(nil, nameNodeSRPCServerAddr, sdfsapitypes.DefaultSDFSRPCProtocol)
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
	nameNodeClient *sdfsapi.NameNodeClient,
	dataNodeClient *sdfsapi.DataNodeClient,
	netBlockDriver *NetBlockDriver,
) {
	MakeDriversForTest(soloOSEnv, mockServerAddr, nameNodeClient, dataNodeClient, netBlockDriver)
	MakeMockServerForTest(soloOSEnv, mockServerAddr, mockServer)
}
