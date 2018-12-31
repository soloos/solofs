package types

import "unsafe"

type FsDirUintptr uintptr

func (u FsDirUintptr) Ptr() *FsDir { return (*FsDir)(unsafe.Pointer(u)) }

type FsDir struct {
}
