package solonn

import "soloos/common/iron"

func (p *WebServer) ctrSolodnHeartBeat(ir *iron.Request) {
	ir.ApiOutput(nil, iron.CODE_OK, "")
}
