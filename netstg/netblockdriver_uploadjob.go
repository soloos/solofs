package netstg

import (
	"soloos/common/sdbapitypes"
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
)

func (p *netBlockDriverUploader) PrepareUploadMemBlockJob(pJob *sdfsapitypes.UploadMemBlockJob,
	uNetINode sdfsapitypes.NetINodeUintptr,
	uNetBlock sdfsapitypes.NetBlockUintptr, netBlockIndex int32,
	uMemBlock sdfsapitypes.MemBlockUintptr, memBlockIndex int32,
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

	pJob.UploadBlockJob.Reset()

	pJob.MetaDataState.Store(sdbapitypes.MetaDataStateInited)
	pJob.MetaDataState.UnlockContext()
}
