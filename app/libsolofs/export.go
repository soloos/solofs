package main

import "C"

import (
	"os"
	"soloos/common/log"
	"strings"
	"sync"
)

var (
	mutexInit sync.Mutex
	isInited  = false
)

func main() {
}

//export GoSolofsInit
func GoSolofsInit() {
	initEnv()
}

func initEnv() {
	mutexInit.Lock()
	defer mutexInit.Unlock()

	if isInited {
		return
	}

	var optionsFile string
	for _, v := range os.Args {
		if strings.Contains(v, ".json") {
			optionsFile = v
			break
		}
	}
	env.Init(optionsFile)
	isInited = true
	log.Info("init success.")
}
