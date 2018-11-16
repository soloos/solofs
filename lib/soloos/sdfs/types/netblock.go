package types

import (
	snettypes "soloos/snet/types"
	"sync"
	"unsafe"
)

const (
	NetBlockStructSize            = unsafe.Sizeof(NetBlock{})
	MaxDataNodesSizeStoreNetBlock = 8
)

type NetBlockUintptr uintptr

func (u NetBlockUintptr) Ptr() *NetBlock { return (*NetBlock)(unsafe.Pointer(u)) }

type NetBlock struct {
	UploadSig sync.WaitGroup
	DataNodes snettypes.PeerUintptrArray8
}
