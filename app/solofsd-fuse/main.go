package main

import (
	"os"
	"soloos/common/snet"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/solofs/sfuse"
	"soloos/solofs/solofssdk"
)

func main() {
	optionsFile := os.Args[1]

	options, err := LoadOptionsFile(optionsFile)
	util.AssertErrIsNil(err)

	var (
		soloosEnv    soloosbase.SoloosEnv
		clientDriver solofssdk.ClientDriver
		sfuseServer  sfuse.Server
	)

	err = soloosEnv.InitWithSNet(options.SNetDriverServeAddr)
	util.AssertErrIsNil(err)

	if options.PProfListenAddr != "" {
		go util.PProfServe(options.PProfListenAddr)
	}

	{
		var solonnSrpcPeerID snet.PeerID
		solonnSrpcPeerID.SetStr(options.SolonnSrpcPeerID)
		err = clientDriver.Init(&soloosEnv,
			solonnSrpcPeerID,
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
