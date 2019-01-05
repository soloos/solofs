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

func (p *Env) Init(optionsFile string) {
	var (
		options Options
		err     error
	)

	options, err = LoadOptionsFile(optionsFile)
	util.AssertErrIsNil(err)

	go func() {
		util.PProfServe(options.PProfListenAddr)
	}()

	util.AssertErrIsNil(p.Client.Init(options.NameNodeSRPCServerAddr,
		options.MemBlockChunkSize, options.MemBlockChunksLimit,
		options.DBDriver, options.Dsn,
	))
}
