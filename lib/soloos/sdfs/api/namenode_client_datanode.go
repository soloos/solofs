package api

import (
	"soloos/sdfs/protocol"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeClient) RegisterDataNode(peerID snettypes.PeerID, serveAddr string) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		peerIDOff       flatbuffers.UOffsetT
		addrOff         flatbuffers.UOffsetT
		err             error
	)

	peerIDOff = protocolBuilder.CreateByteString(peerID[:])
	addrOff = protocolBuilder.CreateString(serveAddr)
	protocol.SNetPeerStart(&protocolBuilder)
	protocol.SNetPeerAddPeerID(&protocolBuilder, peerIDOff)
	protocol.SNetPeerAddAddress(&protocolBuilder, addrOff)
	protocolBuilder.Finish(protocol.SNetPeerEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	err = p.snetClientDriver.Call(p.nameNodePeer,
		"/DataNode/Register", &req, &resp)
	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	p.snetClientDriver.ReadResponse(p.nameNodePeer, &req, &resp, body)
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
