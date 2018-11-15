package types

import (
	snettypes "soloos/snet/types"
	"unsafe"
)

const (
	NetBlockStructSize = unsafe.Sizeof(NetBlock{})
)

type NetBlockUintptr uintptr

func (u NetBlockUintptr) Ptr() *NetBlock { return (*NetBlock)(unsafe.Pointer(u)) }

type NetBlock struct {
	DataNodes snettypes.PeerUintptrArray8
}
