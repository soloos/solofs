package netstg

import (
	"soloos/common/util"
	"soloos/sdfs/types"
	"sync"
)

func (p *netBlockDriverUploader) cronUpload() error {
	var (
		uJob         types.UploadMemBlockJobUintptr
		pJob         *types.UploadMemBlockJob
		pNetINode    *types.NetINode
		pNetBlock    *types.NetBlock
		uploadJobNum int
		uploadWG     sync.WaitGroup
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
		go func() {
			uploadErrors[0] = p.driver.dataNodeClient.UploadMemBlock(uJob, 0, pNetBlock.SyncDataPrimaryBackendTransferCount)
			uploadWG.Done()
		}()

		// upload other backends
		for i = pNetBlock.SyncDataPrimaryBackendTransferCount + 1; i < pNetBlock.SyncDataBackends.Len; i++ {
			go func(i int) {
				uploadErrors[i] = p.driver.dataNodeClient.UploadMemBlock(uJob, i, 0)
				uploadWG.Done()
			}(i)
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
