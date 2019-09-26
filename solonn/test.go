package solonn

import (
	"soloos/common/iron"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/solofs/memstg"
	"soloos/solofs/metastg"
)

func MakeSolonnForTest(soloosEnv *soloosbase.SoloosEnv,
	solonn *Solonn, metaStg *metastg.MetaStg,
	solonnSRPCPeerID snettypes.PeerID,
	solonnSRPCServerAddr string,
	solonnWebPeerID snettypes.PeerID,
	solonnWebServerAddr string,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *memstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) {
	var err error

	netBlockDriver.SetHelper(nil, metaStg.PrepareNetBlockMetaData)
	netINodeDriver.SetHelper(nil,
		metaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		metaStg.PrepareNetINodeMetaDataWithStorDB,
		metaStg.NetINodeCommitSizeInDB,
	)

	var webServerOptions = iron.Options{
		ServeStr:  solonnWebServerAddr,
		ListenStr: solonnWebServerAddr,
	}
	err = solonn.Init(soloosEnv,
		solonnSRPCPeerID, solonnSRPCServerAddr, solonnSRPCServerAddr,
		solonnWebPeerID, webServerOptions,
		metaStg,
		memBlockDriver,
		netBlockDriver,
		netINodeDriver,
	)
	util.AssertErrIsNil(err)
}
