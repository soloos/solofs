package api

import (
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeClient) PRead(uPeer snettypes.PeerUintptr,
	uNetBlock types.NetBlockUintptr,
	netBlockIndex int,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int,
	offset int, length int,
	resp *snettypes.Response,
) error {
	var (
		req             snettypes.Request
		protocolBuilder flatbuffers.Builder
		netINodeIDOff   flatbuffers.UOffsetT
		err             error
	)

	netINodeIDOff = protocolBuilder.CreateByteVector(uNetBlock.Ptr().NetINodeID[:])
	protocol.NetBlockPReadRequestStart(&protocolBuilder)
	protocol.NetBlockPReadRequestAddNetINodeID(&protocolBuilder, netINodeIDOff)
	protocol.NetBlockPReadRequestAddNetBlockIndex(&protocolBuilder, int32(netBlockIndex))
	protocol.NetBlockPReadRequestAddOffset(&protocolBuilder, int32(offset))
	protocol.NetBlockPReadRequestAddLength(&protocolBuilder, int32(length))
	protocolBuilder.Finish(protocol.NetBlockPReadRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	// TODO choose datanode
	err = p.snetClientDriver.Call(uPeer,
		"/NetBlock/PRead", &req, resp)
	if err != nil {
		return err
	}

	var (
		netBlockPReadResp           protocol.NetBlockPReadResponse
		commonResp                  protocol.CommonResponse
		param                       = make([]byte, resp.ParamSize)
		offsetInMemBlock, readedLen int
	)
	err = p.snetClientDriver.ReadResponse(uPeer, &req, resp, param)
	if err != nil {
		return err
	}

	netBlockPReadResp.Init(param, flatbuffers.GetUOffsetT(param))
	netBlockPReadResp.CommonResponse(&commonResp)
	if commonResp.Code() != snettypes.CODE_OK {
		return types.ErrNetBlockPRead
	}

	offsetInMemBlock = int(offset - (uMemBlock.Ptr().Bytes.Cap * memBlockIndex))
	readedLen = int(resp.BodySize - resp.ParamSize)
	err = p.snetClientDriver.ReadResponse(uPeer, &req, resp,
		(*uMemBlock.Ptr().BytesSlice())[offsetInMemBlock:readedLen])

	return nil
}
