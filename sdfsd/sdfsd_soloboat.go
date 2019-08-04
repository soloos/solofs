package sdfsd

import "soloos/common/snettypes"

func (p *SDFSD) initSoloBoat() error {
	return p.soloboatClient.Init(&p.SoloOSEnv, snettypes.StrToPeerID(p.options.SoloBoatWebPeerID))
}
