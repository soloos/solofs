package solodn

import (
	"soloos/common/iron"
	"soloos/common/snettypes"
)

type SolodnOptions struct {
	SrpcPeerID           snettypes.PeerID
	SrpcServerListenAddr string
	SrpcServerServeAddr  string
	WebPeerID            snettypes.PeerID
	WebServer            iron.Options
	LocalFsRoot          string
	SolonnSrpcPeerID   snettypes.PeerID
}
