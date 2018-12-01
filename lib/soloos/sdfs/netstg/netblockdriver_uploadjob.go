package netstg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
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
	MemBlockIndex               int
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

func (p *netBlockDriverUploader) prepareUploadJob(uUploadJob UploadJobUintptr,
	uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int) {
	var (
		pUploadJob = uUploadJob.Ptr()
		i          int
	)
	pUploadJob.UNetBlock = uNetBlock
	pUploadJob.UMemBlock = uMemBlock
	pUploadJob.MemBlockIndex = memBlockIndex

	// TODO add block placement policy
	for i = 0; i < uNetBlock.Ptr().DataNodes.Len; i++ {
		pUploadJob.Backends.Append(uNetBlock.Ptr().DataNodes.Arr[i])
	}
	pUploadJob.PrimaryBackendTransferCount = pUploadJob.UNetBlock.Ptr().DataNodes.Len - 1
}
