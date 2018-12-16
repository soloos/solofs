package types

import (
	"soloos/util/offheap"
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
	SyncDataSig            sync.WaitGroup
	UploadPolicyMutex      sync.Mutex
	IsUploadPolicyPrepared bool
	UNetINode              NetINodeUintptr
	UNetBlock              NetBlockUintptr
	UMemBlock              MemBlockUintptr
	MemBlockIndex          int
	UploadMaskWaitingIndex int
	UploadMask             [2]offheap.ChunkMask
	UploadMaskWaiting      offheap.ChunkMaskUintptr
	UploadMaskProcessing   offheap.ChunkMaskUintptr
}

func (p *UploadMemBlockJob) Reset() {
	p.IsUploadPolicyPrepared = false
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
