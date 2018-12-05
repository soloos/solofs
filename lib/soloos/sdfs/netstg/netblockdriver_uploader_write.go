package netstg

import (
	"soloos/sdfs/types"
)

func (p *netBlockDriverUploader) PWrite(uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int,
	offset, end int) error {

	var (
		isMergeEventHappened    bool
		isMergeWriteMaskSuccess bool = false
	)

	pMemBlock := uMemBlock.Ptr()
	for isMergeWriteMaskSuccess == false {
		p.uploadMemBlockJobMutex.Lock()
		if pMemBlock.UploadJob.IsUploadPolicyPrepared == false {
			p.prepareUploadMemBlockJob(&pMemBlock.UploadJob, uNetBlock, uMemBlock, memBlockIndex)
		}

		isMergeEventHappened, isMergeWriteMaskSuccess =
			pMemBlock.UploadJob.UploadMaskWaiting.Ptr().MergeIncludeNeighbour(offset, end)

		if isMergeWriteMaskSuccess == true {
			if isMergeEventHappened == false {
				pMemBlock.UploadJob.UploadSig.Add(1)
				p.uploadMemBlockJobChan <- pMemBlock.GetUploadMemBlockJobUintptr()
			}
		}
		p.uploadMemBlockJobMutex.Unlock()

		if isMergeWriteMaskSuccess == false {
			pMemBlock.UploadJob.UploadSig.Wait()
		}
	}

	return nil
}

func (p *netBlockDriverUploader) FlushMemBlock(uMemBlock types.MemBlockUintptr) error {
	pMemBlock := uMemBlock.Ptr()
	// TODO add lock in metadb
	pMemBlock.UploadJob.UploadSig.Wait()
	return nil
}
