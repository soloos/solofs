package solodn

import (
	"path/filepath"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/solofs/memstg"
)

func MakeSolodnForTest(soloosEnv *soloosbase.SoloosEnv,
	solodn *Solodn,
	solodnSrpcPeerID snettypes.PeerID, solodnSrpcServerAddr string,
	solonnSrpcPeerID snettypes.PeerID, solonnSrpcServerAddr string,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *memstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) {
	var (
		err error
	)

	var localFsRoot = filepath.Join("/tmp/soloos_test.data", solodnSrpcPeerID.Str())

	var options = SolodnOptions{
		SrpcPeerID:           solodnSrpcPeerID,
		SrpcServerListenAddr: solodnSrpcServerAddr,
		SrpcServerServeAddr:  solodnSrpcServerAddr,
		LocalFsRoot:          localFsRoot,
		SolonnSrpcPeerID:   solonnSrpcPeerID,
	}

	err = solodn.Init(soloosEnv,
		options,
		memBlockDriver, netBlockDriver, netINodeDriver)
	util.AssertErrIsNil(err)
}
