package types

import (
	"soloos/sdbone/offheap"
	"sync"
	"unsafe"
)

const (
	UploadMemBlockJobStructSize = unsafe.Sizeof(UploadMemBlockJob{})
)

type UploadMemBlockJobUintptr uintptr

func (u UploadMemBlockJobUintptr) Ptr() *UploadMemBlockJob {
	return (*UploadMemBlockJob)(unsafe.Pointer(u))
}

type UploadMemBlockJob struct {
	MetaDataStateMutex     sync.Mutex
	MetaDataState          MetaDataState
	SyncDataSig            sync.WaitGroup
	UNetINode              NetINodeUintptr
	UNetBlock              NetBlockUintptr
	UMemBlock              MemBlockUintptr
	MemBlockIndex          int32
	UploadMaskWaitingIndex int
	UploadMask             [2]offheap.ChunkMask
	UploadMaskWaiting      offheap.ChunkMaskUintptr
	UploadMaskProcessing   offheap.ChunkMaskUintptr
}

func (p *UploadMemBlockJob) Reset() {
	p.MetaDataState.Store(MetaDataStateUninited)
}

func (p *UploadMemBlockJob) UploadMaskSwap() {
	if p.UploadMaskWaitingIndex == 0 {
		p.UploadMaskWaiting = offheap.ChunkMaskUintptr(unsafe.Pointer(&p.UploadMask[1]))
		p.UploadMaskProcessing = offheap.ChunkMaskUintptr(unsafe.Pointer(&p.UploadMask[0]))
		p.UploadMaskWaitingIndex = 1
	} else {
		p.UploadMaskWaiting = offheap.ChunkMaskUintptr(unsafe.Pointer(&p.UploadMask[0]))
		p.UploadMaskProcessing = offheap.ChunkMaskUintptr(unsafe.Pointer(&p.UploadMask[1]))
		p.UploadMaskWaitingIndex = 0
	}
}
