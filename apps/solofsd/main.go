package main

import (
	"os"
	"soloos/common/util"
	"soloos/solofs/solofsd"
)

func main() {
	var (
		solofsdIns solofsd.SolofsDaemon
		options    solofsd.Options
		err        error
	)

	optionsFile := os.Args[1]

	err = util.LoadOptionsFile(optionsFile, &options)
	util.AssertErrIsNil(err)

	solofsdIns.Init(options)
	solofsdIns.Start()
}
