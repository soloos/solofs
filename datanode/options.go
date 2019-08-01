package datanode

import "soloos/common/snettypes"

type DataNodeOptions struct {
	SRPCPeerID           snettypes.PeerID
	SRPCServerListenAddr string
	SRPCServerServeAddr  string
	WebServerListenAddr  string
	WebServerServeAddr   string
	LocalFSRoot          string
	NameNodeSRPCPeerID   snettypes.PeerID
}
