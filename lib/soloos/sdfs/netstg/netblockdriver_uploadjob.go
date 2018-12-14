package netstg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

func (p *netBlockDriverUploader) SetUploadMemBlockJobBackends(pJob *types.UploadMemBlockJob,
	backends snettypes.PeerUintptrArray8) {
	// TODO add block placement policy
	var i int
	pJob.Backends.Reset()
	for i = 0; i < backends.Len; i++ {
		pJob.Backends.Append(backends.Arr[i])
	}
	pJob.PrimaryBackendTransferCount = backends.Len - 1
}

func (p *netBlockDriverUploader) prepareUploadMemBlockJob(pJob *types.UploadMemBlockJob,
	uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int) {
	pJob.UploadPolicyMutex.Lock()
	if pJob.IsUploadPolicyPrepared {
		pJob.UploadPolicyMutex.Unlock()
		return
	}
	pJob.UNetINode = uNetINode
	pJob.UNetBlock = uNetBlock
	pJob.UMemBlock = uMemBlock
	pJob.MemBlockIndex = memBlockIndex

	pJob.UploadMaskWaitingIndex = 1
	pJob.UploadMaskSwap()

	p.SetUploadMemBlockJobBackends(pJob, uNetBlock.Ptr().DataNodes)
	pJob.IsUploadPolicyPrepared = true
	pJob.UploadPolicyMutex.Unlock()
}
