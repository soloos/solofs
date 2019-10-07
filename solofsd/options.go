package solofsd

import (
	"soloos/common/iron"
	"soloos/common/snet"
)

type Options struct {
	DefaultNetBlockCap    int
	DefaultMemBlockCap    int
	DefaultMemBlocksLimit int32
	Mode                  string

	SNetDriverServeAddr string
	SoloboatWebPeerID   string

	SrpcServeAddr  string
	SrpcListenAddr string
	WebServer      iron.Options

	SolodnSrpcPeerID  string
	SolodnWebPeerID   string
	SolodnLocalFsRoot string

	SolonnSrpcPeerID string
	SolonnWebPeerID  string

	PProfListenAddr string
	DBDriver        string
	Dsn             string

	HeartBeatServers []snet.HeartBeatServerOptions
}
