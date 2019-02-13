package datanode

import snettypes "soloos/common/snet/types"

type DataNodeOptions struct {
	PeerID               snettypes.PeerID
	SrpcServerListenAddr string
	SrpcServerServeAddr  string
	LocalFsRoot          string
	NameNodePeerID       snettypes.PeerID
	NameNodeSRPCServer   string
}
