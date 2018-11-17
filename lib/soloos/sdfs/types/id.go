package types

import (
	"reflect"
	"unsafe"
)

const (
	INodeBlockIDSize int = INodeIDSize + IntSize
	PtrBindIndexSize int = UintptrSize + IntSize
)

type INodeBlockID = [INodeBlockIDSize]byte
type PtrBindIndex = [PtrBindIndexSize]byte

func EncodeINodeBlockID(inodeBlockID *INodeBlockID, inodeID INodeID, blockIndex int) {
	bytes := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		uintptr(unsafe.Pointer(inodeBlockID)),
		INodeBlockIDSize,
		INodeBlockIDSize,
	}))
	copy(bytes[:INodeIDSize], (*(*[INodeIDSize]byte)((unsafe.Pointer)(&inodeID)))[:INodeIDSize])
	copy(bytes[INodeIDSize:], (*(*[IntSize]byte)((unsafe.Pointer)(&blockIndex)))[:IntSize])
}

func EncodePtrBindIndex(id *PtrBindIndex, u uintptr, index int) {
	*((*uintptr)(unsafe.Pointer(id))) = u
	*((*int)(unsafe.Pointer(uintptr(unsafe.Pointer(id)) + UUintptrSize))) = index
}
