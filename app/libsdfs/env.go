package main

import "soloos/util"

var (
	env Env
)

type Env struct {
	Client Client
}

func (p *Env) Init(nameNodeSRPCServerAddr string,
	memBlockChunkSize int, memBlockChunksLimit int32) {
	util.AssertErrIsNil(p.Client.Init(nameNodeSRPCServerAddr,
		memBlockChunkSize, memBlockChunksLimit))
}
