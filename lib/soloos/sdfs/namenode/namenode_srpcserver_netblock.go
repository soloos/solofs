package namenode

import (
	snettypes "soloos/common/snet/types"
	"soloos/common/sdfsapi"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeSRPCServer) NetBlockPrepareMetaData(serviceReq *snettypes.NetQuery) error {
	var (
		param           = make([]byte, serviceReq.BodySize)
		req             protocol.NetINodeNetBlockInfoRequest
		uNetINode       types.NetINodeUintptr
		netINodeID      types.NetINodeID
		uNetBlock       types.NetBlockUintptr
		protocolBuilder flatbuffers.Builder
		err             error
	)

	err = serviceReq.ReadAll(param)
	if err != nil {
		return err
	}

	// request
	req.Init(param, flatbuffers.GetUOffsetT(param))
	copy(netINodeID[:], req.NetINodeID())
	uNetINode, err = p.nameNode.netINodeDriver.GetNetINode(netINodeID)
	defer p.nameNode.netINodeDriver.ReleaseNetINode(uNetINode)

	if err != nil {
		if err == types.ErrObjectNotExists {
			sdfsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_404, err.Error())
		} else {
			sdfsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		}
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	// response
	uNetBlock, err = p.nameNode.netBlockDriver.MustGetNetBlock(uNetINode, req.NetBlockIndex())
	defer p.nameNode.netBlockDriver.ReleaseNetBlock(uNetBlock)
	if err != nil {
		sdfsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	sdfsapi.SetNetINodeNetBlockInfoResponse(&protocolBuilder,
		uNetBlock.Ptr().StorDataBackends.Slice(), req.Cap(), req.Cap())
	err = serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])

SERVICE_DONE:
	return err
}
