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
		req             protocol.INodeNetBlockInfoRequest
		uINode          types.INodeUintptr
		inodeID         types.INodeID
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
	copy(inodeID[:], req.InodeID())
	uINode, err = p.nameNode.MetaStg.GetINode(inodeID)
	if err != nil {
		if err == types.ErrObjectNotExists {
			api.SetINodeNetBlockInfoResponseError(snettypes.CODE_404, "", &protocolBuilder)
			conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		}
		goto SERVICE_DONE
	}

	// response
	uNetBlock, err = p.nameNode.MetaStg.MustGetNetBlock(uINode, int(req.NetBlockIndex()))
	if err != nil {
		api.SetINodeNetBlockInfoResponseError(snettypes.CODE_502, "", &protocolBuilder)
		conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	api.SetINodeNetBlockInfoResponse(uNetBlock.Ptr().DataNodes.Slice(),
		req.Cap(), req.Cap(), &protocolBuilder)
	err = conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])

SERVICE_DONE:
	return err
}
