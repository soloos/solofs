package namenode

import (
	snettypes "soloos/common/snet/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/common/util"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
)

func MakeNameNodeForTest(soloOSEnv *soloosbase.SoloOSEnv,
	nameNode *NameNode, metaStg *metastg.MetaStg, nameNodeSRPCServerAddr string,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) {
	var err error

	netBlockDriver.SetHelper(nil, metaStg.PrepareNetBlockMetaData)
	netINodeDriver.SetHelper(nil,
		metaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		metaStg.PrepareNetINodeMetaDataWithStorDB,
		metaStg.NetINodeCommitSizeInDB,
	)
	var nameNodePeerID snettypes.PeerID
	snettypes.InitTmpPeerID(&nameNodePeerID)
	err = nameNode.Init(soloOSEnv, nameNodeSRPCServerAddr, nameNodePeerID, metaStg,
		memBlockDriver,
		netBlockDriver,
		netINodeDriver,
	)
	util.AssertErrIsNil(err)
}
