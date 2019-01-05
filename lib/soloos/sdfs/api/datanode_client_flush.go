package api

import (
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

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
	protocol.NetINodeFlushRequestStart(&protocolBuilder)
	protocol.NetINodeFlushRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	protocolBuilder.Finish(protocol.NetINodeFlushRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	err = p.snetClientDriver.Call(uDataNode,
		"/NetINode/Flush", &req, &resp)

	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	p.snetClientDriver.ReadResponse(uDataNode, &req, &resp, body)
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
