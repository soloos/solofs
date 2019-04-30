package namenode

import (
	"soloos/common/log"
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeSRPCServer) DataNodeRegister(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	var (
		param           = make([]byte, reqBodySize)
		req             protocol.SNetPeer
		protocolBuilder flatbuffers.Builder
		err             error
	)

	err = conn.ReadAll(param)
	if err != nil {
		return err
	}

	var (
		peerID snettypes.PeerID
	)
	req.Init(param, flatbuffers.GetUOffsetT(param))
	copy(peerID[:], req.PeerID())
	err = p.nameNode.RegisterDataNode(peerID, string(req.Address()))
	log.Info("datanode resgister:", string(peerID[:]), string(req.Address()))
	if err != nil {
		api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
		conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])

SERVICE_DONE:
	return nil
}
