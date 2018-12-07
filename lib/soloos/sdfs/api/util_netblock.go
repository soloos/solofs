package api

import (
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
	"soloos/util"

	flatbuffers "github.com/google/flatbuffers/go"
)

func SetINodeNetBlockInfoResp(backends []snettypes.PeerUintptr, netBlockLen, netBlockCap int32,
	protocolBuilder *flatbuffers.Builder) {
	var (
		netBlockID    types.NetBlockID
		peerOff       flatbuffers.UOffsetT
		addrOff       flatbuffers.UOffsetT
		backendOff    flatbuffers.UOffsetT
		netBlockIDOff flatbuffers.UOffsetT
		i             int
	)

	backendOffs := make([]flatbuffers.UOffsetT, len(backends))

	for i = 0; i < len(backends); i++ {
		peerOff = protocolBuilder.CreateByteVector(backends[i].Ptr().PeerID[:])
		addrOff = protocolBuilder.CreateString(backends[i].Ptr().AddressStr())
		protocol.NetBlockBackendStart(protocolBuilder)
		protocol.NetBlockBackendAddPeerID(protocolBuilder, peerOff)
		protocol.NetBlockBackendAddAddress(protocolBuilder, addrOff)
		backendOffs[i] = protocol.NetBlockBackendEnd(protocolBuilder)
	}

	protocol.INodeNetBlockInfoResponseStartBackendsVector(protocolBuilder, len(backends))
	for i = len(backends) - 1; i >= 0; i-- {
		protocolBuilder.PrependUOffsetT(backendOffs[i])
	}
	backendOff = protocolBuilder.EndVector(len(backends))

	util.InitUUID64(&netBlockID)

	netBlockIDOff = protocolBuilder.CreateByteVector(netBlockID[:])
	protocol.INodeNetBlockInfoResponseStart(protocolBuilder)
	protocol.INodeNetBlockInfoResponseAddNetBlockID(protocolBuilder, netBlockIDOff)
	protocol.INodeNetBlockInfoResponseAddBackends(protocolBuilder, backendOff)
	protocol.INodeNetBlockInfoResponseAddLen(protocolBuilder, netBlockLen)
	protocol.INodeNetBlockInfoResponseAddCap(protocolBuilder, netBlockCap)
	protocolBuilder.Finish(protocol.INodeNetBlockInfoResponseEnd(protocolBuilder))
}

func SetNetBlockPReadResponse(code int32, length int32,
	protocolBuilder *flatbuffers.Builder) {
	var (
		commonResponseOff flatbuffers.UOffsetT
	)
	protocol.CommonResponseStart(protocolBuilder)
	protocol.CommonResponseAddCode(protocolBuilder, code)
	commonResponseOff = protocol.CommonResponseEnd(protocolBuilder)

	protocol.NetBlockPReadResponseStart(protocolBuilder)
	protocol.NetBlockPReadResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	protocolBuilder.Finish(protocol.NetBlockPReadResponseEnd(protocolBuilder))
}
