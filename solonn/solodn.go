package solonn

import "soloos/common/snettypes"

func (p *Solonn) SolodnRegister(peer snettypes.Peer) error {
	var err = p.SoloosEnv.SNetDriver.RegisterPeer(peer)
	if err != nil {
		return err
	}

	return p.metaStg.SolodnRegister(peer)
}
