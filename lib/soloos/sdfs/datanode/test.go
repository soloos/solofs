package datanode

import (
	"path/filepath"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/netstg"
)

func MakeDataNodeForTest(soloOSEnv *soloosbase.SoloOSEnv,
	dataNode *DataNode,
	dataNodeSRPCServerAddr string,
	metaStg *metastg.MetaStg,
	nameNodeSRPCServerAddr string,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) {
	var (
		err error
	)

	var peerID snettypes.PeerID
	snettypes.InitTmpPeerID(&peerID)
	var localFSRoot = filepath.Join("/tmp/sdfs_test.data", string(peerID[:3]))

	var options = DataNodeOptions{
		PeerID:               peerID,
		SrpcServerListenAddr: dataNodeSRPCServerAddr,
		SrpcServerServeAddr:  dataNodeSRPCServerAddr,
		LocalFSRoot:          localFSRoot,
		NameNodeSRPCServer:   nameNodeSRPCServerAddr,
	}

	err = dataNode.Init(soloOSEnv,
		options,
		metaStg,
		memBlockDriver, netBlockDriver, netINodeDriver)
	util.AssertErrIsNil(err)
}
