package datanode

import (
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
	"soloos/util"
	"soloos/util/offheap"
)

func MakeDataNodeForTest(dataNode *DataNode,
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
	err = dataNode.Init(options, offheapDriver, metaStg, netBlockDriver, memBlockDriver, netINodeDriver)
	util.AssertErrIsNil(err)
}
