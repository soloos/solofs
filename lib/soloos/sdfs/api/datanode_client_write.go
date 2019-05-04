package api

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdbone/offheap"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeClient) UploadMemBlock(uJob types.UploadMemBlockJobUintptr,
	uploadPeerIndex int, transferPeersCount int,
) error {
	var (
		uDataNode snettypes.PeerUintptr
	)
	uDataNode = uJob.Ptr().UNetBlock.Ptr().SyncDataBackends.Arr[uploadPeerIndex]
	switch uDataNode.Ptr().ServiceProtocol {
	case snettypes.ProtocolDisk:
		return p.uploadMemBlockWithDisk(uJob, uploadPeerIndex, transferPeersCount)
	case snettypes.ProtocolSRPC:
		return p.doUploadMemBlockWithSRPC(uJob, uploadPeerIndex, transferPeersCount)
	}

	return nil
}

func (p *DataNodeClient) doUploadMemBlockWithSRPC(uJob types.UploadMemBlockJobUintptr,
	uploadPeerIndex int, transferPeersCount int,
) error {
	var (
		req                 snettypes.Request
		resp                snettypes.Response
		protocolBuilder     flatbuffers.Builder
		netINodeIDOff       flatbuffers.UOffsetT
		backendOff          flatbuffers.UOffsetT
		uNetBlock           types.NetBlockUintptr
		netINodeWriteOffset int
		netINodeWriteLength int
		memBlockCap         int
		peerOff, addrOff    flatbuffers.UOffsetT
		backendOffs         = make([]flatbuffers.UOffsetT, 8)
		pChunkMask          *offheap.ChunkMask
		commonResp          protocol.CommonResponse
		respBody            []byte
		i                   int
		uPeer               snettypes.PeerUintptr
		err                 error
	)

	uNetBlock = uJob.Ptr().UNetBlock
	pChunkMask = uJob.Ptr().UploadMaskProcessing.Ptr()

	req.OffheapBody.OffheapBytes = uJob.Ptr().UMemBlock.Ptr().Bytes.Data
	memBlockCap = uJob.Ptr().UMemBlock.Ptr().Bytes.Len
	for chunkMaskIndex := 0; chunkMaskIndex < pChunkMask.MaskArrayLen; chunkMaskIndex++ {
		req.OffheapBody.CopyOffset = pChunkMask.MaskArray[chunkMaskIndex].Offset
		req.OffheapBody.CopyEnd = pChunkMask.MaskArray[chunkMaskIndex].End
		netINodeWriteOffset = memBlockCap*int(uJob.Ptr().MemBlockIndex) + req.OffheapBody.CopyOffset
		netINodeWriteLength = req.OffheapBody.CopyEnd - req.OffheapBody.CopyOffset

		if transferPeersCount > 0 {
			for i = 0; i < transferPeersCount; i++ {
				uPeer = uNetBlock.Ptr().SyncDataBackends.Arr[uploadPeerIndex+1+i]
				peerOff = protocolBuilder.CreateByteVector(uPeer.Ptr().ID[:])
				addrOff = protocolBuilder.CreateString(uPeer.Ptr().AddressStr())
				protocol.SNetPeerStart(&protocolBuilder)
				protocol.SNetPeerAddPeerID(&protocolBuilder, peerOff)
				protocol.SNetPeerAddAddress(&protocolBuilder, addrOff)
				if i < cap(backendOffs) {
					backendOffs[i] = protocol.SNetPeerEnd(&protocolBuilder)
				} else {
					backendOffs = append(backendOffs, protocol.SNetPeerEnd(&protocolBuilder))
				}
			}

			protocol.NetINodePWriteRequestStartTransferBackendsVector(&protocolBuilder, transferPeersCount)
			for i = transferPeersCount - 1; i >= 0; i-- {
				protocolBuilder.PrependUOffsetT(backendOffs[i])
			}
			backendOff = protocolBuilder.EndVector(transferPeersCount)
		}

		netINodeIDOff = protocolBuilder.CreateByteVector(uNetBlock.Ptr().NetINodeID[:])
		protocol.NetINodePWriteRequestStart(&protocolBuilder)
		if transferPeersCount > 0 {
			protocol.NetINodePWriteRequestAddTransferBackends(&protocolBuilder, backendOff)
		}
		protocol.NetINodePWriteRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
		protocol.NetINodePWriteRequestAddOffset(&protocolBuilder, uint64(netINodeWriteOffset))
		protocol.NetINodePWriteRequestAddLength(&protocolBuilder, int32(netINodeWriteLength))
		protocolBuilder.Finish(protocol.NetINodePWriteRequestEnd(&protocolBuilder))
		req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

		uPeer = uJob.Ptr().UNetBlock.Ptr().SyncDataBackends.Arr[uploadPeerIndex]
		err = p.SNetClientDriver.Call(uPeer,
			"/NetINode/PWrite", &req, &resp)
		if err != nil {
			goto PWRITE_DONE
		}

		respBody = make([]byte, resp.ParamSize)
		err = p.SNetClientDriver.ReadResponse(uPeer, &req, &resp, respBody)
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
