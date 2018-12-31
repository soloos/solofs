package types

import "unsafe"

type FsFileUintptr uintptr

func (u FsFileUintptr) Ptr() *FsFile { return (*FsFile)(unsafe.Pointer(u)) }

type FsFile struct {
}
