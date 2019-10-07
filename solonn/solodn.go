package solonn

import "soloos/common/snet"

func (p *Solonn) SolodnRegister(peer snet.Peer) error {
	var err = p.SoloosEnv.SNetDriver.RegisterPeer(peer)
	if err != nil {
		return err
	}

	return p.metaStg.SolodnRegister(peer)
}
