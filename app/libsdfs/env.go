package main

import (
	"soloos/sdfs/libsdfs"
	"soloos/util"
)

var (
	env Env
)

type Env struct {
	Options Options
	Client  libsdfs.Client
}

func (p *Env) Init(optionsFile string) {
	var (
		err error
	)

	p.Options, err = LoadOptionsFile(optionsFile)
	util.AssertErrIsNil(err)

	go func() {
		util.PProfServe(p.Options.PProfListenAddr)
	}()

	util.AssertErrIsNil(p.Client.Init(p.Options.NameNodeSRPCServerAddr,
		p.Options.MemBlockChunkSize, p.Options.MemBlockChunksLimit,
		p.Options.DBDriver, p.Options.Dsn,
	))
}
