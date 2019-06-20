package memstg

import (
	"soloos/common/sdfsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
)

func NetStgMakeNetBlockDriversForTest(soloOSEnv *soloosbase.SoloOSEnv,
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

func NetStgMakeDriversForTest(soloOSEnv *soloosbase.SoloOSEnv,
	nameNodeSRPCServerAddr string,
	nameNodeClient *sdfsapi.NameNodeClient,
	dataNodeClient *sdfsapi.DataNodeClient,
	netBlockDriver *NetBlockDriver,
) {
	var nameNodePeer snettypes.Peer
	soloOSEnv.SNetDriver.InitPeerID((*snettypes.PeerID)(&nameNodePeer.ID))
	nameNodePeer.SetAddress(nameNodeSRPCServerAddr)
	nameNodePeer.ServiceProtocol = sdfsapitypes.DefaultSDFSRPCProtocol
	util.AssertErrIsNil(soloOSEnv.SNetDriver.RegisterPeer(nameNodePeer))

	util.AssertErrIsNil(nameNodeClient.Init(soloOSEnv, nameNodePeer.ID))
	util.AssertErrIsNil(dataNodeClient.Init(soloOSEnv))
	NetStgMakeNetBlockDriversForTest(soloOSEnv,
		netBlockDriver, nameNodeClient, dataNodeClient,
	)
}

func NetStgMakeDriversWithMockServerForTest(soloOSEnv *soloosbase.SoloOSEnv,
	mockServerAddr string,
	mockServer *MockServer,
	nameNodeClient *sdfsapi.NameNodeClient,
	dataNodeClient *sdfsapi.DataNodeClient,
	netBlockDriver *NetBlockDriver,
) {
	NetStgMakeDriversForTest(soloOSEnv, mockServerAddr, nameNodeClient, dataNodeClient, netBlockDriver)
	MakeMockServerForTest(soloOSEnv, mockServerAddr, mockServer)
}
