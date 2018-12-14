package namenode

import (
	"soloos/sdfs/metastg"
	"soloos/util"
	"soloos/util/offheap"
)

func MakeNameNodeForTest(nameNode *NameNode, metaStg *metastg.MetaStg, nameNodeSRPCServerAddr string) {
	var (
		offheapDriver *offheap.OffheapDriver = &offheap.DefaultOffheapDriver

		options = NameNodeOptions{
			SRPCServer: NameNodeSRPCServerOptions{
				Network:    "tcp",
				ListenAddr: nameNodeSRPCServerAddr,
			},
		}
		err error
	)
	err = nameNode.Init(options, offheapDriver, metaStg)
	util.AssertErrIsNil(err)
}
