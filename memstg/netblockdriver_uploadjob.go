package memstg

import (
	"soloos/common/snettypes"
	"soloos/common/solodbapitypes"
	"soloos/common/solofsapitypes"
)

func (p *netBlockDriverUploader) PrepareUploadMemBlockJob(pJob *solofsapitypes.UploadMemBlockJob,
	uNetINode solofsapitypes.NetINodeUintptr,
	uNetBlock solofsapitypes.NetBlockUintptr, netBlockIndex int32,
	uMemBlock solofsapitypes.MemBlockUintptr, memBlockIndex int32,
	backends snettypes.PeerGroup) {
	pJob.MetaDataState.LockContext()
	if pJob.MetaDataState.Load() == solodbapitypes.MetaDataStateInited {
		pJob.MetaDataState.UnlockContext()
		return
	}
	pJob.UNetINode = uNetINode
	pJob.UNetBlock = uNetBlock
	pJob.UMemBlock = uMemBlock
	pJob.MemBlockIndex = memBlockIndex

	pJob.UploadBlockJob.Reset()

	pJob.MetaDataState.Store(solodbapitypes.MetaDataStateInited)
	pJob.MetaDataState.UnlockContext()
}
