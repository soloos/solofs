package main

import (
	"soloos/common/fsapi"
	"soloos/common/util"
	"soloos/sdfs/libsdfs"
)

var (
	env Env
)

type Env struct {
	Options      Options
	ClientDriver libsdfs.ClientDriver
	Client       libsdfs.Client
	PosixFS      fsapi.PosixFS
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
		p.Options.DBDriver, p.Options.Dsn,
	))

	if p.Options.DefaultNetBlockCap == 0 {
		p.Options.DefaultNetBlockCap = p.Options.DefaultMemBlockCap
	}

	util.AssertErrIsNil(
		p.ClientDriver.InitClient(&p.Client,
			p.Options.DefaultNetBlockCap,
			p.Options.DefaultMemBlockCap,
			p.Options.DefaultMemBlocksLimit,
		))

	p.PosixFS = p.Client.GetPosixFS()
}
