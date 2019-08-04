package sdfsd

import (
	"soloos/common/iron"
	"soloos/common/sdfsapitypes"
)

type Options struct {
	DefaultNetBlockCap    int
	DefaultMemBlockCap    int
	DefaultMemBlocksLimit int32
	Mode                  string

	SNetDriverServeAddr string
	SoloBoatWebPeerID   string

	SRPCServeAddr  string
	SRPCListenAddr string
	WebServer      iron.Options

	DataNodeSRPCPeerID  string
	DataNodeWebPeerID   string
	DataNodeLocalFSRoot string

	NameNodeSRPCPeerID string
	NameNodeWebPeerID  string

	PProfListenAddr string
	DBDriver        string
	Dsn             string

	HeartBeatServers []sdfsapitypes.HeartBeatServerOptions
}
