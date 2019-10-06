package main

import "C"
import (
	"soloos/common/fsapitypes"
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"unsafe"
)

//export GoSolofsOpenFile
func GoSolofsOpenFile(cInodePath *C.char, flags,
	bufferSize C.int, replication C.short, blocksize C.long) (uint64, C.int) {
	var (
		fsINodePath = C.GoString(cInodePath)
		fsINodeMeta solofsapitypes.FsINodeMeta
		err         error
	)

	fsINodeMeta, err = env.PosixFs.SimpleOpenFile(fsINodePath,
		env.Options.DefaultNetBlockCap,
		env.Options.DefaultMemBlockCap)
	if err != nil {
		return 0, solofsapi.CODE_ERR
	}

	return env.PosixFs.FdTableAllocFd(fsINodeMeta.Ino), solofsapi.CODE_OK
}

//export GoSolofsExists
func GoSolofsExists(cInodePath *C.char) C.int {
	var (
		fsINodePath = C.GoString(cInodePath)
		fsINodeMeta solofsapitypes.FsINodeMeta
		err         error
	)
	err = env.PosixFs.FetchFsINodeByPath(&fsINodeMeta, fsINodePath)
	if err != nil {
		// contains err.Error() == solofsapitypes.ErrObjectNotExists.Error()
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

	err = env.PosixFs.ListFsINodeByParentPath(fsINodePath, false,
		func(resultCount int) (uint64, uint64) {
			*ret = C.malloc(C.size_t(resultCount) * C.size_t(unsafe.Sizeof(uintptr(0))))
			*num = C.int(resultCount)
			arr = (*[1<<30 - 1]*C.char)(*ret)
			index = 0
			return uint64(resultCount), 0
		},
		func(fsINodeMeta solofsapitypes.FsINodeMeta) bool {
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
		code        fsapitypes.Status
	)
	code = env.PosixFs.SimpleMkdirAll(0777, fsINodePath, 0, 0)
	if code != fsapitypes.OK {
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
	err = env.PosixFs.DeleteFsINodeByPath(fsINodePath)
	if err != nil {
		return solofsapi.CODE_ERR
	}

	return solofsapi.CODE_OK
}

//export GoSolofsRename
func GoSolofsRename(oldINodePath, newINodePath *C.char) C.int {
	var err error
	err = env.PosixFs.RenameWithFullPath(C.GoString(oldINodePath), C.GoString(newINodePath))
	if err != nil {
		return solofsapi.CODE_ERR
	}
	return solofsapi.CODE_OK
}

//export GoSolofsGetPathInfo
func GoSolofsGetPathInfo(cInodePath *C.char) (inodeID uint64, size uint64, mTime uint64, code C.int) {
	var (
		fsINodeMeta solofsapitypes.FsINodeMeta
		err         error
		status      fsapitypes.Status
	)

	err = env.PosixFs.FetchFsINodeByPath(&fsINodeMeta, C.GoString(cInodePath))
	if err != nil {
		return 0, 0, 0, solofsapi.CODE_ERR
	}

	var (
		getAttrInput fsapitypes.GetAttrIn
		getAttrOut   fsapitypes.AttrOut
	)
	getAttrInput.NodeId = inodeID

	status = env.PosixFs.GetAttr(&getAttrInput, &getAttrOut)
	if status != fsapitypes.OK {
		return 0, 0, 0, solofsapi.CODE_ERR
	}

	return uint64(inodeID), getAttrOut.Size, getAttrOut.Mtime, solofsapi.CODE_OK
}
