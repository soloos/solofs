package minicluster

import (
	"fmt"
	"soloos/sdfs/datanode"
	"soloos/sdfs/metastg"
	"soloos/sdfs/namenode"
	"soloos/util"
)

type MiniCluster struct {
	DataNodes []datanode.DataNode
	NameNodes []namenode.NameNode
}

func (p *MiniCluster) Init(nameNodePorts []int, dataNodePorts []int) {
	p.NameNodes = make([]namenode.NameNode, len(nameNodePorts))
	for i := 0; i < len(nameNodePorts); i++ {
		nameNodePort := nameNodePorts[i]
		options := namenode.NameNodeOptions{
			namenode.NameNodeSRPCServerOptions{
				"tcp",
				fmt.Sprintf("127.0.0.1:%d", nameNodePort),
			},
			metastg.TestMetaStgDBDriver,
			metastg.TestMetaStgDBConnect,
		}
		util.AssertErrIsNil(p.NameNodes[i].Init(options))
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
		util.AssertErrIsNil(p.DataNodes[i].Init(options))
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
