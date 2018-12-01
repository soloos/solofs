package netstg

import "soloos/sdfs/types"

func (p *netBlockDriverUploader) PWrite(uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int,
	offset, end int) error {

	var (
		uUploadJob              UploadJobUintptr
		isMergeEventHappened    bool
		isMergeWriteMaskSuccess bool = false
	)

	uMemBlock.Ptr().UploadRWMutex.RLock()
	for isMergeWriteMaskSuccess == false {
		p.uploadJobMutex.Lock()
		uUploadJob, _ = p.uploadJobs[uMemBlock]
		if uUploadJob == 0 {
			uUploadJob = UploadJobUintptr(p.uploadJobPool.AllocRawObject())
			p.prepareUploadJob(uUploadJob, uNetBlock, uMemBlock, memBlockIndex)
			p.uploadJobs[uMemBlock] = uUploadJob
		}

		isMergeEventHappened, isMergeWriteMaskSuccess =
			uUploadJob.Ptr().UploadMaskWaiting.Ptr().MergeIncludeNeighbour(offset, end)

		if isMergeWriteMaskSuccess == true {
			if isMergeEventHappened == false {
				uUploadJob.Ptr().UNetBlock.Ptr().UploadSig.Add(1)
				p.uploadJobChan <- uUploadJob
			}
		}
		p.uploadJobMutex.Unlock()

		if isMergeWriteMaskSuccess == false {
			uUploadJob.Ptr().UNetBlock.Ptr().UploadSig.Wait()
		}
	}
	uMemBlock.Ptr().UploadRWMutex.RUnlock()

	return nil
}

func (p *netBlockDriverUploader) Flush(uMemBlock types.MemBlockUintptr) error {
	uMemBlock.Ptr().UploadRWMutex.Lock()
	// TODO add lock
	uUploadJob := p.uploadJobs[uMemBlock]
	if uUploadJob != 0 {
		uUploadJob.Ptr().UNetBlock.Ptr().UploadSig.Wait()
		delete(p.uploadJobs, uMemBlock)
		p.uploadJobPool.ReleaseRawObject(uintptr(uUploadJob))
	}
	uMemBlock.Ptr().UploadRWMutex.Unlock()
	return nil
}
