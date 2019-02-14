package types

import (
	"reflect"
	"unsafe"
)

const (
	NetINodeBlockIDSize int = NetINodeIDSize + Int32Size
	PtrBindIndexSize    int = UintptrSize + Int32Size
)

type NetINodeBlockID = [NetINodeBlockIDSize]byte
type PtrBindIndex = [PtrBindIndexSize]byte

func EncodeNetINodeBlockID(netINodeBlockID *NetINodeBlockID, netINodeID NetINodeID, blockIndex int32) {
	bytes := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		uintptr(unsafe.Pointer(netINodeBlockID)),
		NetINodeBlockIDSize,
		NetINodeBlockIDSize,
	}))
	copy(bytes[:NetINodeIDSize], (*(*[NetINodeIDSize]byte)((unsafe.Pointer)(&netINodeID)))[:NetINodeIDSize])
	copy(bytes[NetINodeIDSize:], (*(*[Int32Size]byte)((unsafe.Pointer)(&blockIndex)))[:Int32Size])
}

func EncodePtrBindIndex(id *PtrBindIndex, u uintptr, index int32) {
	*((*uintptr)(unsafe.Pointer(id))) = u
	*((*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(id)) + UUintptrSize))) = index
}
