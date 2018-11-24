package netstg

import (
	"soloos/sdfs/types"
	"soloos/util/offheap"
	"unsafe"
)

const (
	UploadJobStructSize = unsafe.Sizeof(UploadJob{})
)

type UploadJobUintptr uintptr

func (u UploadJobUintptr) Ptr() *UploadJob {
	return (*UploadJob)(unsafe.Pointer(u))
}

type UploadJob struct {
	UNetBlock                   types.NetBlockUintptr
	UMemBlock                   types.MemBlockUintptr
	UploadMaskWaiting           offheap.ChunkMaskUintptr
	UploadMaskProcessing        offheap.ChunkMaskUintptr
	PrimaryBackendTransferCount int
	Backends                    snettypes.PeerUintptrArray8
}

func (p *netBlockDriverUploader) RawChunkPoolInvokePrepareNewUploadJob(uRawChunk uintptr) {
	uUploadJob := UploadJobUintptr(uRawChunk)
	uUploadJob.Ptr().UploadMaskWaiting = offheap.ChunkMaskUintptr(p.uploadChunkMaskPool.AllocRawObject())
	uUploadJob.Ptr().UploadMaskProcessing = offheap.ChunkMaskUintptr(p.uploadChunkMaskPool.AllocRawObject())
}
