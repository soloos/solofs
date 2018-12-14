package netstg

import (
	"soloos/sdfs/types"
	"soloos/util"
)

func (p *netBlockDriverUploader) cronUpload() error {
	var (
		uJob      types.UploadMemBlockJobUintptr
		pJob      *types.UploadMemBlockJob
		pNetINode *types.NetINode
		i         int
		ok        bool
		err       error
	)

	for {
		uJob, ok = <-p.uploadMemBlockJobChan
		if ok == false {
			panic("uploadMemBlockJobChan closed")
		}

		pJob = uJob.Ptr()
		pNetINode = pJob.UNetINode.Ptr()

		// prepare upload job
		pJob.UploadPolicyMutex.Lock()
		if pJob.UploadMaskWaiting.Ptr().MaskArrayLen == 0 {
			// upload done and continue
			pJob.UploadPolicyMutex.Unlock()
			goto ONE_RUN_DONE
		}
		pJob.UploadMaskSwap()
		pJob.UploadPolicyMutex.Unlock()

		util.AssertTrue(pJob.Backends.Len > 0)

		// start upload
		// upload primary backend
		err = p.driver.dataNodeClient.UploadMemBlock(uJob, 0, pJob.PrimaryBackendTransferCount)

		// upload other backends
		for i = pJob.PrimaryBackendTransferCount + 1; i < pJob.Backends.Len; i++ {
			err = p.driver.dataNodeClient.UploadMemBlock(uJob, i, 0)
		}

	ONE_RUN_DONE:
		pJob.SyncDataSig.Done()
		pNetINode.SyncDataSig.Done()

		if pNetINode.LastSyncDataError != err {
			err = pNetINode.LastSyncDataError
		}

		if err == nil {
			pJob.UploadMaskProcessing.Ptr().Reset()
		} else {
			// TODO catch error
			pNetINode.LastSyncDataError = err
		}
	}

	return nil
}
