package main

import "C"
import (
	"unsafe"
)

//export GoSdfsPappend
func GoSdfsPappend(inodeID uint64, buffer unsafe.Pointer, bufferLen, offset int32) (int32, C.int) {
	return bufferLen, 0
}

//export GoSdfsAppend
func GoSdfsAppend(inodeID uint64, buffer unsafe.Pointer, bufferLen int32) (int32, C.int) {
	return bufferLen, 0
}

//export GoSdfsRead
func GoSdfsRead(inodeID uint64, buffer unsafe.Pointer, bufferLen int32) (int32, C.int) {
	return bufferLen, 0
}

//export GoSdfsPread
func GoSdfsPread(inodeID uint64, buffer unsafe.Pointer, bufferLen int32, position int64) (int32, C.int) {
	return bufferLen, 0
}

//export GoSdfsCloseFile
func GoSdfsCloseFile(inodeID uint64) C.int {
	return 0
}

//export GoSdfsFlushFile
func GoSdfsFlushFile(inodeID uint64) C.int {
	return 0
}

//export GoSdfsHFlushINode
func GoSdfsHFlushINode(inodeID uint64) C.int {
	return 0
}

//export GoSdfsHSyncINode
func GoSdfsHSyncINode(inodeID uint64) C.int {
	return 0
}
