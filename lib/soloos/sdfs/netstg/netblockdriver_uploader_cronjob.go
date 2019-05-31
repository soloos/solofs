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
		uploadRetArr chan error
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

		uploadJobNum = pNetBlock.SyncDataBackends.Len - pNetBlock.SyncDataPrimaryBackendTransferCount
		uploadRetArr = make(chan error, uploadJobNum)

		// start upload
		// upload primary backend
		go func(uploadRetArr chan error, uJob sdfsapitypes.UploadMemBlockJobUintptr, transferCount int) {
			uploadRetArr <- p.driver.dataNodeClient.UploadMemBlock(uJob, 0, transferCount)
		}(uploadRetArr, uJob, pNetBlock.SyncDataPrimaryBackendTransferCount)

		// upload other backends
		for i = pNetBlock.SyncDataPrimaryBackendTransferCount + 1; i < pNetBlock.SyncDataBackends.Len; i++ {
			go func(uploadRetArr chan error, i int, uJob sdfsapitypes.UploadMemBlockJobUintptr) {
				uploadRetArr <- p.driver.dataNodeClient.UploadMemBlock(uJob, i, 0)
			}(uploadRetArr, i, uJob)
		}

		{
			var tmpErr error
			for i = 0; i < uploadJobNum; i++ {
				tmpErr = <-uploadRetArr
				if tmpErr != nil {
					err = tmpErr
				}
			}
		}

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
