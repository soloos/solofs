package memstg

import (
	"soloos/common/solofsapitypes"
	"soloos/common/util"
)

func (p *netBlockDriverUploader) cronUpload() error {
	var (
		uJob         solofsapitypes.UploadMemBlockJobUintptr
		pJob         *solofsapitypes.UploadMemBlockJob
		pNetINode    *solofsapitypes.NetINode
		pNetBlock    *solofsapitypes.NetBlock
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

		uploadJobNum = 0
		for i = 0; i < pNetBlock.SyncDataBackends.Len; {
			i += int(pNetBlock.SyncDataBackends.Arr[i].TransferCount + 1)
			uploadJobNum++
		}

		uploadRetArr = make(chan error, uploadJobNum)

		for i = 0; i < pNetBlock.SyncDataBackends.Len; {
			go func(uploadRetArr chan error, i int, uJob solofsapitypes.UploadMemBlockJobUintptr) {
				uploadRetArr <- p.driver.solodnClient.UploadMemBlock(uJob, i)
			}(uploadRetArr, i, uJob)
			i += int(pNetBlock.SyncDataBackends.Arr[i].TransferCount + 1)
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
		if err != nil {
			// TODO catch error
			pNetINode.LastSyncDataError = err
		} else {
			pJob.ResetProcessingChunkMask()
		}

		pJob.SyncDataSig.Done()
		pNetINode.SyncDataSig.Done()
	}

	return nil
}
