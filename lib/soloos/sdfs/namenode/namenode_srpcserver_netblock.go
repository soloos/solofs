package namenode

import (
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeSRPCServer) NetBlockPrepareMetadata(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	var (
		param           = make([]byte, reqBodySize)
		req             protocol.NetINodeNetBlockInfoRequest
		uNetINode       types.NetINodeUintptr
		netINodeID      types.NetINodeID
		uNetBlock       types.NetBlockUintptr
		protocolBuilder flatbuffers.Builder
		err             error
	)

	err = conn.ReadAll(param)
	if err != nil {
		return err
	}

	// request
	req.Init(param, flatbuffers.GetUOffsetT(param))
	copy(netINodeID[:], req.InodeID())
	uNetINode, err = p.nameNode.MetaStg.GetNetINode(netINodeID)
	if err != nil {
		if err == types.ErrObjectNotExists {
			api.SetNetINodeNetBlockInfoResponseError(snettypes.CODE_404, err.Error(), &protocolBuilder)
			conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		}
		goto SERVICE_DONE
	}

	// response
	uNetBlock, err = p.nameNode.MetaStg.MustGetNetBlock(uNetINode, int(req.NetBlockIndex()))
	if err != nil {
		api.SetNetINodeNetBlockInfoResponseError(snettypes.CODE_502, err.Error(), &protocolBuilder)
		conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	api.SetNetINodeNetBlockInfoResponse(uNetBlock.Ptr().DataNodes.Slice(),
		req.Cap(), req.Cap(), &protocolBuilder)
	err = conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])

SERVICE_DONE:
	return err
}
