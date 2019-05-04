package api

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeClient) PReadMemBlock(uNetINode types.NetINodeUintptr,
	uPeer snettypes.PeerUintptr,
	uNetBlock types.NetBlockUintptr,
	netBlockIndex int32,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int32,
	offset uint64, length int,
) (int, error) {
	if uNetBlock.Ptr().LocalDataBackend != 0 {
		return p.preadMemBlockWithDisk(uNetINode, uPeer, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, offset, length)
	}

	switch uPeer.Ptr().ServiceProtocol {
	case snettypes.ProtocolSRPC:
		return p.doPReadMemBlockWithSRPC(uNetINode, uPeer, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, offset, length)
	}

	return 0, types.ErrServiceNotExists
}

func (p *DataNodeClient) doPReadMemBlockWithSRPC(uNetINode types.NetINodeUintptr,
	uPeer snettypes.PeerUintptr,
	uNetBlock types.NetBlockUintptr,
	netBlockIndex int32,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int32,
	offset uint64, length int,
) (int, error) {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		netINodeIDOff   flatbuffers.UOffsetT
		err             error
	)

	netINodeIDOff = protocolBuilder.CreateByteVector(uNetBlock.Ptr().NetINodeID[:])
	protocol.NetINodePReadRequestStart(&protocolBuilder)
	protocol.NetINodePReadRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	protocol.NetINodePReadRequestAddOffset(&protocolBuilder, offset)
	protocol.NetINodePReadRequestAddLength(&protocolBuilder, int32(length))
	protocolBuilder.Finish(protocol.NetINodePReadRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	// TODO choose datanode
	err = p.SNetClientDriver.Call(uPeer,
		"/NetINode/PRead", &req, &resp)
	if err != nil {
		return 0, err
	}

	var (
		netBlockPReadResp           protocol.NetINodePReadResponse
		commonResp                  protocol.CommonResponse
		param                       = make([]byte, resp.ParamSize)
		offsetInMemBlock, readedLen int
	)
	err = p.SNetClientDriver.ReadResponse(uPeer, &req, &resp, param)
	if err != nil {
		return 0, err
	}

	netBlockPReadResp.Init(param, flatbuffers.GetUOffsetT(param))
	netBlockPReadResp.CommonResponse(&commonResp)
	if commonResp.Code() != snettypes.CODE_OK {
		return 0, types.ErrNetBlockPRead
	}

	offsetInMemBlock = int(offset - uint64(uMemBlock.Ptr().Bytes.Cap)*uint64(memBlockIndex))
	readedLen = int(resp.BodySize - resp.ParamSize)
	err = p.SNetClientDriver.ReadResponse(uPeer, &req, &resp,
		(*uMemBlock.Ptr().BytesSlice())[offsetInMemBlock:readedLen])
	if err != nil {
		return 0, err
	}

	return int(netBlockPReadResp.Length()), err
}
