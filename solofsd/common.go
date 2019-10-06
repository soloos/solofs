package solofsd

import "soloos/common/util"

func (p *SolofsDaemon) startCommon() {
	if p.options.PProfListenAddr != "" {
		go util.PProfServe(p.options.PProfListenAddr)
	}
}
