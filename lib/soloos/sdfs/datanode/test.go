package datanode

import (
	"path/filepath"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"
)

func MakeDataNodeForTest(snetDriver *snet.NetDriver, snetClientDriver *snet.ClientDriver,
	dataNode *DataNode,
	dataNodeSRPCServerAddr string,
	metaStg *metastg.MetaStg,
	nameNodeSRPCServerAddr string,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) {
	var (
		offheapDriver *offheap.OffheapDriver = &offheap.DefaultOffheapDriver
		err           error
	)

	var peerID snettypes.PeerID
	util.InitUUID64(&peerID)
	var localFsRoot = filepath.Join("/tmp/sdfs_test.data", string(peerID[:3]))

	var options = DataNodeOptions{
		PeerID:               peerID,
		SrpcServerListenAddr: dataNodeSRPCServerAddr,
		SrpcServerServeAddr:  dataNodeSRPCServerAddr,
		LocalFsRoot:          localFsRoot,
		NameNodeSRPCServer:   nameNodeSRPCServerAddr,
	}

	err = dataNode.Init(offheapDriver, options,
		snetDriver, snetClientDriver, metaStg,
		memBlockDriver, netBlockDriver, netINodeDriver)
	util.AssertErrIsNil(err)
}
