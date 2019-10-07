package solofsd

import (
	"soloos/common/snet"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/soloboat/soloboatsdk"
	"soloos/solodb/offheap"
	"soloos/solofs/memstg"
	"soloos/solofs/metastg"
	"soloos/solofs/solodn"
	"soloos/solofs/solonn"
)

type SolofsDaemon struct {
	options        Options
	SoloosEnv      soloosbase.SoloosEnv
	offheapDriver  *offheap.OffheapDriver
	SNetDriver     snet.NetDriver
	MetaStg        metastg.MetaStg
	MemBlockDriver memstg.MemBlockDriver
	NetBlockDriver memstg.NetBlockDriver
	NetINodeDriver memstg.NetINodeDriver

	solonn     solonn.Solonn
	solodn     solodn.Solodn
	srpcPeerID snet.PeerID
	webPeerID  snet.PeerID

	soloboatClient soloboatsdk.Client
}

func (p *SolofsDaemon) initMetaStg() error {
	return p.MetaStg.Init(&p.SoloosEnv,
		p.options.DBDriver, p.options.Dsn)
}

func (p *SolofsDaemon) initMemBlockDriver() error {
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

func (p *SolofsDaemon) Init(options Options) error {
	p.options = options
	util.AssertErrIsNil(p.SoloosEnv.InitWithSNet(p.options.SNetDriverServeAddr))

	util.AssertErrIsNil(p.initMetaStg())

	util.AssertErrIsNil(p.initMemBlockDriver())

	util.AssertErrIsNil(p.initSoloboat())
	return nil
}

func (p *SolofsDaemon) Serve() error {
	if p.options.Mode == "solonn" {
		p.startCommon()
		p.startSolonn()
	}

	if p.options.Mode == "solodn" {
		p.startCommon()
		p.startSolodn()
	}
	return nil
}

func (p *SolofsDaemon) Close() error {
	return nil
}
