package netstg

import (
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *netBlockDriverUploader) doUpload(uploadJobErr *error, uploadJobSig *sync.WaitGroup,
	uUploadJob UploadJobUintptr,
	uPeer snettypes.PeerUintptr,
	transferBackends []snettypes.PeerUintptr,
	pChunkMask *offheap.ChunkMask,
	request *snettypes.Request, response *snettypes.Response) {
	if uPeer == 0 {
		return
	}

	var (
		protocolBuilder     flatbuffers.Builder
		netBlockIDOff       flatbuffers.UOffsetT
		backendOff          flatbuffers.UOffsetT
		netBlockBytesOffset int
		netBlockBytesEnd    int
		memBlockSize        int
		i                   int
	)

	var err error
	request.OffheapBody.OffheapBytes = uUploadJob.Ptr().UMemBlock.Ptr().Bytes.Data
	memBlockSize = uUploadJob.Ptr().UMemBlock.Ptr().Bytes.Len
	for chunkMaskIndex := 0; chunkMaskIndex < pChunkMask.MaskArrayLen; chunkMaskIndex++ {
		request.OffheapBody.CopyOffset = pChunkMask.MaskArray[chunkMaskIndex].Offset
		request.OffheapBody.CopyEnd = pChunkMask.MaskArray[chunkMaskIndex].End
		netBlockBytesOffset = memBlockSize * uUploadJob.Ptr().MemBlockIndex
		netBlockBytesEnd = netBlockBytesOffset + request.OffheapBody.CopyEnd
		netBlockBytesOffset = netBlockBytesOffset + request.OffheapBody.CopyOffset

		if len(transferBackends) > 0 {
			var backendOffs = make([]flatbuffers.UOffsetT, len(transferBackends))
			var peerOff, addrOff flatbuffers.UOffsetT
			for i = 0; i < len(transferBackends); i++ {
				peerOff = protocolBuilder.CreateByteVector(transferBackends[i].Ptr().ID[:])
				addrOff = protocolBuilder.CreateString(transferBackends[i].Ptr().AddressStr())
				protocol.NetBlockPWriteBackendStart(&protocolBuilder)
				protocol.NetBlockPWriteBackendAddPeerID(&protocolBuilder, peerOff)
				protocol.NetBlockPWriteBackendAddAddress(&protocolBuilder, addrOff)
				backendOffs[i] = protocol.NetBlockPWriteBackendEnd(&protocolBuilder)
			}

			protocol.NetBlockPWriteStartTransferBackendsVector(&protocolBuilder, len(transferBackends))
			for i = len(transferBackends) - 1; i >= 0; i-- {
				protocolBuilder.PrependUOffsetT(backendOffs[i])
			}
			backendOff = protocolBuilder.EndVector(len(transferBackends))
		}

		netBlockIDOff = protocolBuilder.CreateByteVector(uUploadJob.Ptr().UNetBlock.Ptr().ID[:])
		protocol.NetBlockPWriteStart(&protocolBuilder)
		if len(transferBackends) > 0 {
			protocol.NetBlockPWriteAddTransferBackends(&protocolBuilder, backendOff)
		}
		protocol.NetBlockPWriteAddNetBlockID(&protocolBuilder, netBlockIDOff)
		protocol.NetBlockPWriteAddOffset(&protocolBuilder, int32(netBlockBytesOffset))
		protocol.NetBlockPWriteAddLength(&protocolBuilder, int32(netBlockBytesEnd))
		protocolBuilder.Finish(protocol.NetBlockPWriteEnd(&protocolBuilder))
		request.Parameter = protocolBuilder.Bytes[protocolBuilder.Head():]

		err = p.snetClientDriver.Call(uPeer,
			"/NetBlock/PWrite", request, response)
		if err != nil {
			*uploadJobErr = err
			return
		}

		protocolBuilder.Reset()
	}
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
