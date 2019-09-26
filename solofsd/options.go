package solofsd

import (
	"soloos/common/iron"
	"soloos/common/snettypes"
)

type Options struct {
	DefaultNetBlockCap    int
	DefaultMemBlockCap    int
	DefaultMemBlocksLimit int32
	Mode                  string

	SNetDriverServeAddr string
	SoloboatWebPeerID   string

	SRPCServeAddr  string
	SRPCListenAddr string
	WebServer      iron.Options

	SolodnSRPCPeerID  string
	SolodnWebPeerID   string
	SolodnLocalFSRoot string

	SolonnSRPCPeerID string
	SolonnWebPeerID  string

	PProfListenAddr string
	DBDriver        string
	Dsn             string

	HeartBeatServers []snettypes.HeartBeatServerOptions
}
