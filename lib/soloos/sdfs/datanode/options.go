package datanode

type DataNodeSRPCServerOptions struct {
	Network    string
	ListenAddr string
}

type DataNodeOptions struct {
	SRPCServer DataNodeSRPCServerOptions
}
