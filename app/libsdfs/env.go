package main

import (
	"soloos/sdfs/libsdfs"
	"soloos/util"
)

var (
	env Env
)

type Env struct {
	Options      Options
	ClientDriver libsdfs.ClientDriver
	Client       libsdfs.Client
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

	util.AssertErrIsNil(p.ClientDriver.Init(p.Options.NameNodeSRPCServerAddr,
		p.Options.MemBlockChunkSize, p.Options.MemBlockChunksLimit,
		p.Options.DBDriver, p.Options.Dsn,
	))
	util.AssertErrIsNil(p.ClientDriver.InitClient(&p.Client))
}
