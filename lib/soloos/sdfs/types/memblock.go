package types

import (
	"reflect"
	"soloos/util/offheap"
	"sync/atomic"
	"unsafe"
)

const (
	MemBlockStructSize = unsafe.Sizeof(MemBlock{})

	MemBlockRefuseReleaseForErr = -1

	MemBlockUninited = iota
	MemBlockIniteded
	MemBlockReleasable
	MemBlockRelease
)

type MemBlockUintptr uintptr

func (u MemBlockUintptr) Ptr() *MemBlock {
	return (*MemBlock)(unsafe.Pointer(u))
}

type MemBlock struct {
	MemID  PtrBindIndex
	Status int64 // equals 0 if could be release
	Chunk  offheap.ChunkUintptr
	Bytes  reflect.SliceHeader
}

func (p *MemBlock) BytesSlice() *[]byte {
	return (*[]byte)(unsafe.Pointer(&p.Bytes))
}

func (p *MemBlock) Reset() {
	p.Status = MemBlockUninited
}

func (p *MemBlock) CompleteInit() {
	p.Status = MemBlockIniteded
}

func (p *MemBlock) IsInited() bool {
	return p.Status > MemBlockUninited
}

func (p *MemBlock) SetReleasable() {
	p.Status = MemBlockReleasable
}

func (p *MemBlock) EnsureRelease() bool {
	return atomic.CompareAndSwapInt64(&p.Status, MemBlockReleasable, MemBlockRelease)
}
