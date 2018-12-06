package netstg

import (
	"soloos/sdfs/types"
	"soloos/util/offheap"
)

func (p *netBlockDriverUploader) cronUpload() error {
	var (
		uUploadMemBlockJob types.UploadMemBlockJobUintptr
		pUploadMemBlockJob *types.UploadMemBlockJob
		pChunkMask         *offheap.ChunkMask
		i                  int
		ok                 bool
		err                error
	)

	for {
		uUploadMemBlockJob, ok = <-p.uploadMemBlockJobChan
		if ok == false {
			panic("uploadMemBlockJobChan closed")
		}

		pUploadMemBlockJob = uUploadMemBlockJob.Ptr()

		p.uploadMemBlockJobMutex.Lock()
		if pUploadMemBlockJob.UploadMaskWaiting.Ptr().MaskArrayLen == 0 {
			// upload done and continue
			pUploadMemBlockJob.UploadSig.Done()
			p.uploadMemBlockJobMutex.Unlock()
			continue
		}

		// start upload
		pUploadMemBlockJob.UploadMaskSwap()
		p.uploadMemBlockJobMutex.Unlock()

		pChunkMask = pUploadMemBlockJob.UploadMaskProcessing.Ptr()

		// upload primary backend
		{
			if pUploadMemBlockJob.PrimaryBackendTransferCount > 0 {
				err = p.driver.dataNodeClient.UploadMemBlock(uUploadMemBlockJob,
					pUploadMemBlockJob.Backends.Arr[0],
					pUploadMemBlockJob.Backends.Arr[1:1+pUploadMemBlockJob.PrimaryBackendTransferCount])
			} else {
				err = p.driver.dataNodeClient.UploadMemBlock(uUploadMemBlockJob,
					pUploadMemBlockJob.Backends.Arr[0],
					nil)
			}
		}

		// upload other backends
		for i = pUploadMemBlockJob.PrimaryBackendTransferCount + 1; i < pUploadMemBlockJob.Backends.Len; i++ {
			err = p.driver.dataNodeClient.UploadMemBlock(uUploadMemBlockJob,
				pUploadMemBlockJob.Backends.Arr[i],
				nil)
		}

		pUploadMemBlockJob.UploadSig.Done()

		// TODO catch error
		if err != nil {
			return err
		}

		pChunkMask.Reset()
	}

	return nil
}
