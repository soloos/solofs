package minicluster

import (
	"fmt"
	"soloos/sdfs/datanode"
	"soloos/sdfs/metastg"
	"soloos/sdfs/namenode"
	"soloos/util"
	"soloos/util/offheap"
)

type MiniCluster struct {
	offheapDriver    *offheap.OffheapDriver
	DataNodes        []datanode.DataNode
	DataNodesMetaStg []metastg.MetaStg
	NameNodes        []namenode.NameNode
	NameNodesMetaStg []metastg.MetaStg
}

func (p *MiniCluster) Init(nameNodePorts []int, dataNodePorts []int) {
	p.offheapDriver = &offheap.DefaultOffheapDriver
	p.NameNodes = make([]namenode.NameNode, len(nameNodePorts))
	p.NameNodesMetaStg = make([]metastg.MetaStg, len(nameNodePorts))
	for i := 0; i < len(nameNodePorts); i++ {
		nameNodePort := nameNodePorts[i]
		util.AssertErrIsNil(p.NameNodesMetaStg[i].Init(p.offheapDriver,
			metastg.TestMetaStgDBDriver,
			metastg.TestMetaStgDBConnect,
		))
		nameNodeOptions := namenode.NameNodeOptions{
			namenode.NameNodeSRPCServerOptions{
				"tcp",
				fmt.Sprintf("127.0.0.1:%d", nameNodePort),
			},
		}
		util.AssertErrIsNil(p.NameNodes[i].Init(nameNodeOptions, p.offheapDriver, &p.NameNodesMetaStg[i]))
		go func() {
			util.AssertErrIsNil(p.NameNodes[i].Serve())
		}()
	}

	p.DataNodes = make([]datanode.DataNode, len(dataNodePorts))
	for i := 0; i < len(dataNodePorts); i++ {
		dataNodePort := dataNodePorts[i]
		options := datanode.DataNodeOptions{
			datanode.DataNodeSRPCServerOptions{
				"tcp",
				fmt.Sprintf("127.0.0.1:%d", dataNodePort),
			},
		}
		util.AssertErrIsNil(p.DataNodes[i].Init(options, p.offheapDriver))
		go func() {
			util.AssertErrIsNil(p.DataNodes[i].Serve())
		}()
	}
}

func (p *MiniCluster) Shutdown() {
	for i := 0; i < len(p.NameNodes); i++ {
		util.AssertErrIsNil(p.NameNodes[i].Close())
	}
	for i := 0; i < len(p.DataNodes); i++ {
		util.AssertErrIsNil(p.DataNodes[i].Close())
	}
}
