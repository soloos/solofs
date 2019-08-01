package datanode

import (
	"soloos/common/sdbapitypes"
	"soloos/common/sdfsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/sdfsprotocol"
	"soloos/common/snettypes"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *SRPCServer) NetINodePWrite(serviceReq *snettypes.NetQuery) error {
	var (
		reqParamData     = make([]byte, serviceReq.ParamSize)
		reqParam         sdfsprotocol.NetINodePWriteRequest
		syncDataBackends snettypes.PeerGroup
		peerID           snettypes.PeerID
		uNetBlock        sdfsapitypes.NetBlockUintptr
		i                int
		err              error
	)

	// request param
	err = serviceReq.ReadAll(reqParamData)
	if err != nil {
		return err
	}
	reqParam.Init(reqParamData[:serviceReq.ParamSize], flatbuffers.GetUOffsetT(reqParamData[:serviceReq.ParamSize]))

	// response

	// get uNetINode
	var (
		protocolBuilder    flatbuffers.Builder
		netINodeID         sdfsapitypes.NetINodeID
		uNetINode          sdfsapitypes.NetINodeUintptr
		firstNetBlockIndex int32
		lastNetBlockIndex  int32
		netBlockIndex      int32
	)
	copy(netINodeID[:], reqParam.NetINodeID())
	uNetINode, err = p.dataNode.netINodeDriver.GetNetINode(netINodeID)
	defer p.dataNode.netINodeDriver.ReleaseNetINode(uNetINode)
	if err != nil {
		if err == sdfsapitypes.ErrObjectNotExists {
			sdfsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_404)
			goto SERVICE_DONE
		} else {
			sdfsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
			goto SERVICE_DONE
		}
	}

	// TODO no need prepare syncDataBackends every pwrite
	syncDataBackends.Reset()
	syncDataBackends.Append(p.dataNode.localFsSNetPeer.ID)
	for i = 0; i < reqParam.TransferBackendsLength(); i++ {
		copy(peerID[:], reqParam.TransferBackends(i))
		syncDataBackends.Append(peerID)
	}

	// prepare uNetBlock
	firstNetBlockIndex = int32(reqParam.Offset() / uint64(uNetINode.Ptr().NetBlockCap))
	lastNetBlockIndex = int32((reqParam.Offset() + uint64(reqParam.Length())) / uint64(uNetINode.Ptr().NetBlockCap))
	for netBlockIndex = firstNetBlockIndex; netBlockIndex <= lastNetBlockIndex; netBlockIndex++ {
		uNetBlock, err = p.dataNode.netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
		defer p.dataNode.netBlockDriver.ReleaseNetBlock(uNetBlock)
		if err != nil {
			sdfsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
			goto SERVICE_DONE
		}

		if uNetBlock.Ptr().IsSyncDataBackendsInited.Load() == sdbapitypes.MetaDataStateUninited {
			p.dataNode.PrepareNetBlockSyncDataBackends(uNetBlock, syncDataBackends)
		}
	}

	// request file data
	err = p.dataNode.netINodeDriver.PWriteWithNetQuery(uNetINode, serviceReq,
		int(reqParam.Length()), reqParam.Offset())
	if err != nil {
		return err
	}

SERVICE_DONE:
	if err != nil {
		return nil
	}

	if err == nil {
		sdfsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	}

	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	err = serviceReq.SimpleResponse(serviceReq.ReqID, respBody)
	if err != nil {
		return err
	}

	return nil
}
