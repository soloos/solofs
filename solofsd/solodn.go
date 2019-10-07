package solofsd

import (
	"soloos/common/snettypes"
	"soloos/common/util"
	"soloos/solofs/solodn"
)

func (p *SolofsDaemon) startSolodn() {
	var (
		solonnSrpcPeerID snettypes.PeerID
		solodnOptions    solodn.SolodnOptions
	)

	copy(p.srpcPeerID[:], []byte(p.options.SolodnSrpcPeerID))
	copy(p.webPeerID[:], []byte(p.options.SolodnWebPeerID))
	copy(solonnSrpcPeerID[:], []byte(p.options.SolonnSrpcPeerID))

	solodnOptions = solodn.SolodnOptions{
		SrpcPeerID:           p.srpcPeerID,
		SrpcServerListenAddr: p.options.SrpcListenAddr,
		SrpcServerServeAddr:  p.options.SrpcServeAddr,
		WebPeerID:            p.webPeerID,
		WebServer:            p.options.WebServer,
		LocalFsRoot:          p.options.SolodnLocalFsRoot,
		SolonnSrpcPeerID:     solonnSrpcPeerID,
	}

	util.AssertErrIsNil(p.NetBlockDriver.Init(&p.SoloosEnv,
		nil, &p.SolodnClient, nil, nil, nil, nil))

	util.AssertErrIsNil(p.NetINodeDriver.Init(&p.SoloosEnv,
		&p.NetBlockDriver,
		&p.MemBlockDriver,
		nil, nil, nil, nil,
	))

	util.AssertErrIsNil(p.solodn.Init(&p.SoloosEnv, solodnOptions,
		&p.MemBlockDriver,
		&p.NetBlockDriver,
		&p.NetINodeDriver,
	))
	util.AssertErrIsNil(p.solodn.SetHeartBeatServers(p.options.HeartBeatServers))
	util.AssertErrIsNil(p.solodn.Serve())
	util.AssertErrIsNil(p.solodn.Close())
}
