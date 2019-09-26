package memstg

import (
	"soloos/common/solodbapitypes"
	"soloos/common/solofsapitypes"
)

func (p *netBlockDriverUploader) PWrite(uNetINode solofsapitypes.NetINodeUintptr,
	uNetBlock solofsapitypes.NetBlockUintptr, netBlockIndex int32,
	uMemBlock solofsapitypes.MemBlockUintptr, memBlockIndex int32,
	offset, end int) error {

	var (
		isMergeEventHappened    bool
		isMergeWriteMaskSuccess bool = false
		pMemBlock                    = uMemBlock.Ptr()
		pUploadJob                   = &pMemBlock.UploadJob
	)

	pUploadJob = pMemBlock.GetUploadMemBlockJobUintptr().Ptr()

	if pUploadJob.MetaDataState.Load() == solodbapitypes.MetaDataStateUninited {
		// TODO: refine me
		p.PrepareUploadMemBlockJob(pUploadJob,
			uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, uNetBlock.Ptr().StorDataBackends)
	}

	for isMergeWriteMaskSuccess == false {
		isMergeEventHappened, isMergeWriteMaskSuccess = pUploadJob.WaitingQueueMergeIncludeNeighbour(offset, end)

		if isMergeWriteMaskSuccess {
			if isMergeEventHappened == false {
				pUploadJob.UNetINode.Ptr().SyncDataSig.Add(1)
				pUploadJob.SyncDataSig.Add(1)
				p.uploadMemBlockJobChan <- pMemBlock.GetUploadMemBlockJobUintptr()
			}
		}

		if isMergeWriteMaskSuccess == false {
			pMemBlock.UploadJob.SyncDataSig.Wait()
		}
	}

	return nil
}
