package solonn

import (
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
	"soloos/common/snettypes"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *SRPCServer) NetBlockPrepareMetaData(serviceReq *snettypes.NetQuery) error {
	var (
		param           = make([]byte, serviceReq.BodySize)
		req             solofsprotocol.NetINodeNetBlockInfoRequest
		uNetINode       solofsapitypes.NetINodeUintptr
		netINodeID      solofsapitypes.NetINodeID
		uNetBlock       solofsapitypes.NetBlockUintptr
		protocolBuilder flatbuffers.Builder
		err             error
	)

	err = serviceReq.ReadAll(param)
	if err != nil {
		return err
	}

	// request
	req.Init(param, flatbuffers.GetUOffsetT(param))
	copy(netINodeID[:], req.NetINodeID())
	uNetINode, err = p.solonn.netINodeDriver.GetNetINode(netINodeID)
	defer p.solonn.netINodeDriver.ReleaseNetINode(uNetINode)

	if err != nil {
		if err == solofsapitypes.ErrObjectNotExists {
			solofsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_404, err.Error())
		} else {
			solofsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		}
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	// response
	uNetBlock, err = p.solonn.netBlockDriver.MustGetNetBlock(uNetINode, req.NetBlockIndex())
	defer p.solonn.netBlockDriver.ReleaseNetBlock(uNetBlock)
	if err != nil {
		solofsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	solofsapi.SetNetINodeNetBlockInfoResponse(&protocolBuilder,
		uNetBlock.Ptr().StorDataBackends.Slice(), req.Cap(), req.Cap())
	err = serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])

SERVICE_DONE:
	return err
}
