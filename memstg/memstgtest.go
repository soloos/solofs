package memstg

import (
	"soloos/common/solofsapi"
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
	solonnSRPCServerAddr string,
	memBlockDriver *MemBlockDriver,
	netBlockDriver *NetBlockDriver,
	netINodeDriver *NetINodeDriver,
	blockSize int, blocksLimit int32) {
	var (
		solonnClient solofsapi.SolonnClient
		solodnClient solofsapi.SolodnClient
	)

	NetStgMakeDriversForTest(soloOSEnv,
		solonnSRPCServerAddr,
		&solonnClient, &solodnClient,
		netBlockDriver,
	)

	MemStgMakeMemBlockDriversForTest(memBlockDriver, soloOSEnv, blockSize, blocksLimit)

	util.AssertErrIsNil(netINodeDriver.Init(soloOSEnv, netBlockDriver, memBlockDriver, &solonnClient,
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
		solonnClient solofsapi.SolonnClient
		solodnClient solofsapi.SolodnClient
	)

	NetStgMakeDriversWithMockServerForTest(soloOSEnv,
		mockServerAddr, mockServer,
		&solonnClient, &solodnClient,
		netBlockDriver,
	)

	MemStgMakeMemBlockDriversForTest(memBlockDriver, soloOSEnv, blockSize, blocksLimit)

	util.AssertErrIsNil(netINodeDriver.Init(soloOSEnv, netBlockDriver, memBlockDriver, &solonnClient,
		netINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		netINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		netINodeDriver.NetINodeCommitSizeInDB,
	))
}
