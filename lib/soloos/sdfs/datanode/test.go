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
	metaStg *metastg.MetaStg,
	netBlockDriver *netstg.NetBlockDriver,
	memBlockDriver *memstg.MemBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
	dataNodeSRPCServerAddr string,
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
	err = dataNode.Init(options, offheapDriver, snetDriver, metaStg, netBlockDriver, memBlockDriver)
	util.AssertErrIsNil(err)
}
