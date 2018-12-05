package netstg

import (
	"soloos/sdfs/types"
)

func (p *netBlockDriverUploader) prepareUploadMemBlockJob(pUploadMemBlockJob *types.UploadMemBlockJob,
	uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int) {
	var i int
	pUploadMemBlockJob.UNetBlock = uNetBlock
	pUploadMemBlockJob.UMemBlock = uMemBlock
	pUploadMemBlockJob.MemBlockIndex = memBlockIndex

	// TODO add block placement policy
	for i = 0; i < uNetBlock.Ptr().DataNodes.Len; i++ {
		pUploadMemBlockJob.Backends.Append(uNetBlock.Ptr().DataNodes.Arr[i])
	}
	pUploadMemBlockJob.PrimaryBackendTransferCount = pUploadMemBlockJob.UNetBlock.Ptr().DataNodes.Len - 1
}
