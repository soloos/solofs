package memstg

import (
	"soloos/common/solofsapi"
	"soloos/common/soloosbase"
	"soloos/common/util"
)

func MemStgMakeMemBlockDriversForTest(memBlockDriver *MemBlockDriver, soloosEnv *soloosbase.SoloosEnv,
	blockSize int, blocksLimit int32) {
	memBlockDriverOptions := MemBlockDriverOptions{
		[]MemBlockTableOptions{
			MemBlockTableOptions{
				blockSize, blocksLimit,
			},
		},
	}
	util.AssertErrIsNil(memBlockDriver.Init(soloosEnv, memBlockDriverOptions))
}

func MemStgMakeDriversForTest(soloosEnv *soloosbase.SoloosEnv,
	solonnSrpcServerAddr string,
	memBlockDriver *MemBlockDriver,
	netBlockDriver *NetBlockDriver,
	netINodeDriver *NetINodeDriver,
	blockSize int, blocksLimit int32) {
	var (
		solonnClient solofsapi.SolonnClient
		solodnClient solofsapi.SolodnClient
	)

	NetStgMakeDriversForTest(soloosEnv,
		solonnSrpcServerAddr,
		&solonnClient, &solodnClient,
		netBlockDriver,
	)

	MemStgMakeMemBlockDriversForTest(memBlockDriver, soloosEnv, blockSize, blocksLimit)

	util.AssertErrIsNil(netINodeDriver.Init(soloosEnv, netBlockDriver, memBlockDriver, &solonnClient,
		netINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		netINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		netINodeDriver.NetINodeCommitSizeInDB,
	))
}

func MemStgMakeDriversWithMockServerForTest(soloosEnv *soloosbase.SoloosEnv,
	mockServerAddr string,
	mockServer *MockServer,
	netBlockDriver *NetBlockDriver,
	memBlockDriver *MemBlockDriver,
	netINodeDriver *NetINodeDriver,
	blockSize int, blocksLimit int32) {
	var (
		solonnClient solofsapi.SolonnClient
		solodnClient solofsapi.SolodnClient
	)

	NetStgMakeDriversWithMockServerForTest(soloosEnv,
		mockServerAddr, mockServer,
		&solonnClient, &solodnClient,
		netBlockDriver,
	)

	MemStgMakeMemBlockDriversForTest(memBlockDriver, soloosEnv, blockSize, blocksLimit)

	util.AssertErrIsNil(netINodeDriver.Init(soloosEnv, netBlockDriver, memBlockDriver, &solonnClient,
		netINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		netINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		netINodeDriver.NetINodeCommitSizeInDB,
	))
}
