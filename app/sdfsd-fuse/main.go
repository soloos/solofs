package main

import (
	"os"
	"soloos/sdfs/libsdfs"
	"soloos/sfuse"
	"soloos/util"
)

func main() {
	optionsFile := os.Args[1]

	options, err := LoadOptionsFile(optionsFile)
	util.AssertErrIsNil(err)

	var (
		clientDriver libsdfs.ClientDriver
		sfuseServer  sfuse.Server
	)

	err = clientDriver.Init(options.NameNodeSRPCServerAddr,
		options.MemBlockChunkSize,
		options.MemBlockChunksLimit,
		options.DBDriver, options.Dsn)
	util.AssertErrIsNil(err)

	err = sfuseServer.Init(options.SFuseOptions, &clientDriver)
	util.AssertErrIsNil(err)

	err = sfuseServer.Serve()
	util.AssertErrIsNil(err)
}
