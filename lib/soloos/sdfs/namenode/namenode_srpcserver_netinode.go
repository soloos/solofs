package namenode

import (
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeSRPCServer) NetINodeMustGet(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	var (
		param           = make([]byte, reqBodySize)
		req             protocol.NetINodeInfoRequest
		uNetINode       types.NetINodeUintptr
		netINodeID      types.NetINodeID
		protocolBuilder flatbuffers.Builder
		err             error
	)

	err = conn.ReadAll(param)
	if err != nil {
		return err
	}

	// request
	req.Init(param, flatbuffers.GetUOffsetT(param))

	copy(netINodeID[:], req.NetINodeID())
	uNetINode, err = p.nameNode.metaStg.MustGetNetINode(netINodeID,
		req.Size(), int(req.NetBlockCap()), int(req.MemBlockCap()))
	if err != nil {
		api.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	// response
	api.SetNetINodeInfoResponse(&protocolBuilder,
		uNetINode.Ptr().Size, int32(uNetINode.Ptr().NetBlockCap), int32(uNetINode.Ptr().MemBlockCap))
	conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
	err = nil

SERVICE_DONE:
	return err
}
