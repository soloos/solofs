package memstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
	"soloos/common/snet"
	"soloos/common/util"
	"soloos/sdbone/offheap"
)

func MakeMemBlockDriversForTest(memBlockDriver *MemBlockDriver, offheapDriver *offheap.OffheapDriver,
	blockSize int, blocksLimit int32) {
	memBlockDriverOptions := MemBlockDriverOptions{
		[]MemBlockTableOptions{
			MemBlockTableOptions{
				blockSize, blocksLimit,
			},
		},
	}
	util.AssertErrIsNil(memBlockDriver.Init(offheapDriver, memBlockDriverOptions))
}

func MakeDriversForTest(snetDriver *snet.NetDriver, snetClientDriver *snet.ClientDriver,
	nameNodeSRPCServerAddr string,
	memBlockDriver *MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *NetINodeDriver,
	blockSize int, blocksLimit int32) {
	var (
		offheapDriver  = &offheap.DefaultOffheapDriver
		nameNodeClient api.NameNodeClient
		dataNodeClient api.DataNodeClient
	)

	netstg.MakeDriversForTest(snetDriver, snetClientDriver,
		nameNodeSRPCServerAddr,
		&nameNodeClient, &dataNodeClient,
		netBlockDriver,
	)

	MakeMemBlockDriversForTest(memBlockDriver, offheapDriver, blockSize, blocksLimit)

	util.AssertErrIsNil(netINodeDriver.Init(offheapDriver, netBlockDriver, memBlockDriver, &nameNodeClient,
		netINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		netINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		netINodeDriver.NetINodeCommitSizeInDB,
	))
}

func MakeDriversWithMockServerForTest(mockServerAddr string,
	mockServer *netstg.MockServer,
	snetDriver *snet.NetDriver,
	netBlockDriver *netstg.NetBlockDriver,
	memBlockDriver *MemBlockDriver,
	netINodeDriver *NetINodeDriver,
	blockSize int, blocksLimit int32) {
	var (
		offheapDriver    = &offheap.DefaultOffheapDriver
		snetClientDriver snet.ClientDriver
		nameNodeClient   api.NameNodeClient
		dataNodeClient   api.DataNodeClient
	)

	netstg.MakeDriversWithMockServerForTest(snetDriver, &snetClientDriver,
		mockServerAddr, mockServer,
		&nameNodeClient, &dataNodeClient,
		netBlockDriver,
	)

	MakeMemBlockDriversForTest(memBlockDriver, offheapDriver, blockSize, blocksLimit)

	util.AssertErrIsNil(netINodeDriver.Init(offheapDriver, netBlockDriver, memBlockDriver, &nameNodeClient,
		netINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		netINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		netINodeDriver.NetINodeCommitSizeInDB,
	))
}
