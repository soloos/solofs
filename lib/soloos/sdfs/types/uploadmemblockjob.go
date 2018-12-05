package types

import (
	snettypes "soloos/snet/types"
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
	IsUploadPolicyPrepared      bool
	UNetBlock                   NetBlockUintptr
	UMemBlock                   MemBlockUintptr
	UploadSig                   sync.WaitGroup
	MemBlockIndex               int
	UploadMaskWaitingIndex      int
	UploadMask                  [2]offheap.ChunkMask
	UploadMaskWaiting           offheap.ChunkMaskUintptr
	UploadMaskProcessing        offheap.ChunkMaskUintptr
	PrimaryBackendTransferCount int
	Backends                    snettypes.PeerUintptrArray8
}

func (p *UploadMemBlockJob) Reset() {
	p.UploadMaskWaitingIndex = 1
	p.UploadMaskSwap()
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
