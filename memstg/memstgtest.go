package memstg

import (
	"soloos/common/sdfsapi"
	"soloos/common/soloosbase"
	"soloos/common/util"
)

func MemStgMakeMemBlockDriversForTest(memBlockDriver *MemBlockDriver, soloOSEnv *soloosbase.SoloOSEnv,
	blockSize int, blocksLimit int32) {
	memBlockDriverOptions := MemBlockDriverOptions{
		[]MemBlockTableOptions{
			MemBlockTableOptions{
				blockSize, blocksLimit,
			},
		},
	}
	util.AssertErrIsNil(memBlockDriver.Init(soloOSEnv, memBlockDriverOptions))
}

func MemStgMakeDriversForTest(soloOSEnv *soloosbase.SoloOSEnv,
	nameNodeSRPCServerAddr string,
	memBlockDriver *MemBlockDriver,
	netBlockDriver *NetBlockDriver,
	netINodeDriver *NetINodeDriver,
	blockSize int, blocksLimit int32) {
	var (
		nameNodeClient sdfsapi.NameNodeClient
		dataNodeClient sdfsapi.DataNodeClient
	)

	NetStgMakeDriversForTest(soloOSEnv,
		nameNodeSRPCServerAddr,
		&nameNodeClient, &dataNodeClient,
		netBlockDriver,
	)

	MemStgMakeMemBlockDriversForTest(memBlockDriver, soloOSEnv, blockSize, blocksLimit)

	util.AssertErrIsNil(netINodeDriver.Init(soloOSEnv, netBlockDriver, memBlockDriver, &nameNodeClient,
		netINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		netINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		netINodeDriver.NetINodeCommitSizeInDB,
	))
}

func MemStgMakeDriversWithMockServerForTest(soloOSEnv *soloosbase.SoloOSEnv,
	mockServerAddr string,
	mockServer *MockServer,
	netBlockDriver *NetBlockDriver,
	memBlockDriver *MemBlockDriver,
	netINodeDriver *NetINodeDriver,
	blockSize int, blocksLimit int32) {
	var (
		nameNodeClient sdfsapi.NameNodeClient
		dataNodeClient sdfsapi.DataNodeClient
	)

	NetStgMakeDriversWithMockServerForTest(soloOSEnv,
		mockServerAddr, mockServer,
		&nameNodeClient, &dataNodeClient,
		netBlockDriver,
	)

	MemStgMakeMemBlockDriversForTest(memBlockDriver, soloOSEnv, blockSize, blocksLimit)

	util.AssertErrIsNil(netINodeDriver.Init(soloOSEnv, netBlockDriver, memBlockDriver, &nameNodeClient,
		netINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		netINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		netINodeDriver.NetINodeCommitSizeInDB,
	))
}
