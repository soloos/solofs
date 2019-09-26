package memstg

import (
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
)

func NetStgMakeNetBlockDriversForTest(soloOSEnv *soloosbase.SoloOSEnv,
	netBlockDriver *NetBlockDriver,
	solonnClient *solofsapi.SolonnClient,
	solodnClient *solofsapi.SolodnClient,
) {
	util.AssertErrIsNil(netBlockDriver.Init(soloOSEnv,
		solonnClient, solodnClient,
		// TODO test netBlockDriver.PrepareNetBlockMetaDataWithTransfer,
		netBlockDriver.PrepareNetBlockMetaData,
	))
}

func NetStgMakeDriversForTest(soloOSEnv *soloosbase.SoloOSEnv,
	solonnSRPCServerAddr string,
	solonnClient *solofsapi.SolonnClient,
	solodnClient *solofsapi.SolodnClient,
	netBlockDriver *NetBlockDriver,
) {
	var solonnPeer snettypes.Peer
	soloOSEnv.SNetDriver.InitPeerID((*snettypes.PeerID)(&solonnPeer.ID))
	solonnPeer.SetAddress(solonnSRPCServerAddr)
	solonnPeer.ServiceProtocol = solofsapitypes.DefaultSOLOFSRPCProtocol
	util.AssertErrIsNil(soloOSEnv.SNetDriver.RegisterPeer(solonnPeer))

	util.AssertErrIsNil(solonnClient.Init(soloOSEnv, solonnPeer.ID))
	util.AssertErrIsNil(solodnClient.Init(soloOSEnv))
	NetStgMakeNetBlockDriversForTest(soloOSEnv,
		netBlockDriver, solonnClient, solodnClient,
	)
}

func NetStgMakeDriversWithMockServerForTest(soloOSEnv *soloosbase.SoloOSEnv,
	mockServerAddr string,
	mockServer *MockServer,
	solonnClient *solofsapi.SolonnClient,
	solodnClient *solofsapi.SolodnClient,
	netBlockDriver *NetBlockDriver,
) {
	NetStgMakeDriversForTest(soloOSEnv, mockServerAddr, solonnClient, solodnClient, netBlockDriver)
	MakeMockServerForTest(soloOSEnv, mockServerAddr, mockServer)
}
