package netstg

import (
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"
	"sync"
)

type netBlockDriverUploader struct {
	driver *NetBlockDriver

	snetDriver       *snet.SNetDriver
	snetClientDriver *snet.ClientDriver

	uploadJobMutex      sync.Mutex
	uploadJobChan       chan UploadJobUintptr
	uploadJobs          map[types.MemBlockUintptr]UploadJobUintptr
	uploadChunkMaskPool offheap.RawObjectPool
	uploadJobPool       offheap.RawObjectPool
}

func (p *netBlockDriverUploader) Init(driver *NetBlockDriver,
	snetDriver *snet.SNetDriver,
	snetClientDriver *snet.ClientDriver) error {
	var err error
	p.driver = driver

	p.snetDriver = snetDriver
	p.snetClientDriver = snetClientDriver

	p.uploadJobChan = make(chan UploadJobUintptr, 2048)
	p.uploadJobs = make(map[types.MemBlockUintptr]UploadJobUintptr, 2048)
	err = p.driver.offheapDriver.InitRawObjectPool(&p.uploadChunkMaskPool,
		int(offheap.ChunkMaskStructSize), -1, nil, nil)
	if err != nil {
		return err
	}
	err = p.driver.offheapDriver.InitRawObjectPool(&p.uploadJobPool,
		int(UploadJobStructSize), -1, p.RawChunkPoolInvokePrepareNewUploadJob, nil)
	if err != nil {
		return err
	}

	go func() {
		util.AssertErrIsNil(p.cronUpload())
	}()

	return nil
}

func (p *netBlockDriverUploader) doUpload(uploadJobErr *error, uploadJobSig *sync.WaitGroup,
	uUploadJob UploadJobUintptr,
	uPeer snettypes.PeerUintptr,
	transferBackends []snettypes.PeerUintptr,
	pChunkMask *offheap.ChunkMask,
	request *snettypes.Request, response *snettypes.Response) {
	if uPeer == 0 {
		return
	}

	var err error
	request.OffheapBody.OffheapBytes = uUploadJob.Ptr().UMemBlock.Ptr().Bytes.Data
	for chunkMaskIndex := 0; chunkMaskIndex < pChunkMask.MaskArrayLen; chunkMaskIndex++ {
		request.OffheapBody.CopyOffset = pChunkMask.MaskArray[chunkMaskIndex].Offset
		request.OffheapBody.CopyEnd = pChunkMask.MaskArray[chunkMaskIndex].End
		err = p.snetClientDriver.Call(uPeer,
			"/NetBlock/PWrite", request, response)
		if err != nil {
			*uploadJobErr = err
			return
		}
	}

	// var protocolBuilder flatbuffers.Builder
	// protocolBuilder.Reset()

	// off = protocolBuilder.CreateByteVector([]byte("fuck"))
	// protocol.UploadJobBackendStart(&protocolBuilder)
	// protocol.UploadJobBackendAddPeerID(&protocolBuilder, off)
	// off = protocol.UploadJobBackendEnd(&protocolBuilder)

	// protocol.UploadJobStart(&protocolBuilder)
	// protocol.UploadJobStartTransferBackendsVector(&protocolBuilder, 1)
	// protocolBuilder.EndVector()
	// protocol.UploadJobAddOffset(&protocolBuilder)
	// protocol.UploadJobAddTransferBackends(&protocolBuilder)
}

func (p *netBlockDriverUploader) cronUpload() error {
	var (
		uUploadJob    UploadJobUintptr
		pUploadJob    *UploadJob
		uTmpChunkMask offheap.ChunkMaskUintptr
		request       [types.MaxDataNodesSizeStoreNetBlock]snettypes.Request
		response      [types.MaxDataNodesSizeStoreNetBlock]snettypes.Response
		pChunkMask    *offheap.ChunkMask
		pNetBlock     *types.NetBlock
		uploadJobSig  sync.WaitGroup
		dataNodeIndex int
		i             int
		ok            bool
		err           error
	)

	for {
		uUploadJob, ok = <-p.uploadJobChan
		if !ok {
			panic("uploadJobChan closed")
		}

		pUploadJob = uUploadJob.Ptr()
		pNetBlock = pUploadJob.UNetBlock.Ptr()

		p.uploadJobMutex.Lock()
		if pUploadJob.UploadMaskWaiting.Ptr().MaskArrayLen == 0 {
			pNetBlock.UploadSig.Done()
			p.uploadJobMutex.Unlock()
			continue
		}
		uTmpChunkMask = pUploadJob.UploadMaskProcessing
		pUploadJob.UploadMaskProcessing = pUploadJob.UploadMaskWaiting
		pUploadJob.UploadMaskWaiting = uTmpChunkMask
		p.uploadJobMutex.Unlock()

		pChunkMask = pUploadJob.UploadMaskProcessing.Ptr()

		// upload primary backend
		if pUploadJob.PrimaryBackendTransferCount > 0 {
			p.doUpload(&err, &uploadJobSig, uUploadJob,
				pUploadJob.Backends.Arr[0],
				pUploadJob.Backends.Arr[1:1+pUploadJob.PrimaryBackendTransferCount],
				pChunkMask,
				&request[dataNodeIndex], &response[dataNodeIndex])
		} else {
			p.doUpload(&err, &uploadJobSig, uUploadJob,
				pUploadJob.Backends.Arr[0],
				nil,
				pChunkMask,
				&request[dataNodeIndex], &response[dataNodeIndex])
		}

		// upload other backends
		for i = pUploadJob.PrimaryBackendTransferCount + 1; i < pUploadJob.Backends.Len; i++ {
			p.doUpload(&err, &uploadJobSig, uUploadJob,
				pUploadJob.Backends.Arr[i],
				nil,
				pChunkMask,
				&request[dataNodeIndex], &response[dataNodeIndex])
		}

		uploadJobSig.Wait()

		if err != nil {
			break
		}

		uUploadJob.Ptr().UNetBlock.Ptr().UploadSig.Done()

		// TODO catch error
		if err != nil {
			return err
		}

		pChunkMask.Reset()
	}

	return nil
}

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
			uUploadJob.Ptr().UNetBlock = uNetBlock
			uUploadJob.Ptr().UMemBlock = uMemBlock
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
