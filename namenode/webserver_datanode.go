package namenode

import "soloos/common/iron"

func (p *WebServer) ctrDataNodeHeartBeat(ir *iron.Request) {
	ir.ApiOutput(nil, iron.CODE_OK, "")
}
