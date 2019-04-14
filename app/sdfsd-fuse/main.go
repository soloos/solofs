package main

import (
	"os"
	"soloos/common/util"
	"soloos/sdfs/libsdfs"
	"soloos/sdfs/sfuse"
)

func main() {
	optionsFile := os.Args[1]

	options, err := LoadOptionsFile(optionsFile)
	util.AssertErrIsNil(err)

	var (
		clientDriver libsdfs.ClientDriver
		sfuseServer  sfuse.Server
	)

	if options.PProfListenAddr != "" {
		go util.PProfServe(options.PProfListenAddr)
	}

	err = clientDriver.Init(options.NameNodeSRPCServerAddr,
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
