package main

import "C"
import "unsafe"

//export GoSdfsOpenFile
func GoSdfsOpenFile(cInodePath *C.char, flags,
	bufferSize C.int, replication C.short, blocksize C.long) (uint64, C.int) {
	return 0, 0
}

//export GoSdfsExists
func GoSdfsExists(inodePath *C.char) C.int {
	return 0
}

//export GoSdfsListDirectory
func GoSdfsListDirectory(inodePath *C.char, ret *unsafe.Pointer, num *C.int) {
}

//export GoSdfsCreateDirectory
func GoSdfsCreateDirectory(inodePath *C.char) C.int {
	return 0
}

//export GoSdfsDelete
func GoSdfsDelete(inodePath *C.char, recursive C.int) C.int {
	return 0
}

//export GoSdfsRename
func GoSdfsRename(oldINodePath, newINodePath *C.char) C.int {
	return 0
}

//export GoSdfsGetPathInfo
func GoSdfsGetPathInfo(cInodePath *C.char) (inodeID uint64, size int64, mTime uint64, code C.int) {
	return 0, 0, 0, 0
}

//export GoSdfsGetINodeInfo
func GoSdfsGetINodeInfo(inodeID uint64) (size int64, mTime uint64, code C.int) {
	return 0, 0, 0
}
