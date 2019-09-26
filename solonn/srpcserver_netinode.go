package solonn

import (
	"soloos/common/log"
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/snettypes"
	"soloos/common/solofsprotocol"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *SRPCServer) doNetINodeGet(isMustGet bool, serviceReq *snettypes.NetQuery) error {
	var (
		param           = make([]byte, serviceReq.BodySize)
		req             solofsprotocol.NetINodeInfoRequest
		uNetINode       solofsapitypes.NetINodeUintptr
		netINodeID      solofsapitypes.NetINodeID
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
	if isMustGet {
		uNetINode, err = p.solonn.netINodeDriver.MustGetNetINode(netINodeID,
			req.Size(), int(req.NetBlockCap()), int(req.MemBlockCap()))
	} else {
		uNetINode, err = p.solonn.netINodeDriver.GetNetINode(netINodeID)
	}
	defer p.solonn.netINodeDriver.ReleaseNetINode(uNetINode)

	if err != nil {
		if err == solofsapitypes.ErrObjectNotExists {
			solofsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_404, err.Error())
			serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
			err = nil
			goto SERVICE_DONE
		}

		log.Info("get netinode from db error:", err, string(netINodeID[:]))
		solofsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	// response
	solofsapi.SetNetINodeInfoResponse(&protocolBuilder,
		uNetINode.Ptr().Size, int32(uNetINode.Ptr().NetBlockCap), int32(uNetINode.Ptr().MemBlockCap))
	serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
	err = nil

SERVICE_DONE:
	return err
}

func (p *SRPCServer) NetINodeGet(serviceReq *snettypes.NetQuery) error {
	return p.doNetINodeGet(false, serviceReq)
}

func (p *SRPCServer) NetINodeMustGet(serviceReq *snettypes.NetQuery) error {
	return p.doNetINodeGet(true, serviceReq)
}

func (p *SRPCServer) NetINodeCommitSizeInDB(serviceReq *snettypes.NetQuery) error {
	var (
		param           = make([]byte, serviceReq.BodySize)
		req             solofsprotocol.NetINodeCommitSizeInDBRequest
		uNetINode       solofsapitypes.NetINodeUintptr
		netINodeID      solofsapitypes.NetINodeID
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
		solofsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		err = nil
		goto SERVICE_DONE
	}

	err = p.solonn.metaStg.NetINodeDriver.NetINodeTruncate(uNetINode, req.Size())
	if err != nil {
		solofsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		err = nil
		goto SERVICE_DONE
	}

	// response
	solofsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])

SERVICE_DONE:
	return err
}
