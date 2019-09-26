package solonn

import (
	"soloos/common/log"
	"soloos/common/solofsapi"
	"soloos/common/solofsprotocol"
	"soloos/common/snettypes"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *SRPCServer) SolodnRegister(serviceReq *snettypes.NetQuery) error {
	var (
		param           = make([]byte, serviceReq.BodySize)
		req             solofsprotocol.SNetPeer
		protocolBuilder flatbuffers.Builder
		err             error
	)

	err = serviceReq.ReadAll(param)
	if err != nil {
		return err
	}

	var (
		peer snettypes.Peer
	)
	req.Init(param, flatbuffers.GetUOffsetT(param))
	copy(peer.ID[:], req.PeerID())
	peer.SetAddressBytes(req.Address())
	peer.ServiceProtocol.SetProtocolBytes(req.Protocol())

	err = p.solonn.SolodnRegister(peer)
	log.Info("solodn resgister:", peer.PeerIDStr(), peer.AddressStr())
	if err != nil {
		solofsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	solofsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])

SERVICE_DONE:
	return nil
}
