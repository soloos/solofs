package datanode

import (
	"path/filepath"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/sdfs/memstg"
)

func MakeDataNodeForTest(soloOSEnv *soloosbase.SoloOSEnv,
	dataNode *DataNode,
	dataNodePeerID snettypes.PeerID, dataNodeSRPCServerAddr string,
	nameNodePeerID snettypes.PeerID, nameNodeSRPCServerAddr string,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *memstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) {
	var (
		err error
	)

	var localFSRoot = filepath.Join("/tmp/sdfs_test.data", dataNodePeerID.Str())

	var options = DataNodeOptions{
		PeerID:               dataNodePeerID,
		SrpcServerListenAddr: dataNodeSRPCServerAddr,
		SrpcServerServeAddr:  dataNodeSRPCServerAddr,
		LocalFSRoot:          localFSRoot,
		NameNodePeerID:       nameNodePeerID,
	}

	err = dataNode.Init(soloOSEnv,
		options,
		memBlockDriver, netBlockDriver, netINodeDriver)
	util.AssertErrIsNil(err)
}
