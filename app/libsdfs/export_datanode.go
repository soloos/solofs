package main

import "C"
import (
	"io"
	"reflect"
	"soloos/common/log"
	"soloos/common/sdfsapi"
	"soloos/sdfs/types"
	"unsafe"
)

//export GoSdfsPappend
func GoSdfsPappend(fdID uint64, buffer unsafe.Pointer, bufferLen, offset int32) (int32, C.int) {
	var (
		fsINode types.FsINode
		fd      = env.PosixFS.FdTableGetFd(fdID)
		err     error
	)

	err = env.PosixFS.FetchFsINodeByID(fd.FsINodeID, &fsINode)
	if err != nil {
		return 0, sdfsapi.CODE_ERR
	}

	var data = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufferLen),
		Cap:  int(bufferLen),
	}))

	err = env.PosixFS.SimpleWriteWithMem(fsINode.UNetINode, data, uint64(offset))
	if err != nil {
		return 0, sdfsapi.CODE_ERR
	}

	return bufferLen, 0
}

//export GoSdfsAppend
func GoSdfsAppend(fdID uint64, buffer unsafe.Pointer, bufferLen int32) (int32, C.int) {
	var (
		fsINode types.FsINode
		fd      = env.PosixFS.FdTableGetFd(fdID)
		err     error
	)

	err = env.PosixFS.FetchFsINodeByID(fd.FsINodeID, &fsINode)
	if err != nil {
		return 0, sdfsapi.CODE_ERR
	}

	var data = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufferLen),
		Cap:  int(bufferLen),
	}))
	err = env.PosixFS.SimpleWriteWithMem(fsINode.UNetINode, data, fd.AppendPosition)
	if err != nil {
		log.Warn(err)
		return 0, sdfsapi.CODE_ERR
	}

	env.PosixFS.FdTableFdAddAppendPosition(fdID, uint64(bufferLen))

	return bufferLen, sdfsapi.CODE_OK
}

//export GoSdfsRead
func GoSdfsRead(fdID uint64, buffer unsafe.Pointer, bufferLen int32) (int32, C.int) {
	var (
		fsINode        types.FsINode
		fd             = env.PosixFS.FdTableGetFd(fdID)
		readDataLength int
		err            error
	)

	err = env.PosixFS.FetchFsINodeByID(fd.FsINodeID, &fsINode)
	if err != nil {
		return 0, sdfsapi.CODE_ERR
	}

	var data = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufferLen),
		Cap:  int(bufferLen),
	}))
	readDataLength, err = env.PosixFS.SimpleReadWithMem(fsINode.UNetINode, data, fd.ReadPosition)
	if err != nil && err != io.EOF {
		log.Warn(err, readDataLength)
		return int32(readDataLength), sdfsapi.CODE_ERR
	}

	env.PosixFS.FdTableFdAddReadPosition(fdID, uint64(bufferLen))

	return int32(readDataLength), sdfsapi.CODE_OK
}

//export GoSdfsPread
func GoSdfsPread(fdID uint64, buffer unsafe.Pointer, bufferLen int32, position uint64) (int32, C.int) {
	var (
		fsINode        types.FsINode
		fd             = env.PosixFS.FdTableGetFd(fdID)
		readDataLength int
		err            error
	)

	var data = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufferLen),
		Cap:  int(bufferLen),
	}))
	err = env.PosixFS.FetchFsINodeByID(fd.FsINodeID, &fsINode)
	if err != nil {
		log.Warn(err)
		return 0, sdfsapi.CODE_ERR
	}

	readDataLength, err = env.PosixFS.SimpleReadWithMem(fsINode.UNetINode, data, position)
	if err != nil {
		return int32(readDataLength), sdfsapi.CODE_ERR
	}

	return int32(readDataLength), sdfsapi.CODE_OK
}

//export GoSdfsCloseFile
func GoSdfsCloseFile(fdID uint64) C.int {
	ret := doFlushINode(fdID)
	// env.PosixFS.FdTableReleaseFd(fdID)
	return ret
}

//export GoSdfsFlushFile
func GoSdfsFlushFile(fdID uint64) C.int {
	return doFlushINode(fdID)
}

//export GoSdfsHFlushINode
func GoSdfsHFlushINode(fdID uint64) C.int {
	return doFlushINode(fdID)
}

//export GoSdfsHSyncINode
func GoSdfsHSyncINode(fdID uint64) C.int {
	return doFlushINode(fdID)
}

func doFlushINode(fdID uint64) C.int {
	var (
		fsINode types.FsINode
		fd      = env.PosixFS.FdTableGetFd(fdID)
		err     error
	)

	err = env.PosixFS.FetchFsINodeByID(fd.FsINodeID, &fsINode)
	if err != nil {
		return sdfsapi.CODE_ERR
	}

	if fsINode.UNetINode != 0 {
		env.PosixFS.SimpleFlush(fsINode.UNetINode)
	}

	return sdfsapi.CODE_OK
}
