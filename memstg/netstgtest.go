package memstg

import (
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/snettypes"
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
	))
}

func NetStgMakeDriversForTest(soloosEnv *soloosbase.SoloosEnv,
	solonnSRPCServerAddr string,
	solonnClient *solofsapi.SolonnClient,
	solodnClient *solofsapi.SolodnClient,
	netBlockDriver *NetBlockDriver,
) {
	var solonnPeer snettypes.Peer
	soloosEnv.SNetDriver.InitPeerID((*snettypes.PeerID)(&solonnPeer.ID))
	solonnPeer.SetAddress(solonnSRPCServerAddr)
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
