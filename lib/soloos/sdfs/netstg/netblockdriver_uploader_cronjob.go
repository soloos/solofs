package netstg

import (
	"soloos/sdfs/types"
	"soloos/util"
	"soloos/util/offheap"
)

func (p *netBlockDriverUploader) cronUpload() error {
	var (
		uJob       types.UploadMemBlockJobUintptr
		pJob       *types.UploadMemBlockJob
		pNetINode  *types.NetINode
		pChunkMask *offheap.ChunkMask
		i          int
		ok         bool
		err        error
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

		// start upload

		pChunkMask = pJob.UploadMaskProcessing.Ptr()

		util.AssertTrue(pJob.Backends.Len > 0)

		// upload primary backend
		if pJob.PrimaryBackendTransferCount > 0 {
			err = p.driver.dataNodeClient.UploadMemBlock(uJob, 0, pJob.PrimaryBackendTransferCount)
		} else {
			err = p.driver.dataNodeClient.UploadMemBlock(uJob, 0, 0)
		}

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
			pChunkMask.Reset()
		} else {
			// TODO catch error
			pNetINode.LastSyncDataError = err
		}
	}

	return nil
}
