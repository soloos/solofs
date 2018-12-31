package types

import (
	"unsafe"
)

type FsINodeID = int64

type FsINodeUintptr uintptr

func (u FsINodeUintptr) Ptr() *FsINode { return (*FsINode)(unsafe.Pointer(u)) }

type FsINode struct {
	ID         FsINodeID
	ParentID   FsINodeID
	Name       string
	Flag       int
	Permission int
	NetINodeID NetINodeID
	Type       int
}
