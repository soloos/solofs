package types

import (
	"unsafe"
)

const (
	INodeStructSize = unsafe.Sizeof(INode{})
)

type INodeUintptr uintptr

func (p INodeUintptr) Ptr() *INode { return (*INode)(unsafe.Pointer(p)) }

type INode struct {
	ID           INodeID
	Size         uint64
	NetBlockSize int
	MemBlockSize int
	INodeSize     int
}
