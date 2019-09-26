package metastg

import (
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
)

type SolodnDriver struct {
	*soloosbase.SoloosEnv
	metaStg *MetaStg

	chooseSolodnIndex         uint32
	solodnsForBlockRegistered map[snettypes.PeerID]int64
	solodnsForBlock           []snettypes.PeerID
	solodnsForBlockRWMutex    util.RWMutex
}

func (p *SolodnDriver) Init(metaStg *MetaStg) error {
	p.SoloosEnv = metaStg.SoloosEnv
	p.metaStg = metaStg
	p.solodnsForBlockRegistered = make(map[snettypes.PeerID]int64)
	return nil
}

func (p *SolodnDriver) SolodnRegister(peer snettypes.Peer) error {
	var (
		registered bool
	)

	p.solodnsForBlockRWMutex.Lock()
	_, registered = p.solodnsForBlockRegistered[peer.ID]
	if registered == false {
		p.SNetDriver.RegisterPeer(peer)
		p.solodnsForBlockRegistered[peer.ID] = 0
		p.solodnsForBlock = append(p.solodnsForBlock, peer.ID)
	}
	p.solodnsForBlockRWMutex.Unlock()

	return nil
}

func (p *SolodnDriver) GetSolodn(peerID snettypes.PeerID) (snettypes.Peer, error) {
	return p.SNetDriver.GetPeer(peerID)
}
