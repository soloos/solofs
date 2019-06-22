package main

import (
	"os"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/sdfs/sdfssdk"
	"soloos/sdfs/sfuse"
)

func main() {
	optionsFile := os.Args[1]

	options, err := LoadOptionsFile(optionsFile)
	util.AssertErrIsNil(err)

	var (
		soloOSEnv    soloosbase.SoloOSEnv
		clientDriver sdfssdk.ClientDriver
		sfuseServer  sfuse.Server
	)

	err = soloOSEnv.Init()
	util.AssertErrIsNil(err)
	err = soloOSEnv.SNetDriver.StartClient(options.SNetDriverServeAddr)
	util.AssertErrIsNil(err)

	if options.PProfListenAddr != "" {
		go util.PProfServe(options.PProfListenAddr)
	}

	{
		var nameNodePeerID snettypes.PeerID
		nameNodePeerID.SetStr(options.NameNodePeerID)
		err = clientDriver.Init(&soloOSEnv,
			nameNodePeerID,
			options.DBDriver, options.Dsn)
		util.AssertErrIsNil(err)
	}

	err = sfuseServer.Init(options.SFuseOptions,
		options.DefaultNetBlockCap,
		options.DefaultMemBlockCap,
		options.DefaultMemBlocksLimit,
		&clientDriver)
	util.AssertErrIsNil(err)

	err = sfuseServer.Serve()
	util.AssertErrIsNil(err)

	err = sfuseServer.Close()
	util.AssertErrIsNil(err)
}
