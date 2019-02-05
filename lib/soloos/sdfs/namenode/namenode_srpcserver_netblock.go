package namenode

import (
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeSRPCServer) NetBlockPrepareMetaData(reqID uint64,
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
	copy(netINodeID[:], req.NetINodeID())
	uNetINode, err = p.nameNode.netINodeDriver.GetNetINodeWithReadAcquire(false, netINodeID)
	defer p.nameNode.netINodeDriver.ReleaseNetINodeWithReadRelease(uNetINode)

	if err != nil {
		if err == types.ErrObjectNotExists {
			api.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_404, err.Error())
		} else {
			api.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		}
		conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	// response
	uNetBlock, err = p.nameNode.netBlockDriver.MustGetNetBlock(uNetINode, int(req.NetBlockIndex()))
	if err != nil {
		api.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	api.SetNetINodeNetBlockInfoResponse(&protocolBuilder,
		uNetBlock.Ptr().StorDataBackends.Slice(), req.Cap(), req.Cap())
	err = conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])

SERVICE_DONE:
	return err
}
