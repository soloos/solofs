package datanode

import "soloos/common/snettypes"

type DataNodeOptions struct {
	PeerID               snettypes.PeerID
	SrpcServerListenAddr string
	SrpcServerServeAddr  string
	LocalFSRoot          string
	NameNodePeerID       snettypes.PeerID
}
