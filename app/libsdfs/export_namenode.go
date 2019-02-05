package main

import "C"
import (
	"soloos/sdfs/libsdfs"
	"soloos/sdfs/types"
	"unsafe"

	"github.com/hanwen/go-fuse/fuse"
)

//export GoSdfsOpenFile
func GoSdfsOpenFile(cInodePath *C.char, flags,
	bufferSize C.int, replication C.short, blocksize C.long) (uint64, C.int) {
	var (
		fsINodePath = C.GoString(cInodePath)
		fsINode     types.FsINode
		err         error
	)

	fsINode, err = env.Client.MemDirTreeStg.OpenFile(fsINodePath,
		types.DefaultNetBlockCap,
		env.Options.DefaultMemBlockCap)
	if err != nil {
		return 0, libsdfs.CODE_ERR
	}

	return env.Client.MemDirTreeStg.FdTable.AllocFd(fsINode.Ino), libsdfs.CODE_OK
}

//export GoSdfsExists
func GoSdfsExists(cInodePath *C.char) C.int {
	var (
		fsINodePath = C.GoString(cInodePath)
		fsINode     types.FsINode
		err         error
	)
	err = env.Client.MemDirTreeStg.FetchFsINodeByPath(fsINodePath, &fsINode)
	if err != nil {
		// contains err == types.ErrObjectNotExists
		return libsdfs.CODE_ERR
	}
	return libsdfs.CODE_OK
}

//export GoSdfsListDirectory
func GoSdfsListDirectory(cInodePath *C.char, ret *unsafe.Pointer, num *C.int) {
	var (
		fsINodePath = C.GoString(cInodePath)
		arr         *[1<<30 - 1]*C.char
		index       int
		err         error
	)

	err = env.Client.MemDirTreeStg.ListFsINodeByParentPath(fsINodePath, false,
		func(resultCount int) (uint64, uint64) {
			*ret = C.malloc(C.size_t(resultCount) * C.size_t(unsafe.Sizeof(uintptr(0))))
			*num = C.int(resultCount)
			arr = (*[1<<30 - 1]*C.char)(*ret)
			index = 0
			return uint64(resultCount), 0
		},
		func(fsINode types.FsINode) bool {
			arr[index] = C.CString(fsINode.Name)
			index += 1
			return true
		},
	)

	if err != nil {
		return
	}

	return
}

//export GoSdfsCreateDirectory
func GoSdfsCreateDirectory(cInodePath *C.char) C.int {
	var (
		fsINodePath = C.GoString(cInodePath)
		code        fuse.Status
	)
	code = env.Client.MemDirTreeStg.MkdirAll(0777, fsINodePath, 0, 0)
	if code != fuse.OK {
		return libsdfs.CODE_ERR
	}

	return libsdfs.CODE_OK
}

//export GoSdfsDelete
func GoSdfsDelete(cInodePath *C.char, recursive C.int) C.int {
	var (
		fsINodePath = C.GoString(cInodePath)
		err         error
	)
	err = env.Client.MemDirTreeStg.DeleteFsINodeByPath(fsINodePath)
	if err != nil {
		return libsdfs.CODE_ERR
	}

	return libsdfs.CODE_OK
}

//export GoSdfsRename
func GoSdfsRename(oldINodePath, newINodePath *C.char) C.int {
	var err error
	err = env.Client.MemDirTreeStg.RenameWithFullPath(C.GoString(oldINodePath), C.GoString(newINodePath))
	if err != nil {
		return libsdfs.CODE_ERR
	}
	return libsdfs.CODE_OK
}

//export GoSdfsGetPathInfo
func GoSdfsGetPathInfo(cInodePath *C.char) (inodeID uint64, size uint64, mTime uint64, code C.int) {
	var (
		fsINode   types.FsINode
		uNetINode types.NetINodeUintptr
		err       error
	)

	err = env.Client.MemDirTreeStg.FetchFsINodeByPath(C.GoString(cInodePath), &fsINode)
	if err != nil {
		return 0, 0, 0, libsdfs.CODE_ERR
	}

	uNetINode, err = env.Client.MemStg.NetINodeDriver.GetNetINodeWithReadAcquire(false, fsINode.NetINodeID)
	defer env.Client.MemStg.NetINodeDriver.ReleaseNetINodeWithReadRelease(uNetINode)
	if err != nil {
		return 0, 0, 0, libsdfs.CODE_ERR
	}

	return uint64(fsINode.Ino), uNetINode.Ptr().Size, uint64(fsINode.Mtime), libsdfs.CODE_OK
}
