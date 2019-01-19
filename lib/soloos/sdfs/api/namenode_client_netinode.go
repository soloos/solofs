package api

import (
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeClient) doGetNetINodeMetaData(isMustGet bool,
	uNetINode types.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int,
) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		netINodeIDOff   flatbuffers.UOffsetT
		err             error
	)

	netINodeIDOff = protocolBuilder.CreateByteString(uNetINode.Ptr().ID[:])
	protocol.NetINodeInfoRequestStart(&protocolBuilder)
	protocol.NetINodeInfoRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	protocol.NetINodeInfoRequestAddSize(&protocolBuilder, size)
	protocol.NetINodeInfoRequestAddNetBlockCap(&protocolBuilder, int32(netBlockCap))
	protocol.NetINodeInfoRequestAddMemBlockCap(&protocolBuilder, int32(memBlockCap))
	protocolBuilder.Finish(protocol.NetINodeNetBlockInfoRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	if isMustGet {
		err = p.snetClientDriver.Call(p.nameNodePeer,
			"/NetINode/MustGet", &req, &resp)
	} else {
		err = p.snetClientDriver.Call(p.nameNodePeer,
			"/NetINode/Get", &req, &resp)
	}

	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	p.snetClientDriver.ReadResponse(p.nameNodePeer, &req, &resp, body)
	if err != nil {
		return err
	}

	var (
		netINodeInfo   protocol.NetINodeInfoResponse
		commonResponse protocol.CommonResponse
	)

	netINodeInfo.Init(body, flatbuffers.GetUOffsetT(body))
	netINodeInfo.CommonResponse(&commonResponse)
	err = CommonResponseToError(&commonResponse)
	if err != nil {
		return err
	}

	uNetINode.Ptr().Size = netINodeInfo.Size()
	uNetINode.Ptr().NetBlockCap = int(netINodeInfo.NetBlockCap())
	uNetINode.Ptr().MemBlockCap = int(netINodeInfo.MemBlockCap())

	return nil
}

func (p *NameNodeClient) GetNetINodeMetaData(uNetINode types.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int,
) error {
	return p.doGetNetINodeMetaData(false, uNetINode, size, netBlockCap, memBlockCap)
}

func (p *NameNodeClient) MustGetNetINodeMetaData(uNetINode types.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int,
) error {
	return p.doGetNetINodeMetaData(true, uNetINode, size, netBlockCap, memBlockCap)
}
