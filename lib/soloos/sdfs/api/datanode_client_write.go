package api

import (
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeClient) UploadMemBlock(uUploadMemBlockJob types.UploadMemBlockJobUintptr,
	uPeer snettypes.PeerUintptr,
	transferBackends []snettypes.PeerUintptr,
) error {
	if uPeer == 0 {
		return nil
	}

	var (
		req                 snettypes.Request
		resp                snettypes.Response
		protocolBuilder     flatbuffers.Builder
		netINodeIDOff       flatbuffers.UOffsetT
		backendOff          flatbuffers.UOffsetT
		uNetBlock           types.NetBlockUintptr
		netBlockBytesOffset int
		netBlockBytesEnd    int
		memBlockCap         int
		peerOff, addrOff    flatbuffers.UOffsetT
		backendOffs         = make([]flatbuffers.UOffsetT, 8)
		pChunkMask          *offheap.ChunkMask
		commonResp          protocol.CommonResponse
		respBody            = make([]byte, 64)
		i                   int
		err                 error
	)

	uNetBlock = uUploadMemBlockJob.Ptr().UNetBlock
	pChunkMask = uUploadMemBlockJob.Ptr().UploadMaskProcessing.Ptr()

	req.OffheapBody.OffheapBytes = uUploadMemBlockJob.Ptr().UMemBlock.Ptr().Bytes.Data
	memBlockCap = uUploadMemBlockJob.Ptr().UMemBlock.Ptr().Bytes.Len
	for chunkMaskIndex := 0; chunkMaskIndex < pChunkMask.MaskArrayLen; chunkMaskIndex++ {
		req.OffheapBody.CopyOffset = pChunkMask.MaskArray[chunkMaskIndex].Offset
		req.OffheapBody.CopyEnd = pChunkMask.MaskArray[chunkMaskIndex].End
		netBlockBytesOffset = memBlockCap * uUploadMemBlockJob.Ptr().MemBlockIndex
		netBlockBytesEnd = netBlockBytesOffset + req.OffheapBody.CopyEnd
		netBlockBytesOffset = netBlockBytesOffset + req.OffheapBody.CopyOffset

		if len(transferBackends) > 0 {
			for i = 0; i < len(transferBackends); i++ {
				peerOff = protocolBuilder.CreateByteVector(transferBackends[i].Ptr().PeerID[:])
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

		netINodeIDOff = protocolBuilder.CreateByteVector(uNetBlock.Ptr().NetINodeID[:])
		protocol.NetBlockPWriteRequestStart(&protocolBuilder)
		if len(transferBackends) > 0 {
			protocol.NetBlockPWriteRequestAddTransferBackends(&protocolBuilder, backendOff)
		}
		protocol.NetBlockPWriteRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
		protocol.NetBlockPWriteRequestAddNetBlockIndex(&protocolBuilder, int32(uNetBlock.Ptr().IndexInNetINode))
		protocol.NetBlockPWriteRequestAddOffset(&protocolBuilder, int32(netBlockBytesOffset))
		protocol.NetBlockPWriteRequestAddLength(&protocolBuilder, int32(netBlockBytesEnd))
		protocolBuilder.Finish(protocol.NetBlockPWriteRequestEnd(&protocolBuilder))
		req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

		err = p.snetClientDriver.Call(uPeer,
			"/NetBlock/PWrite", &req, &resp)
		if err != nil {
			goto PWRITE_DONE
		}

		if resp.ParamSize > uint32(cap(respBody)) {
			respBody = append(respBody, util.DevNullBuf[:int(resp.ParamSize-uint32(cap(respBody)))]...)
		}
		err = p.snetClientDriver.ReadResponse(uPeer, &req, &resp, respBody)
		if err != nil {
			goto PWRITE_DONE
		}
		commonResp.Init(respBody[:(resp.ParamSize)], flatbuffers.GetUOffsetT(respBody[:(resp.ParamSize)]))
		if commonResp.Code() != snettypes.CODE_OK {
			err = types.ErrNetBlockPWrite
			goto PWRITE_DONE
		}
		protocolBuilder.Reset()
	}

PWRITE_DONE:
	return err
}
