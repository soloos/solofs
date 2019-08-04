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
	dataNodeSRPCPeerID snettypes.PeerID, dataNodeSRPCServerAddr string,
	nameNodeSRPCPeerID snettypes.PeerID, nameNodeSRPCServerAddr string,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *memstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) {
	var (
		err error
	)

	var localFSRoot = filepath.Join("/tmp/soloos_test.data", dataNodeSRPCPeerID.Str())

	var options = DataNodeOptions{
		SRPCPeerID:           dataNodeSRPCPeerID,
		SRPCServerListenAddr: dataNodeSRPCServerAddr,
		SRPCServerServeAddr:  dataNodeSRPCServerAddr,
		LocalFSRoot:          localFSRoot,
		NameNodeSRPCPeerID:   nameNodeSRPCPeerID,
	}

	err = dataNode.Init(soloOSEnv,
		options,
		memBlockDriver, netBlockDriver, netINodeDriver)
	util.AssertErrIsNil(err)
}
