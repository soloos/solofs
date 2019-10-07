package main

import "C"
import (
	"soloos/common/fsapi"
	"soloos/common/solofsapi"
	"soloos/common/solofstypes"
	"unsafe"
)

//export GoSolofsOpenFile
func GoSolofsOpenFile(cInodePath *C.char, flags,
	bufferSize C.int, replication C.short, blocksize C.long) (uint64, C.int) {
	var (
		fsINodePath = C.GoString(cInodePath)
		fsINodeMeta solofstypes.FsINodeMeta
		err         error
	)

	fsINodeMeta, err = env.solofsClient.SimpleOpenFile(fsINodePath,
		env.Options.DefaultNetBlockCap,
		env.Options.DefaultMemBlockCap)
	if err != nil {
		return 0, solofsapi.CODE_ERR
	}

	return env.solofsClient.FdTableAllocFd(fsINodeMeta.Ino), solofsapi.CODE_OK
}

//export GoSolofsExists
func GoSolofsExists(cInodePath *C.char) C.int {
	var (
		fsINodePath = C.GoString(cInodePath)
		fsINodeMeta solofstypes.FsINodeMeta
		err         error
	)
	err = env.solofsClient.FetchFsINodeByPath(&fsINodeMeta, fsINodePath)
	if err != nil {
		// contains err.Error() == solofstypes.ErrObjectNotExists.Error()
		return solofsapi.CODE_ERR
	}
	return solofsapi.CODE_OK
}

//export GoSolofsListDirectory
func GoSolofsListDirectory(cInodePath *C.char, ret *unsafe.Pointer, num *C.int) {
	var (
		fsINodePath = C.GoString(cInodePath)
		arr         *[1<<30 - 1]*C.char
		index       int
		err         error
	)

	err = env.solofsClient.ListFsINodeByParentPath(fsINodePath, false,
		func(resultCount int64) (uint64, uint64) {
			*ret = C.malloc(C.size_t(resultCount) * C.size_t(unsafe.Sizeof(uintptr(0))))
			*num = C.int(resultCount)
			arr = (*[1<<30 - 1]*C.char)(*ret)
			index = 0
			return uint64(resultCount), 0
		},
		func(fsINodeMeta solofstypes.FsINodeMeta) bool {
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

//export GoSolofsCreateDirectory
func GoSolofsCreateDirectory(cInodePath *C.char) C.int {
	var (
		fsINodePath = C.GoString(cInodePath)
		code        fsapi.Status
	)
	code = env.solofsClient.SimpleMkdirAll(0777, fsINodePath, 0, 0)
	if code != fsapi.OK {
		return solofsapi.CODE_ERR
	}

	return solofsapi.CODE_OK
}

//export GoSolofsDelete
func GoSolofsDelete(cInodePath *C.char, recursive C.int) C.int {
	var (
		fsINodePath = C.GoString(cInodePath)
		err         error
	)
	err = env.solofsClient.DeleteFsINodeByPath(fsINodePath)
	if err != nil {
		return solofsapi.CODE_ERR
	}

	return solofsapi.CODE_OK
}

//export GoSolofsRename
func GoSolofsRename(oldINodePath, newINodePath *C.char) C.int {
	var err error
	err = env.solofsClient.RenameWithFullPath(C.GoString(oldINodePath), C.GoString(newINodePath))
	if err != nil {
		return solofsapi.CODE_ERR
	}
	return solofsapi.CODE_OK
}

//export GoSolofsGetPathInfo
func GoSolofsGetPathInfo(cInodePath *C.char) (inodeID uint64, size uint64, mTime uint64, code C.int) {
	var (
		fsINodeMeta solofstypes.FsINodeMeta
		err         error
		status      fsapi.Status
	)

	err = env.solofsClient.FetchFsINodeByPath(&fsINodeMeta, C.GoString(cInodePath))
	if err != nil {
		return 0, 0, 0, solofsapi.CODE_ERR
	}

	var (
		getAttrInput fsapi.GetAttrIn
		getAttrOut   fsapi.AttrOut
	)
	getAttrInput.NodeId = inodeID

	status = env.solofsClient.GetAttr(&getAttrInput, &getAttrOut)
	if status != fsapi.OK {
		return 0, 0, 0, solofsapi.CODE_ERR
	}

	return uint64(inodeID), getAttrOut.Size, getAttrOut.Mtime, solofsapi.CODE_OK
}
