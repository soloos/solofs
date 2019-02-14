package types

import (
	"reflect"
	"unsafe"
)

const (
	UintptrSize = int(unsafe.Sizeof(uintptr(0)))
	Int32Size   = int(unsafe.Sizeof(int32(0)))

	UUintptrSize = unsafe.Sizeof(uintptr(0))
)

type BytesUintptr uintptr

func (u BytesUintptr) MaxIntBytes() *[MaxInt]byte {
	return (*[MaxInt]byte)(unsafe.Pointer(u))
}

func (u BytesUintptr) Ptr() *[]byte { return (*[]byte)(unsafe.Pointer(u)) }

func (u *BytesUintptr) Change(addr uintptr, len, cap int) {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(u))
	sliceHeader.Data = addr
	sliceHeader.Len = len
	sliceHeader.Cap = cap
}
