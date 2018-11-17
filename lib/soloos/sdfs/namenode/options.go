package namenode

type NameNodeSRPCServerOptions struct {
	Network    string
	ListenAddr string
}

type NameNodeOptions struct {
	SRPCServer NameNodeSRPCServerOptions
}
