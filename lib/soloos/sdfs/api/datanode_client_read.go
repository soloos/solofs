package api

import (
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeClient) PRead(uPeer snettypes.PeerUintptr,
	uNetBlock types.NetBlockUintptr,
	offset int, length int,
	resp *snettypes.Response,
) error {
	var (
		req             snettypes.Request
		protocolBuilder flatbuffers.Builder
		err             error
	)

	peerOff := protocolBuilder.CreateByteVector(uNetBlock.Ptr().ID[:])
	protocol.NetBlockPReadRequestStart(&protocolBuilder)
	protocol.NetBlockPReadRequestAddNetBlockID(&protocolBuilder, peerOff)
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

	return nil
}
