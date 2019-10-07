package solonn

import (
	"soloos/common/iron"
	"soloos/common/snet"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/solofs/memstg"
	"soloos/solofs/metastg"
)

func MakeSolonnForTest(soloosEnv *soloosbase.SoloosEnv,
	solonn *Solonn, metaStg *metastg.MetaStg,
	solonnSrpcPeerID snet.PeerID,
	solonnSrpcServerAddr string,
	solonnWebPeerID snet.PeerID,
	solonnWebServerAddr string,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *memstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) {
	var err error

	netBlockDriver.SetHelper(metaStg.PrepareNetBlockMetaData, nil, nil, nil)
	netINodeDriver.SetHelper(
		metaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		metaStg.PrepareNetINodeMetaDataWithStorDB,
		metaStg.NetINodeCommitSizeInDB,
	)

	var webServerOptions = iron.Options{
		ServeStr:  solonnWebServerAddr,
		ListenStr: solonnWebServerAddr,
	}
	err = solonn.Init(soloosEnv,
		solonnSrpcPeerID, solonnSrpcServerAddr, solonnSrpcServerAddr,
		solonnWebPeerID, webServerOptions,
		metaStg,
		memBlockDriver,
		netBlockDriver,
		netINodeDriver,
	)
	util.AssertErrIsNil(err)
}
