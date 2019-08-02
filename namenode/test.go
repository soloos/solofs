package namenode

import (
	"soloos/common/iron"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
)

func MakeNameNodeForTest(soloOSEnv *soloosbase.SoloOSEnv,
	nameNode *NameNode, metaStg *metastg.MetaStg,
	nameNodeSRPCPeerID snettypes.PeerID,
	nameNodeSRPCServerAddr string,
	nameNodeWebPeerID snettypes.PeerID,
	nameNodeWebServerAddr string,
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
		ServeStr:  nameNodeWebServerAddr,
		ListenStr: nameNodeWebServerAddr,
	}
	err = nameNode.Init(soloOSEnv,
		nameNodeSRPCPeerID, nameNodeSRPCServerAddr, nameNodeSRPCServerAddr,
		nameNodeWebPeerID, webServerOptions,
		metaStg,
		memBlockDriver,
		netBlockDriver,
		netINodeDriver,
	)
	util.AssertErrIsNil(err)
}
