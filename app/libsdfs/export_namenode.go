package main

import "C"
import (
	fsapitypes "soloos/common/fsapi/types"
	"soloos/common/sdfsapi"
	sdfsapitypes "soloos/common/sdfsapi/types"
	"soloos/sdfs/types"
	"unsafe"
)

//export GoSdfsOpenFile
func GoSdfsOpenFile(cInodePath *C.char, flags,
	bufferSize C.int, replication C.short, blocksize C.long) (uint64, C.int) {
	var (
		fsINodePath = C.GoString(cInodePath)
		fsINodeMeta sdfsapitypes.FsINodeMeta
		err         error
	)

	fsINodeMeta, err = env.PosixFS.SimpleOpenFile(fsINodePath,
		types.DefaultNetBlockCap,
		env.Options.DefaultMemBlockCap)
	if err != nil {
		return 0, sdfsapi.CODE_ERR
	}

	return env.PosixFS.FdTableAllocFd(fsINodeMeta.Ino), sdfsapi.CODE_OK
}

//export GoSdfsExists
func GoSdfsExists(cInodePath *C.char) C.int {
	var (
		fsINodePath = C.GoString(cInodePath)
		fsINodeMeta sdfsapitypes.FsINodeMeta
		err         error
	)
	err = env.PosixFS.FetchFsINodeByPath(&fsINodeMeta, fsINodePath)
	if err != nil {
		// contains err == types.ErrObjectNotExists
		return sdfsapi.CODE_ERR
	}
	return sdfsapi.CODE_OK
}

//export GoSdfsListDirectory
func GoSdfsListDirectory(cInodePath *C.char, ret *unsafe.Pointer, num *C.int) {
	var (
		fsINodePath = C.GoString(cInodePath)
		arr         *[1<<30 - 1]*C.char
		index       int
		err         error
	)

	err = env.PosixFS.ListFsINodeByParentPath(fsINodePath, false,
		func(resultCount int) (uint64, uint64) {
			*ret = C.malloc(C.size_t(resultCount) * C.size_t(unsafe.Sizeof(uintptr(0))))
			*num = C.int(resultCount)
			arr = (*[1<<30 - 1]*C.char)(*ret)
			index = 0
			return uint64(resultCount), 0
		},
		func(fsINodeMeta sdfsapitypes.FsINodeMeta) bool {
			arr[index] = C.CString(fsINodeMeta.Name())
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
		code        fsapitypes.Status
	)
	code = env.PosixFS.SimpleMkdirAll(0777, fsINodePath, 0, 0)
	if code != fsapitypes.OK {
		return sdfsapi.CODE_ERR
	}

	return sdfsapi.CODE_OK
}

//export GoSdfsDelete
func GoSdfsDelete(cInodePath *C.char, recursive C.int) C.int {
	var (
		fsINodePath = C.GoString(cInodePath)
		err         error
	)
	err = env.PosixFS.DeleteFsINodeByPath(fsINodePath)
	if err != nil {
		return sdfsapi.CODE_ERR
	}

	return sdfsapi.CODE_OK
}

//export GoSdfsRename
func GoSdfsRename(oldINodePath, newINodePath *C.char) C.int {
	var err error
	err = env.PosixFS.RenameWithFullPath(C.GoString(oldINodePath), C.GoString(newINodePath))
	if err != nil {
		return sdfsapi.CODE_ERR
	}
	return sdfsapi.CODE_OK
}

//export GoSdfsGetPathInfo
func GoSdfsGetPathInfo(cInodePath *C.char) (inodeID uint64, size uint64, mTime uint64, code C.int) {
	var (
		fsINodeMeta sdfsapitypes.FsINodeMeta
		err         error
		status      fsapitypes.Status
	)

	err = env.PosixFS.FetchFsINodeByPath(&fsINodeMeta, C.GoString(cInodePath))
	if err != nil {
		return 0, 0, 0, sdfsapi.CODE_ERR
	}

	var (
		getAttrInput fsapitypes.GetAttrIn
		getAttrOut   fsapitypes.AttrOut
	)
	getAttrInput.NodeId = inodeID

	status = env.PosixFS.GetAttr(&getAttrInput, &getAttrOut)
	if status != fsapitypes.OK {
		return 0, 0, 0, sdfsapi.CODE_ERR
	}

	return uint64(inodeID), getAttrOut.Size, getAttrOut.Mtime, sdfsapi.CODE_OK
}
