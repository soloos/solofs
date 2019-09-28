package solodn

import (
	"soloos/common/iron"
	"soloos/common/snettypes"
)

type SolodnOptions struct {
	SRPCPeerID           snettypes.PeerID
	SRPCServerListenAddr string
	SRPCServerServeAddr  string
	WebPeerID            snettypes.PeerID
	WebServer            iron.Options
	LocalFsRoot          string
	SolonnSRPCPeerID   snettypes.PeerID
}
