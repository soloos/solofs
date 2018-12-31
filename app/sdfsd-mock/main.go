package main

import (
	"flag"
	"soloos/util"
)

var (
	FlagWebPort  = flag.Int("webport", 10020, "listen port")
	FlagSRPCPort = flag.Int("srpcport", 10021, "listen port")
)

func main() {
	flag.Parse()

	var (
		err error
		env Env
	)

	err = env.Init(*FlagSRPCPort, *FlagWebPort)
	util.AssertErrIsNil(err)
	util.AssertErrIsNil(env.SRPCServer.Serve())
	go util.AssertErrIsNil(env.WebServer.Serve())
}
