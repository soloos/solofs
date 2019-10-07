package solofsd

import "soloos/common/snet"

func (p *SolofsDaemon) initSoloboat() error {
	return p.soloboatClient.Init(&p.SoloosEnv, snet.StrToPeerID(p.options.SoloboatWebPeerID))
}
