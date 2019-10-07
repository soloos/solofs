package memstg

import (
	"soloos/common/snet"
	"soloos/common/solodbtypes"
	"soloos/common/solofstypes"
)

func (p *netBlockDriverUploader) PrepareUploadMemBlockJob(pJob *solofstypes.UploadMemBlockJob,
	uNetINode solofstypes.NetINodeUintptr,
	uNetBlock solofstypes.NetBlockUintptr, netBlockIndex int32,
	uMemBlock solofstypes.MemBlockUintptr, memBlockIndex int32,
	backends snet.PeerGroup) {
	pJob.MetaDataState.LockContext()
	if pJob.MetaDataState.Load() == solodbtypes.MetaDataStateInited {
		pJob.MetaDataState.UnlockContext()
		return
	}
	pJob.UNetINode = uNetINode
	pJob.UNetBlock = uNetBlock
	pJob.UMemBlock = uMemBlock
	pJob.MemBlockIndex = memBlockIndex

	pJob.UploadBlockJob.Reset()

	pJob.MetaDataState.Store(solodbtypes.MetaDataStateInited)
	pJob.MetaDataState.UnlockContext()
}
