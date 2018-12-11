package types

import (
	"reflect"
	"unsafe"
)

const (
	NetINodeBlockIDSize int = NetINodeIDSize + IntSize
	PtrBindIndexSize int = UintptrSize + IntSize
)

type NetINodeBlockID = [NetINodeBlockIDSize]byte
type PtrBindIndex = [PtrBindIndexSize]byte

func EncodeNetINodeBlockID(netINodeBlockID *NetINodeBlockID, netINodeID NetINodeID, blockIndex int) {
	bytes := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		uintptr(unsafe.Pointer(netINodeBlockID)),
		NetINodeBlockIDSize,
		NetINodeBlockIDSize,
	}))
	copy(bytes[:NetINodeIDSize], (*(*[NetINodeIDSize]byte)((unsafe.Pointer)(&netINodeID)))[:NetINodeIDSize])
	copy(bytes[NetINodeIDSize:], (*(*[IntSize]byte)((unsafe.Pointer)(&blockIndex)))[:IntSize])
}

func EncodePtrBindIndex(id *PtrBindIndex, u uintptr, index int) {
	*((*uintptr)(unsafe.Pointer(id))) = u
	*((*int)(unsafe.Pointer(uintptr(unsafe.Pointer(id)) + UUintptrSize))) = index
}
