package main

import (
	"soloos/common/fsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/solofs/solofssdk"
)

var (
	env Env
)

type Env struct {
	Options      Options
	SoloOSEnv    soloosbase.SoloOSEnv
	ClientDriver solofssdk.ClientDriver
	Client       solofssdk.Client
	PosixFS      fsapi.PosixFS
}

func (p *Env) Init(optionsFile string) {
	var (
		err error
	)

	p.Options, err = LoadOptionsFile(optionsFile)
	util.AssertErrIsNil(err)

	err = p.SoloOSEnv.InitWithSNet(p.Options.SNetDriverServeAddr)
	util.AssertErrIsNil(err)

	go func() {
		util.PProfServe(p.Options.PProfListenAddr)
	}()

	var solonnSRPCPeerID snettypes.PeerID
	solonnSRPCPeerID.SetStr(p.Options.SolonnSRPCPeerID)
	util.AssertErrIsNil(p.ClientDriver.Init(&p.SoloOSEnv,
		solonnSRPCPeerID,
		p.Options.DBDriver, p.Options.Dsn,
	))

	if p.Options.DefaultNetBlockCap == 0 {
		p.Options.DefaultNetBlockCap = p.Options.DefaultMemBlockCap
	}

	util.AssertErrIsNil(
		p.ClientDriver.InitClient(&p.Client,
			solofsapitypes.NameSpaceID(p.Options.NameSpaceID),
			p.Options.DefaultNetBlockCap,
			p.Options.DefaultMemBlockCap,
			p.Options.DefaultMemBlocksLimit,
		))

	p.PosixFS = p.Client.GetPosixFS()
}
