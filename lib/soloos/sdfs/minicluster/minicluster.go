package minicluster

import (
	"fmt"
	"soloos/sdfs/datanode"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/namenode"
	"soloos/sdfs/netstg"
	"soloos/common/util"
	"soloos/sdbone/offheap"
)

type MiniCluster struct {
	offheapDriver *offheap.OffheapDriver

	NameNodes               []namenode.NameNode
	NameNodeMetaStgs        []metastg.MetaStg
	NameNodeMemBlockDrivers []memstg.MemBlockDriver
	NameNodeNetBlockDrivers []netstg.NetBlockDriver
	NameNodeNetINodeDrivers []memstg.NetINodeDriver

	DataNodes               []datanode.DataNode
	DataNodeMetaStgs        []metastg.MetaStg
	DataNodeMemBlockDrivers []memstg.MemBlockDriver
	DataNodeNetBlockDrivers []netstg.NetBlockDriver
	DataNodeNetINodeDrivers []memstg.NetINodeDriver
}

func (p *MiniCluster) Init(nameNodePorts []int, dataNodePorts []int) {
	p.offheapDriver = &offheap.DefaultOffheapDriver
	p.NameNodes = make([]namenode.NameNode, len(nameNodePorts))
	p.NameNodeMetaStgs = make([]metastg.MetaStg, len(nameNodePorts))
	p.NameNodeMemBlockDrivers = make([]memstg.MemBlockDriver, len(nameNodePorts))
	p.NameNodeNetBlockDrivers = make([]netstg.NetBlockDriver, len(nameNodePorts))
	p.NameNodeNetINodeDrivers = make([]memstg.NetINodeDriver, len(nameNodePorts))
	for i := 0; i < len(nameNodePorts); i++ {
		metastg.MakeMetaStgForTest(p.offheapDriver, &p.NameNodeMetaStgs[i])

		namenode.MakeNameNodeForTest(&p.NameNodes[i], &p.NameNodeMetaStgs[i],
			fmt.Sprintf("127.0.0.1:%d", nameNodePorts[i]),
			&p.NameNodeMemBlockDrivers[i],
			&p.NameNodeNetBlockDrivers[i],
			&p.NameNodeNetINodeDrivers[i],
		)
		go func() {
			util.AssertErrIsNil(p.NameNodes[i].Serve())
		}()
	}

	p.DataNodes = make([]datanode.DataNode, len(dataNodePorts))
	p.DataNodeMetaStgs = make([]metastg.MetaStg, len(dataNodePorts))
	p.DataNodeMemBlockDrivers = make([]memstg.MemBlockDriver, len(dataNodePorts))
	p.DataNodeNetBlockDrivers = make([]netstg.NetBlockDriver, len(dataNodePorts))
	p.DataNodeNetINodeDrivers = make([]memstg.NetINodeDriver, len(dataNodePorts))
	for i := 0; i < len(dataNodePorts); i++ {
		metastg.MakeMetaStgForTest(p.offheapDriver, &p.NameNodeMetaStgs[i])

		// namenode.MakeNameNodeForTest(&p.DataNodes[i], &p.DataNodeMetaStgs[i],
		// fmt.Sprintf("127.0.0.1:%d", dataNodePorts[i]))
		// util.AssertErrIsNil(p.DataNodes[i].Init(options, p.offheapDriver))
		// go func() {
		// util.AssertErrIsNil(p.DataNodes[i].Serve())
		// }()
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
