package memstg

import (
	"soloos/common/snet"
	"soloos/common/solofsprotocol"
	"soloos/common/solofstypes"
)

func (p *NetINodeDriver) NetINodeSyncToNet(peerID snet.PeerID,
	uNetINode solofstypes.NetINodeUintptr) error {
	var req = solofsprotocol.NetINodeSyncReq{
		NetINodeID: uNetINode.Ptr().ID,
	}
	return p.SrpcClientDriver.SimpleCall(peerID,
		"/NetINode/Sync", nil, req)
}
