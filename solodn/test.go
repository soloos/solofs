package solodn

import (
	"path/filepath"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/solofs/memstg"
)

func MakeSolodnForTest(soloOSEnv *soloosbase.SoloOSEnv,
	solodn *Solodn,
	solodnSRPCPeerID snettypes.PeerID, solodnSRPCServerAddr string,
	solonnSRPCPeerID snettypes.PeerID, solonnSRPCServerAddr string,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *memstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) {
	var (
		err error
	)

	var localFSRoot = filepath.Join("/tmp/soloos_test.data", solodnSRPCPeerID.Str())

	var options = SolodnOptions{
		SRPCPeerID:           solodnSRPCPeerID,
		SRPCServerListenAddr: solodnSRPCServerAddr,
		SRPCServerServeAddr:  solodnSRPCServerAddr,
		LocalFSRoot:          localFSRoot,
		SolonnSRPCPeerID:   solonnSRPCPeerID,
	}

	err = solodn.Init(soloOSEnv,
		options,
		memBlockDriver, netBlockDriver, netINodeDriver)
	util.AssertErrIsNil(err)
}
