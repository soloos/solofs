package main

import (
	"os"
	"soloos/common/util"
	"soloos/sdfs/sdfsd"
)

func main() {
	var (
		sdfsdIns sdfsd.SDFSD
		options  sdfsd.Options
		err      error
	)

	optionsFile := os.Args[1]

	err = util.LoadOptionsFile(optionsFile, &options)
	util.AssertErrIsNil(err)

	sdfsdIns.Init(options)
	sdfsdIns.Start()
}
