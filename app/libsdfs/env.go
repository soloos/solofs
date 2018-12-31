package main

import (
	"soloos/sdfs/libsdfs"
	"soloos/util"
)

var (
	env Env
)

type Env struct {
	Client libsdfs.Client
}

func (p *Env) Init(nameNodeSRPCServerAddr string,
	memBlockChunkSize int, memBlockChunksLimit int32,
	dbDriver, dsn string,
) {
	util.AssertErrIsNil(p.Client.Init(nameNodeSRPCServerAddr,
		memBlockChunkSize, memBlockChunksLimit,
		dbDriver, dsn,
	))
}
