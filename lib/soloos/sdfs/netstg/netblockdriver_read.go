package netstg

import (
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NetBlockDriver) PRead(uINode types.INodeUintptr,
	uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int,
	offset, length int) error {
	if uNetBlock.Ptr().DataNodes.Len == 0 {
		// TODO should return types.ErrBackendListIsEmpty
		return nil
		// return types.ErrBackendListIsEmpty
	}

	var (
		request         snettypes.Request
		response        snettypes.Response
		protocolBuilder flatbuffers.Builder
		err             error
	)

	peerOff := protocolBuilder.CreateByteVector(uNetBlock.Ptr().ID[:])
	protocol.NetBlockPReadRequestStart(&protocolBuilder)
	protocol.NetBlockPReadRequestAddNetBlockID(&protocolBuilder, peerOff)
	protocol.NetBlockPReadRequestAddOffset(&protocolBuilder, int32(offset))
	protocol.NetBlockPReadRequestAddLength(&protocolBuilder, int32(length))
	protocolBuilder.Finish(protocol.NetBlockPReadRequestEnd(&protocolBuilder))
	request.Parameter = protocolBuilder.Bytes[protocolBuilder.Head():]

	// TODO choose datanode
	err = p.snetClientDriver.Call(uNetBlock.Ptr().DataNodes.Arr[0],
		"/NetBlock/PRead", &request, &response)
	if err != nil {
		return err
	}

	return nil
}
