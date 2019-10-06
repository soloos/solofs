package solofsd

import "soloos/common/util"

func (p *SolofsDaemon) startSolonn() {
	copy(p.srpcPeerID[:], []byte(p.options.SolonnSrpcPeerID))
	copy(p.webPeerID[:], []byte(p.options.SolonnWebPeerID))

	util.AssertErrIsNil(p.NetBlockDriver.Init(&p.SoloosEnv,
		nil, nil, p.MetaStg.PrepareNetBlockMetaData))

	util.AssertErrIsNil(p.NetINodeDriver.Init(&p.SoloosEnv,
		&p.NetBlockDriver,
		&p.MemBlockDriver,
		nil,
		p.MetaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		p.MetaStg.PrepareNetINodeMetaDataWithStorDB,
		p.MetaStg.NetINodeCommitSizeInDB,
	))

	util.AssertErrIsNil(p.solonn.Init(&p.SoloosEnv,
		p.srpcPeerID, p.options.SrpcListenAddr, p.options.SrpcServeAddr,
		p.webPeerID, p.options.WebServer,
		&p.MetaStg,
		&p.MemBlockDriver,
		&p.NetBlockDriver,
		&p.NetINodeDriver,
	))
	util.AssertErrIsNil(p.solonn.SetHeartBeatServers(p.options.HeartBeatServers))
	util.AssertErrIsNil(p.solonn.Serve())
	util.AssertErrIsNil(p.solonn.Close())
}
