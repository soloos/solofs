package main

import "C"

import (
	"os"
	"soloos/common/log"
	"sync"
)

var (
	mutexInit sync.Mutex
	isInited  = false
)

func main() {
}

//export GoSdfsInit
func GoSdfsInit() {
	initEnv()
}

func initEnv() {
	mutexInit.Lock()
	defer mutexInit.Unlock()

	if isInited {
		return
	}

	env.Init(os.Args[1])
	isInited = true
	log.Info("init success.")
}
