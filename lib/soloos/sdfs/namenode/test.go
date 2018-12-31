package namenode

import (
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
	"soloos/util"
	"soloos/util/offheap"
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
		metaStg.PrepareNetINodeMetaDataOnlyLoadDB, metaStg.PrepareNetINodeMetaDataWithStorDB)
	err = nameNode.Init(offheapDriver, nameNodeSRPCServerAddr, metaStg,
		memBlockDriver,
		netBlockDriver,
		netINodeDriver,
	)
	util.AssertErrIsNil(err)
}
