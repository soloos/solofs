package netstg

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/util"
)

func (p *netBlockDriverUploader) cronUpload() error {
	var (
		uJob         sdfsapitypes.UploadMemBlockJobUintptr
		pJob         *sdfsapitypes.UploadMemBlockJob
		pNetINode    *sdfsapitypes.NetINode
		pNetBlock    *sdfsapitypes.NetBlock
		uploadJobNum int
		uploadWG     util.RawWaitGroup
		uploadErrors []error
		i            int
		ok           bool
		err          error
	)

	for {
		uJob, ok = <-p.uploadMemBlockJobChan
		if ok == false {
			panic("uploadMemBlockJobChan closed")
		}

		pJob = uJob.Ptr()
		pNetINode = pJob.UNetINode.Ptr()
		pNetBlock = pJob.UNetBlock.Ptr()

		if pJob.PrepareUploadMask() {
			goto ONE_RUN_DONE
		}

		util.AssertTrue(pNetBlock.SyncDataBackends.Len > 0)

		uploadErrors = uploadErrors[:0]
		uploadJobNum = pNetBlock.SyncDataBackends.Len - pNetBlock.SyncDataPrimaryBackendTransferCount
		for i = 0; i < uploadJobNum; i++ {
			uploadErrors = append(uploadErrors, nil)
		}

		uploadWG.Add(uploadJobNum)

		// start upload
		// upload primary backend
		go func(uJob sdfsapitypes.UploadMemBlockJobUintptr, transferCount int) {
			uploadErrors[0] = p.driver.dataNodeClient.UploadMemBlock(uJob, 0, transferCount)
			uploadWG.Done()
		}(uJob, pNetBlock.SyncDataPrimaryBackendTransferCount)

		// upload other backends
		for i = pNetBlock.SyncDataPrimaryBackendTransferCount + 1; i < pNetBlock.SyncDataBackends.Len; i++ {
			go func(i int, uJob sdfsapitypes.UploadMemBlockJobUintptr) {
				uploadErrors[i] = p.driver.dataNodeClient.UploadMemBlock(uJob, i, 0)
				uploadWG.Done()
			}(i, uJob)
		}

		for i, _ = range uploadErrors {
			if uploadErrors[i] != nil {
				err = uploadErrors[i]
				break
			}
		}

		uploadWG.Wait()

	ONE_RUN_DONE:
		pJob.SyncDataSig.Done()
		pNetINode.SyncDataSig.Done()

		if err != nil {
			// TODO catch error
			pNetINode.LastSyncDataError = err
		} else {
			pJob.ResetProcessingChunkMask()
		}
	}

	return nil
}
