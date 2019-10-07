package solodn

import (
	"soloos/common/iron"
	"soloos/common/snet"
)

type SolodnOptions struct {
	SrpcPeerID           snet.PeerID
	SrpcServerListenAddr string
	SrpcServerServeAddr  string
	WebPeerID            snet.PeerID
	WebServer            iron.Options
	LocalFsRoot          string
	SolonnSrpcPeerID     snet.PeerID
}
