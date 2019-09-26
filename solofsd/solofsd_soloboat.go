package solofsd

import "soloos/common/snettypes"

func (p *SolofsDaemon) initSoloboat() error {
	return p.soloboatClient.Init(&p.SoloosEnv, snettypes.StrToPeerID(p.options.SoloboatWebPeerID))
}
