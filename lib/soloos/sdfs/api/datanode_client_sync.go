package api

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeClient) GetNetINodeMetaData(uDataNode snettypes.PeerUintptr, uNetINode types.NetINodeUintptr) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		netINodeIDOff   flatbuffers.UOffsetT
		err             error
	)

	netINodeIDOff = protocolBuilder.CreateByteString(uNetINode.Ptr().ID[:])
	protocol.NetINodeSyncRequestStart(&protocolBuilder)
	protocol.NetINodeSyncRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	protocolBuilder.Finish(protocol.NetINodeSyncRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	err = p.SNetClientDriver.Call(uDataNode,
		"/NetINode/Sync", &req, &resp)
	if err != nil {
		return err
	}

	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	err = p.SNetClientDriver.ReadResponse(uDataNode, &req, &resp, body)
	if err != nil {
		return err
	}

	var (
		commonResponse protocol.CommonResponse
	)

	commonResponse.Init(body, flatbuffers.GetUOffsetT(body))
	err = CommonResponseToError(&commonResponse)
	if err != nil {
		return err
	}

	return nil
}
