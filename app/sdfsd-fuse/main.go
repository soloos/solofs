package main

import (
	"os"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/sdfs/libsdfs"
	"soloos/sdfs/sfuse"
)

func main() {
	optionsFile := os.Args[1]

	options, err := LoadOptionsFile(optionsFile)
	util.AssertErrIsNil(err)

	var (
		soloOSEnv    soloosbase.SoloOSEnv
		clientDriver libsdfs.ClientDriver
		sfuseServer  sfuse.Server
	)

	err = soloOSEnv.Init()
	util.AssertErrIsNil(err)

	if options.PProfListenAddr != "" {
		go util.PProfServe(options.PProfListenAddr)
	}

	err = clientDriver.Init(&soloOSEnv, options.NameNodeSRPCServerAddr,
		options.DBDriver, options.Dsn)
	util.AssertErrIsNil(err)

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
