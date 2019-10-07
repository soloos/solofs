package memstg

import (
	"soloos/common/snet"
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
)

func NetStgMakeNetBlockDriversForTest(soloosEnv *soloosbase.SoloosEnv,
	netBlockDriver *NetBlockDriver,
	solonnClient *solofsapi.SolonnClient,
	solodnClient *solofsapi.SolodnClient,
) {
	util.AssertErrIsNil(netBlockDriver.Init(soloosEnv,
		solonnClient, solodnClient,
		// TODO test netBlockDriver.PrepareNetBlockMetaDataWithTransfer,
		netBlockDriver.PrepareNetBlockMetaData,
		nil, nil, nil,
	))
}

func NetStgMakeDriversForTest(soloosEnv *soloosbase.SoloosEnv,
	solonnSrpcServerAddr string,
	solonnClient *solofsapi.SolonnClient,
	solodnClient *solofsapi.SolodnClient,
	netBlockDriver *NetBlockDriver,
) {
	var solonnPeer snet.Peer
	soloosEnv.SNetDriver.InitPeerID((*snet.PeerID)(&solonnPeer.ID))
	solonnPeer.SetAddress(solonnSrpcServerAddr)
	solonnPeer.ServiceProtocol = solofsapitypes.DefaultSolofsRPCProtocol
	util.AssertErrIsNil(soloosEnv.SNetDriver.RegisterPeer(solonnPeer))

	util.AssertErrIsNil(solonnClient.Init(soloosEnv, solonnPeer.ID))
	util.AssertErrIsNil(solodnClient.Init(soloosEnv))
	NetStgMakeNetBlockDriversForTest(soloosEnv,
		netBlockDriver, solonnClient, solodnClient,
	)
}

func NetStgMakeDriversWithMockServerForTest(soloosEnv *soloosbase.SoloosEnv,
	mockServerAddr string,
	mockServer *MockServer,
	solonnClient *solofsapi.SolonnClient,
	solodnClient *solofsapi.SolodnClient,
	netBlockDriver *NetBlockDriver,
) {
	NetStgMakeDriversForTest(soloosEnv, mockServerAddr, solonnClient, solodnClient, netBlockDriver)
	MakeMockServerForTest(soloosEnv, mockServerAddr, mockServer)
}
