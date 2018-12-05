package netstg

import (
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *netBlockDriverUploader) doUpload(uploadMemBlockJobErr *error, uploadMemBlockJobSig *sync.WaitGroup,
	uUploadMemBlockJob types.UploadMemBlockJobUintptr,
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
		peerOff, addrOff    flatbuffers.UOffsetT
		backendOffs         = make([]flatbuffers.UOffsetT, 8)
		commonResp          protocol.CommonResponse
		respBody            = make([]byte, 64)
		i                   int
		err                 error
	)

	request.OffheapBody.OffheapBytes = uUploadMemBlockJob.Ptr().UMemBlock.Ptr().Bytes.Data
	memBlockSize = uUploadMemBlockJob.Ptr().UMemBlock.Ptr().Bytes.Len
	for chunkMaskIndex := 0; chunkMaskIndex < pChunkMask.MaskArrayLen; chunkMaskIndex++ {
		request.OffheapBody.CopyOffset = pChunkMask.MaskArray[chunkMaskIndex].Offset
		request.OffheapBody.CopyEnd = pChunkMask.MaskArray[chunkMaskIndex].End
		netBlockBytesOffset = memBlockSize * uUploadMemBlockJob.Ptr().MemBlockIndex
		netBlockBytesEnd = netBlockBytesOffset + request.OffheapBody.CopyEnd
		netBlockBytesOffset = netBlockBytesOffset + request.OffheapBody.CopyOffset

		if len(transferBackends) > 0 {
			for i = 0; i < len(transferBackends); i++ {
				peerOff = protocolBuilder.CreateByteVector(transferBackends[i].Ptr().ID[:])
				addrOff = protocolBuilder.CreateString(transferBackends[i].Ptr().AddressStr())
				protocol.NetBlockBackendStart(&protocolBuilder)
				protocol.NetBlockBackendAddPeerID(&protocolBuilder, peerOff)
				protocol.NetBlockBackendAddAddress(&protocolBuilder, addrOff)
				if i < cap(backendOffs) {
					backendOffs[i] = protocol.NetBlockBackendEnd(&protocolBuilder)
				} else {
					backendOffs = append(backendOffs, protocol.NetBlockBackendEnd(&protocolBuilder))
				}
			}

			protocol.NetBlockPWriteRequestStartTransferBackendsVector(&protocolBuilder, len(transferBackends))
			for i = len(transferBackends) - 1; i >= 0; i-- {
				protocolBuilder.PrependUOffsetT(backendOffs[i])
			}
			backendOff = protocolBuilder.EndVector(len(transferBackends))
		}

		netBlockIDOff = protocolBuilder.CreateByteVector(uUploadMemBlockJob.Ptr().UNetBlock.Ptr().ID[:])
		protocol.NetBlockPWriteRequestStart(&protocolBuilder)
		if len(transferBackends) > 0 {
			protocol.NetBlockPWriteRequestAddTransferBackends(&protocolBuilder, backendOff)
		}
		protocol.NetBlockPWriteRequestAddNetBlockID(&protocolBuilder, netBlockIDOff)
		protocol.NetBlockPWriteRequestAddOffset(&protocolBuilder, int32(netBlockBytesOffset))
		protocol.NetBlockPWriteRequestAddLength(&protocolBuilder, int32(netBlockBytesEnd))
		protocolBuilder.Finish(protocol.NetBlockPWriteRequestEnd(&protocolBuilder))
		request.Parameter = protocolBuilder.Bytes[protocolBuilder.Head():]

		err = p.driver.snetClientDriver.Call(uPeer,
			"/NetBlock/PWrite", request, response)
		if err != nil {
			*uploadMemBlockJobErr = err
			return
		}

		if response.ParameterSize > uint32(cap(respBody)) {
			respBody = append(respBody, util.DevNullBuf[:int(response.ParameterSize-uint32(cap(respBody)))]...)
		}
		err = p.driver.snetClientDriver.ReadResponse(uPeer, request, response, respBody)
		if err != nil {
			*uploadMemBlockJobErr = err
			return
		}
		commonResp.Init(respBody[:(response.ParameterSize)], flatbuffers.GetUOffsetT(respBody[:(response.ParameterSize)]))
		if commonResp.Code() != snettypes.CODE_OK {
			err = types.ErrNetBlockPWrite
			return
		}

		protocolBuilder.Reset()
	}
}

func (p *netBlockDriverUploader) cronUpload() error {
	var (
		uUploadMemBlockJob   types.UploadMemBlockJobUintptr
		pUploadMemBlockJob   *types.UploadMemBlockJob
		request              [types.MaxDataNodesSizeStoreNetBlock]snettypes.Request
		response             [types.MaxDataNodesSizeStoreNetBlock]snettypes.Response
		pChunkMask           *offheap.ChunkMask
		uploadMemBlockJobSig sync.WaitGroup
		dataNodeIndex        int
		i                    int
		ok                   bool
		err                  error
	)

	for {
		uUploadMemBlockJob, ok = <-p.uploadMemBlockJobChan
		if !ok {
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
		if pUploadMemBlockJob.PrimaryBackendTransferCount > 0 {
			p.doUpload(&err, &uploadMemBlockJobSig, uUploadMemBlockJob,
				pUploadMemBlockJob.Backends.Arr[0],
				pUploadMemBlockJob.Backends.Arr[1:1+pUploadMemBlockJob.PrimaryBackendTransferCount],
				pChunkMask,
				&request[dataNodeIndex], &response[dataNodeIndex])
		} else {
			p.doUpload(&err, &uploadMemBlockJobSig, uUploadMemBlockJob,
				pUploadMemBlockJob.Backends.Arr[0],
				nil,
				pChunkMask,
				&request[dataNodeIndex], &response[dataNodeIndex])
		}

		// upload other backends
		for i = pUploadMemBlockJob.PrimaryBackendTransferCount + 1; i < pUploadMemBlockJob.Backends.Len; i++ {
			p.doUpload(&err, &uploadMemBlockJobSig, uUploadMemBlockJob,
				pUploadMemBlockJob.Backends.Arr[i],
				nil,
				pChunkMask,
				&request[dataNodeIndex], &response[dataNodeIndex])
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
