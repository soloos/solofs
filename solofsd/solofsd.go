package solofsd

import (
	"soloos/common/solofsapi"
	"soloos/common/snet"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/solodb/offheap"
	"soloos/solofs/solodn"
	"soloos/solofs/memstg"
	"soloos/solofs/metastg"
	"soloos/solofs/solonn"
	"soloos/soloboat/soloboatsdk"
)

type SolofsDaemon struct {
	options          Options
	SoloosEnv        soloosbase.SoloosEnv
	offheapDriver    *offheap.OffheapDriver
	SNetDriver       snet.NetDriver
	SNetClientDriver snet.SRPCClientDriver
	MetaStg          metastg.MetaStg
	SolodnClient   solofsapi.SolodnClient
	MemBlockDriver   memstg.MemBlockDriver
	NetBlockDriver   memstg.NetBlockDriver
	NetINodeDriver   memstg.NetINodeDriver

	solonn   solonn.Solonn
	solodn   solodn.Solodn
	srpcPeerID snettypes.PeerID
	webPeerID  snettypes.PeerID

	soloboatClient soloboatsdk.Client
}

func (p *SolofsDaemon) initMetaStg() error {
	return p.MetaStg.Init(&p.SoloosEnv,
		p.options.DBDriver, p.options.Dsn)
}

func (p *SolofsDaemon) initMemStg() error {
	var memBlockDriverOptions = memstg.MemBlockDriverOptions{
		[]memstg.MemBlockTableOptions{
			memstg.MemBlockTableOptions{
				p.options.DefaultMemBlockCap,
				p.options.DefaultMemBlocksLimit,
			},
		},
	}
	return (p.MemBlockDriver.Init(&p.SoloosEnv, memBlockDriverOptions))
}

func (p *SolofsDaemon) Init(options Options) {
	p.options = options
	util.AssertErrIsNil(p.SoloosEnv.InitWithSNet(p.options.SNetDriverServeAddr))

	p.SolodnClient.Init(&p.SoloosEnv)

	util.AssertErrIsNil(p.initMetaStg())

	util.AssertErrIsNil(p.initMemStg())

	util.AssertErrIsNil(p.initSoloboat())
}

func (p *SolofsDaemon) startCommon() {
	if p.options.PProfListenAddr != "" {
		go util.PProfServe(p.options.PProfListenAddr)
	}
}

func (p *SolofsDaemon) startSolonn() {
	copy(p.srpcPeerID[:], []byte(p.options.SolonnSRPCPeerID))
	copy(p.webPeerID[:], []byte(p.options.SolonnWebPeerID))

	util.AssertErrIsNil(p.NetBlockDriver.Init(&p.SoloosEnv,
		nil, &p.SolodnClient, p.MetaStg.PrepareNetBlockMetaData))

	util.AssertErrIsNil(p.NetINodeDriver.Init(&p.SoloosEnv,
		&p.NetBlockDriver,
		&p.MemBlockDriver,
		nil,
		p.MetaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		p.MetaStg.PrepareNetINodeMetaDataWithStorDB,
		p.MetaStg.NetINodeCommitSizeInDB,
	))

	util.AssertErrIsNil(p.solonn.Init(&p.SoloosEnv,
		p.srpcPeerID, p.options.SRPCListenAddr, p.options.SRPCServeAddr,
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

func (p *SolofsDaemon) startSolodn() {
	var (
		solonnSRPCPeerID snettypes.PeerID
		solodnOptions    solodn.SolodnOptions
	)

	copy(p.srpcPeerID[:], []byte(p.options.SolodnSRPCPeerID))
	copy(p.webPeerID[:], []byte(p.options.SolodnWebPeerID))
	copy(solonnSRPCPeerID[:], []byte(p.options.SolonnSRPCPeerID))

	solodnOptions = solodn.SolodnOptions{
		SRPCPeerID:           p.srpcPeerID,
		SRPCServerListenAddr: p.options.SRPCListenAddr,
		SRPCServerServeAddr:  p.options.SRPCServeAddr,
		WebPeerID:            p.webPeerID,
		WebServer:            p.options.WebServer,
		LocalFsRoot:          p.options.SolodnLocalFsRoot,
		SolonnSRPCPeerID:   solonnSRPCPeerID,
	}

	util.AssertErrIsNil(p.NetBlockDriver.Init(&p.SoloosEnv,
		nil, &p.SolodnClient, p.MetaStg.PrepareNetBlockMetaData))

	util.AssertErrIsNil(p.NetINodeDriver.Init(&p.SoloosEnv,
		&p.NetBlockDriver,
		&p.MemBlockDriver,
		nil,
		p.MetaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		p.MetaStg.PrepareNetINodeMetaDataWithStorDB,
		p.MetaStg.NetINodeCommitSizeInDB,
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

func (p *SolofsDaemon) Start() {
	if p.options.Mode == "solonn" {
		p.startCommon()
		p.startSolonn()
	}

	if p.options.Mode == "solodn" {
		p.startCommon()
		p.startSolodn()
	}
}
