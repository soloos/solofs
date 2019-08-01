package main

import (
	"soloos/common/fsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/sdfs/sdfssdk"
)

var (
	env Env
)

type Env struct {
	Options      Options
	SoloOSEnv    soloosbase.SoloOSEnv
	ClientDriver sdfssdk.ClientDriver
	Client       sdfssdk.Client
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

	var nameNodeSRPCPeerID snettypes.PeerID
	nameNodeSRPCPeerID.SetStr(p.Options.NameNodeSRPCPeerID)
	util.AssertErrIsNil(p.ClientDriver.Init(&p.SoloOSEnv,
		nameNodeSRPCPeerID,
		p.Options.DBDriver, p.Options.Dsn,
	))

	if p.Options.DefaultNetBlockCap == 0 {
		p.Options.DefaultNetBlockCap = p.Options.DefaultMemBlockCap
	}

	util.AssertErrIsNil(
		p.ClientDriver.InitClient(&p.Client,
			sdfsapitypes.NameSpaceID(p.Options.NameSpaceID),
			p.Options.DefaultNetBlockCap,
			p.Options.DefaultMemBlockCap,
			p.Options.DefaultMemBlocksLimit,
		))

	p.PosixFS = p.Client.GetPosixFS()
}
