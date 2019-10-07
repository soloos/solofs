package memstg

import (
	"soloos/common/snet"
	"soloos/common/solofsprotocol"
	"soloos/common/solofstypes"
	"soloos/common/util"
	"soloos/solodb/offheap"
)

func (p *NetBlockDriver) doUploadMemBlockWithSolofs(uJob solofstypes.UploadMemBlockJobUintptr,
	uploadPeerIndex int,
) error {
	var (
		snetReq            snet.SNetReq
		snetResp           snet.SNetResp
		req                solofsprotocol.NetINodePWriteReq
		transferPeersCount int
		memBlockCap        int
		uploadChunkMask    offheap.ChunkMask
		respParamBs        []byte
		backendPeer        snet.Peer
		i                  int
		err                error
	)

	var pJob = uJob.Ptr()
	var pNetBlock = pJob.UNetBlock.Ptr()
	var pMemBlock = pJob.UMemBlock.Ptr()
	uploadChunkMask = pJob.GetProcessingChunkMask()
	transferPeersCount = int(pNetBlock.SyncDataBackends.Arr[uploadPeerIndex].TransferCount)

	snetReq.OffheapBody.OffheapBytes = pMemBlock.Bytes.Data
	memBlockCap = pMemBlock.Bytes.Len
	for chunkMaskIndex := 0; chunkMaskIndex < uploadChunkMask.MaskArrayLen; chunkMaskIndex++ {
		snetReq.OffheapBody.CopyOffset = uploadChunkMask.MaskArray[chunkMaskIndex].Offset
		snetReq.OffheapBody.CopyEnd = uploadChunkMask.MaskArray[chunkMaskIndex].End

		req.NetINodeID = pNetBlock.NetINodeID
		req.Offset = uint64(memBlockCap)*uint64(pJob.MemBlockIndex) + uint64(snetReq.OffheapBody.CopyOffset)
		req.Length = int32(snetReq.OffheapBody.CopyEnd - snetReq.OffheapBody.CopyOffset)
		req.TransferBackends = req.TransferBackends[:0]
		for i = 0; i < transferPeersCount; i++ {
			backendPeer, _ = p.SNetDriver.GetPeer(pNetBlock.SyncDataBackends.Arr[uploadPeerIndex+1+i].PeerID)
			req.TransferBackends = append(req.TransferBackends, backendPeer.PeerIDStr())
		}

		backendPeer, err = p.SNetDriver.GetPeer(pNetBlock.SyncDataBackends.Arr[uploadPeerIndex].PeerID)
		if err != nil {
			goto PWRITE_DONE
		}

		snetReq.Param = snet.MustSpecMarshalRequest(req)
		err = p.solodnClient.Call(backendPeer.ID,
			"/NetINode/PWrite", &snetReq, &snetResp)
		if err != nil {
			goto PWRITE_DONE
		}

		util.ChangeBytesArraySize(&respParamBs, int(snetResp.ParamSize))
		err = p.solodnClient.ReadResponse(backendPeer.ID, &snetReq, &snetResp, respParamBs, nil)
		if err != nil {
			goto PWRITE_DONE
		}
	}

PWRITE_DONE:
	return err
}

func (p *NetBlockDriver) UploadMemBlockToNet(uJob solofstypes.UploadMemBlockJobUintptr,
	uploadPeerIndex int,
) error {
	var solodn, _ = p.SoloosEnv.SNetDriver.GetPeer(
		uJob.Ptr().UNetBlock.Ptr().SyncDataBackends.Arr[uploadPeerIndex].PeerID)
	switch solodn.ServiceProtocol {
	case snet.ProtocolLocalFs:
		return p.helper.UploadMemBlockWithDisk(uJob, uploadPeerIndex)
	case snet.ProtocolSolomq:
		return p.helper.UploadMemBlockWithSolomq(uJob, uploadPeerIndex)
	case snet.ProtocolSolofs:
		return p.doUploadMemBlockWithSolofs(uJob, uploadPeerIndex)
	}

	return nil
}
