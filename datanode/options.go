package datanode

import (
	"soloos/common/iron"
	"soloos/common/snettypes"
)

type DataNodeOptions struct {
	SRPCPeerID           snettypes.PeerID
	SRPCServerListenAddr string
	SRPCServerServeAddr  string
	WebPeerID            snettypes.PeerID
	WebServer            iron.Options
	LocalFSRoot          string
	NameNodeSRPCPeerID   snettypes.PeerID
}
