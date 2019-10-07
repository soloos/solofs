package solonn

import (
	"soloos/common/log"
	"soloos/common/snet"
	"soloos/common/solofsprotocol"
)

func (p *SrpcServer) SolodnRegister(reqCtx *snet.SNetReqContext,
	req solofsprotocol.SNetPeer,
) error {
	var (
		peer snet.Peer
		err  error
	)
	copy(peer.ID[:], []byte(req.PeerID))
	peer.SetAddress(req.Address)
	peer.ServiceProtocol.SetProtocolBytes([]byte(req.Protocol))

	err = p.solonn.SolodnRegister(peer)
	log.Info("solodn resgister:", peer.PeerIDStr(), peer.AddressStr())
	if err != nil {
		return err
	}

	return nil
}
