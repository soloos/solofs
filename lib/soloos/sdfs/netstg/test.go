package netstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/common/snet"
	snettypes "soloos/common/snet/types"
	"soloos/common/util"
	"soloos/sdbone/offheap"
	"time"
)

func MakeNetBlockDriversForTest(netBlockDriver *NetBlockDriver,
	offheapDriver *offheap.OffheapDriver,
	snetDriver *snet.NetDriver,
	snetClientDriver *snet.ClientDriver,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
) {
	util.AssertErrIsNil(netBlockDriver.Init(offheapDriver,
		snetDriver, snetClientDriver,
		nameNodeClient, dataNodeClient,
		netBlockDriver.PrepareNetBlockMetaDataWithTransfer,
	))
}

func MakeDriversForTest(snetDriver *snet.NetDriver, snetClientDriver *snet.ClientDriver,
	nameNodeSRPCServerAddr string,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
	netBlockDriver *NetBlockDriver,
) {
	var (
		offheapDriver = &offheap.DefaultOffheapDriver
		nameNodePeer  snettypes.PeerUintptr
	)

	util.AssertErrIsNil(snetDriver.Init(offheapDriver))
	util.AssertErrIsNil(snetClientDriver.Init(offheapDriver))

	nameNodePeer, _ = snetDriver.MustGetPeer(nil, nameNodeSRPCServerAddr, types.DefaultSDFSRPCProtocol)
	util.AssertErrIsNil(nameNodeClient.Init(snetClientDriver, nameNodePeer))
	util.AssertErrIsNil(dataNodeClient.Init(snetClientDriver))
	MakeNetBlockDriversForTest(netBlockDriver, offheapDriver,
		snetDriver, snetClientDriver,
		nameNodeClient, dataNodeClient,
	)
}

func MakeMockServerForTest(snetDriver *snet.NetDriver,
	mockServerAddr string, mockServer *MockServer) {
	util.AssertErrIsNil(mockServer.Init(snetDriver, "tcp", mockServerAddr))
	go func() {
		util.AssertErrIsNil(mockServer.Serve())
	}()
	time.Sleep(time.Millisecond * 300)
}

func MakeDriversWithMockServerForTest(snetDriver *snet.NetDriver, snetClientDriver *snet.ClientDriver,
	mockServerAddr string,
	mockServer *MockServer,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
	netBlockDriver *NetBlockDriver,
) {
	MakeDriversForTest(snetDriver, snetClientDriver, mockServerAddr, nameNodeClient, dataNodeClient, netBlockDriver)
	MakeMockServerForTest(snetDriver, mockServerAddr, mockServer)
}
