package types

import (
	"reflect"
	"soloos/util/offheap"
	"sync"
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
	MemID         PtrBindIndex
	Status        int64 // equals 0 if could be release
	UploadRWMutex sync.RWMutex
	Chunk         offheap.ChunkUintptr
	Bytes         reflect.SliceHeader
	AvailMask     offheap.ChunkMask
}

func (p *MemBlock) Contains(offset, end int) bool {
	return p.AvailMask.Contains(offset, end)
}

func (p *MemBlock) PWrite(data []byte, offset int) (isSuccess bool) {
	_, isSuccess = p.AvailMask.MergeIncludeNeighbour(offset, offset+len(data))
	if isSuccess {
		copy((*(*[]byte)(unsafe.Pointer(&p.Bytes)))[offset:], data)
	}
	return
}

func (p *MemBlock) BytesSlice() *[]byte {
	return (*[]byte)(unsafe.Pointer(&p.Bytes))
}

func (p *MemBlock) Reset() {
	p.Status = MemBlockUninited
	p.AvailMask.Reset()
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
