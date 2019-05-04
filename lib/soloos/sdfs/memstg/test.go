package memstg

import (
	soloosbase "soloos/common/soloosapi/base"
	"soloos/common/util"
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
)

func MakeMemBlockDriversForTest(memBlockDriver *MemBlockDriver, soloOSEnv *soloosbase.SoloOSEnv,
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

func MakeDriversForTest(soloOSEnv *soloosbase.SoloOSEnv,
	nameNodeSRPCServerAddr string,
	memBlockDriver *MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *NetINodeDriver,
	blockSize int, blocksLimit int32) {
	var (
		nameNodeClient api.NameNodeClient
		dataNodeClient api.DataNodeClient
	)

	netstg.MakeDriversForTest(soloOSEnv,
		nameNodeSRPCServerAddr,
		&nameNodeClient, &dataNodeClient,
		netBlockDriver,
	)

	MakeMemBlockDriversForTest(memBlockDriver, soloOSEnv, blockSize, blocksLimit)

	util.AssertErrIsNil(netINodeDriver.Init(soloOSEnv, netBlockDriver, memBlockDriver, &nameNodeClient,
		netINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		netINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		netINodeDriver.NetINodeCommitSizeInDB,
	))
}

func MakeDriversWithMockServerForTest(soloOSEnv *soloosbase.SoloOSEnv,
	mockServerAddr string,
	mockServer *netstg.MockServer,
	netBlockDriver *netstg.NetBlockDriver,
	memBlockDriver *MemBlockDriver,
	netINodeDriver *NetINodeDriver,
	blockSize int, blocksLimit int32) {
	var (
		nameNodeClient api.NameNodeClient
		dataNodeClient api.DataNodeClient
	)

	netstg.MakeDriversWithMockServerForTest(soloOSEnv,
		mockServerAddr, mockServer,
		&nameNodeClient, &dataNodeClient,
		netBlockDriver,
	)

	MakeMemBlockDriversForTest(memBlockDriver, soloOSEnv, blockSize, blocksLimit)

	util.AssertErrIsNil(netINodeDriver.Init(soloOSEnv, netBlockDriver, memBlockDriver, &nameNodeClient,
		netINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		netINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		netINodeDriver.NetINodeCommitSizeInDB,
	))
}
