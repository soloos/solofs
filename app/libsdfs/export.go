package main

import "C"

import (
	"soloos/log"
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

	env.Init("127.0.0.1:9096", 1024*1024*2, 256, "sqlite3", "/tmp/sdfs.db")
	isInited = true
	log.Info("init success.")
}
