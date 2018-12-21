package api

import (
	"soloos/sdfs/protocol"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func SetNetINodeNetBlockInfoResponseError(protocolBuilder *flatbuffers.Builder, code int, err string) {
	protocolBuilder.Reset()
	var (
		errOff            flatbuffers.UOffsetT
		commonResponseOff flatbuffers.UOffsetT
	)
	errOff = protocolBuilder.CreateString(err)
	protocol.CommonResponseStart(protocolBuilder)
	protocol.CommonResponseAddCode(protocolBuilder, int32(code))
	protocol.CommonResponseAddError(protocolBuilder, errOff)
	commonResponseOff = protocol.CommonResponseEnd(protocolBuilder)

	protocol.NetINodeNetBlockInfoResponseStart(protocolBuilder)
	protocol.NetINodeNetBlockInfoResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	protocolBuilder.Finish(protocol.NetINodeNetBlockInfoResponseEnd(protocolBuilder))
}

func SetNetINodeNetBlockInfoResponse(protocolBuilder *flatbuffers.Builder,
	backends []snettypes.PeerUintptr, netBlockLen, netBlockCap int32) {
	var (
		peerOff           flatbuffers.UOffsetT
		addrOff           flatbuffers.UOffsetT
		backendOff        flatbuffers.UOffsetT
		commonResponseOff flatbuffers.UOffsetT
		i                 int
	)

	backendOffs := make([]flatbuffers.UOffsetT, len(backends))

	protocol.CommonResponseStart(protocolBuilder)
	protocol.CommonResponseAddCode(protocolBuilder, snettypes.CODE_OK)
	commonResponseOff = protocol.CommonResponseEnd(protocolBuilder)

	for i = 0; i < len(backends); i++ {
		peerOff = protocolBuilder.CreateByteVector(backends[i].Ptr().PeerID[:])
		addrOff = protocolBuilder.CreateString(backends[i].Ptr().AddressStr())
		protocol.NetBlockBackendStart(protocolBuilder)
		protocol.NetBlockBackendAddPeerID(protocolBuilder, peerOff)
		protocol.NetBlockBackendAddAddress(protocolBuilder, addrOff)
		backendOffs[i] = protocol.NetBlockBackendEnd(protocolBuilder)
	}

	protocol.NetINodeNetBlockInfoResponseStartBackendsVector(protocolBuilder, len(backends))
	for i = len(backends) - 1; i >= 0; i-- {
		protocolBuilder.PrependUOffsetT(backendOffs[i])
	}
	backendOff = protocolBuilder.EndVector(len(backends))

	protocol.NetINodeNetBlockInfoResponseStart(protocolBuilder)
	protocol.NetINodeNetBlockInfoResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	protocol.NetINodeNetBlockInfoResponseAddBackends(protocolBuilder, backendOff)
	protocol.NetINodeNetBlockInfoResponseAddLen(protocolBuilder, netBlockLen)
	protocol.NetINodeNetBlockInfoResponseAddCap(protocolBuilder, netBlockCap)
	protocolBuilder.Finish(protocol.NetINodeNetBlockInfoResponseEnd(protocolBuilder))
}
