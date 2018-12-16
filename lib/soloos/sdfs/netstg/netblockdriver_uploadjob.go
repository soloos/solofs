package netstg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

func (p *netBlockDriverUploader) PrepareUploadMemBlockJob(pJob *types.UploadMemBlockJob,
	uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr, netBlockIndex int,
	uMemBlock types.MemBlockUintptr, memBlockIndex int,
	backends snettypes.PeerUintptrArray8) {
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

	pJob.IsUploadPolicyPrepared = true
	pJob.UploadPolicyMutex.Unlock()
}
