package namenode

import (
	snettypes "soloos/common/snet/types"
	"soloos/common/util"
	"soloos/sdbone/offheap"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
)

func MakeNameNodeForTest(nameNode *NameNode, metaStg *metastg.MetaStg, nameNodeSRPCServerAddr string,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) {
	var (
		offheapDriver *offheap.OffheapDriver = &offheap.DefaultOffheapDriver
		err           error
	)

	netBlockDriver.SetHelper(nil, metaStg.PrepareNetBlockMetaData)
	netINodeDriver.SetHelper(nil,
		metaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		metaStg.PrepareNetINodeMetaDataWithStorDB,
		metaStg.NetINodeCommitSizeInDB,
	)
	var nameNodePeerID snettypes.PeerID
	snettypes.InitTmpPeerID(&nameNodePeerID)
	err = nameNode.Init(offheapDriver, nameNodeSRPCServerAddr, nameNodePeerID, metaStg,
		memBlockDriver,
		netBlockDriver,
		netINodeDriver,
	)
	util.AssertErrIsNil(err)
}
