package memstg

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
)

func (p *NetINodeDriver) NetINodeSyncToNet(peerID snettypes.PeerID,
	uNetINode solofsapitypes.NetINodeUintptr) error {
	var req = solofsprotocol.NetINodeSyncReq{
		NetINodeID: uNetINode.Ptr().ID,
	}
	return p.SNetClientDriver.SimpleCall(peerID,
		"/NetINode/Sync", nil, req)
}
