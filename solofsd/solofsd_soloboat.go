package solofsd

import "soloos/common/snettypes"

func (p *SolofsDaemon) initSoloBoat() error {
	return p.soloboatClient.Init(&p.SoloOSEnv, snettypes.StrToPeerID(p.options.SoloBoatWebPeerID))
}
