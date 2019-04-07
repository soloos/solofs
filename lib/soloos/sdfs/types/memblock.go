package types

import (
	"reflect"
	"soloos/common/log"
	snettypes "soloos/common/snet/types"
	"soloos/sdbone/offheap"
	"sync"
	"unsafe"
)

const (
	MemBlockStructSize = unsafe.Sizeof(MemBlock{})

	MemBlockRefuseReleaseForErr = -1
)

type MemBlockUintptr uintptr

func (u MemBlockUintptr) Ptr() *MemBlock {
	return (*MemBlock)(unsafe.Pointer(u))
}

type MemBlock struct {
	offheap.HKVTableObjectWithBytes12
	RebaseNetBlockMutex sync.Mutex
	Bytes               reflect.SliceHeader
	AvailMask           offheap.ChunkMask
	UploadJob           UploadMemBlockJob
}

func (p *MemBlock) Contains(offset, end int) bool {
	return p.AvailMask.Contains(offset, end)
}

func (p *MemBlock) PWriteWithConn(conn *snettypes.Connection, length int, offset int) (isSuccess bool) {
	_, isSuccess = p.AvailMask.MergeIncludeNeighbour(offset, offset+length)
	if isSuccess {
		var err error
		if offset+length > p.Bytes.Cap {
			length = p.Bytes.Cap - offset
		}
		bytes := (*(*[]byte)(unsafe.Pointer(&p.Bytes)))
		err = conn.ReadAll(bytes[offset : offset+length])
		if err != nil {
			log.Warn("PWriteWithConn error", err)
			isSuccess = false
		}
	}
	return
}

func (p *MemBlock) PWriteWithMem(data []byte, offset int) (isSuccess bool) {
	_, isSuccess = p.AvailMask.MergeIncludeNeighbour(offset, offset+len(data))
	if isSuccess {
		copy((*(*[]byte)(unsafe.Pointer(&p.Bytes)))[offset:], data)
	}
	return
}

func (p *MemBlock) PReadWithConn(conn *snettypes.Connection, length int, offset int) error {
	var err error
	err = conn.WriteAll((*(*[]byte)(unsafe.Pointer(&p.Bytes)))[offset : offset+length])
	if err != nil {
		return err
	}
	return nil
}

func (p *MemBlock) PReadWithMem(data []byte, offset int) {
	copy(data, (*(*[]byte)(unsafe.Pointer(&p.Bytes)))[offset:])
}

func (p *MemBlock) GetUploadMemBlockJobUintptr() UploadMemBlockJobUintptr {
	return UploadMemBlockJobUintptr(unsafe.Pointer(p)) + UploadMemBlockJobUintptr(unsafe.Offsetof(p.UploadJob))
}

func (p *MemBlock) BytesSlice() *[]byte {
	return (*[]byte)(unsafe.Pointer(&p.Bytes))
}

func (p *MemBlock) Reset() {
	p.AvailMask.Reset()
	p.UploadJob.Reset()
	p.HKVTableObjectWithBytes12.Reset()
}
