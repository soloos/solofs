package datanode

import (
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
	"soloos/snet"
	"soloos/util"
	"soloos/util/offheap"
)

func MakeDataNodeForTest(snetDriver *snet.SNetDriver,
	dataNode *DataNode,
	dataNodeSRPCServerAddr string,
	metaStg *metastg.MetaStg,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) {
	var (
		offheapDriver *offheap.OffheapDriver = &offheap.DefaultOffheapDriver

		options = DataNodeOptions{
			SRPCServer: DataNodeSRPCServerOptions{
				Network:    "tcp",
				ListenAddr: dataNodeSRPCServerAddr,
			},
		}
		err error
	)

	netBlockDriver.SetHelper(nil, metaStg.PrepareNetBlockMetaData)
	netINodeDriver.SetHelper(nil,
		metaStg.PrepareNetINodeMetaDataOnlyLoadDB, metaStg.PrepareNetINodeMetaDataWithStorDB)
	err = dataNode.Init(options, offheapDriver, snetDriver, metaStg,
		memBlockDriver, netBlockDriver, netINodeDriver)
	util.AssertErrIsNil(err)
}
