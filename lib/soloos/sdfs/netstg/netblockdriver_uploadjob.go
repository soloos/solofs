package netstg

import (
	sdbapitypes "soloos/common/sdbapi/types"
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/types"
)

func (p *netBlockDriverUploader) PrepareUploadMemBlockJob(pJob *types.UploadMemBlockJob,
	uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr, netBlockIndex int32,
	uMemBlock types.MemBlockUintptr, memBlockIndex int32,
	backends snettypes.PeerGroup) {
	pJob.MetaDataState.LockContext()
	if pJob.MetaDataState.Load() == sdbapitypes.MetaDataStateInited {
		pJob.MetaDataState.UnlockContext()
		return
	}
	pJob.UNetINode = uNetINode
	pJob.UNetBlock = uNetBlock
	pJob.UMemBlock = uMemBlock
	pJob.MemBlockIndex = memBlockIndex

	pJob.UploadMaskWaitingIndex = 1
	pJob.UploadMaskSwap()

	pJob.MetaDataState.Store(sdbapitypes.MetaDataStateInited)
	pJob.MetaDataState.UnlockContext()
}
