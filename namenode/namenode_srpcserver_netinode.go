package namenode

import (
	"soloos/common/log"
	"soloos/common/sdfsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/sdfsprotocol"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeSRPCServer) doNetINodeGet(isMustGet bool, serviceReq *snettypes.NetQuery) error {
	var (
		param           = make([]byte, serviceReq.BodySize)
		req             sdfsprotocol.NetINodeInfoRequest
		uNetINode       sdfsapitypes.NetINodeUintptr
		netINodeID      sdfsapitypes.NetINodeID
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
		uNetINode, err = p.nameNode.netINodeDriver.MustGetNetINode(netINodeID,
			req.Size(), int(req.NetBlockCap()), int(req.MemBlockCap()))
	} else {
		uNetINode, err = p.nameNode.netINodeDriver.GetNetINode(netINodeID)
	}
	defer p.nameNode.netINodeDriver.ReleaseNetINode(uNetINode)

	if err != nil {
		if err == sdfsapitypes.ErrObjectNotExists {
			sdfsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_404, err.Error())
			serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
			err = nil
			goto SERVICE_DONE
		}

		log.Info("get netinode from db error:", err, string(netINodeID[:]))
		sdfsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	// response
	sdfsapi.SetNetINodeInfoResponse(&protocolBuilder,
		uNetINode.Ptr().Size, int32(uNetINode.Ptr().NetBlockCap), int32(uNetINode.Ptr().MemBlockCap))
	serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
	err = nil

SERVICE_DONE:
	return err
}

func (p *NameNodeSRPCServer) NetINodeGet(serviceReq *snettypes.NetQuery) error {
	return p.doNetINodeGet(false, serviceReq)
}

func (p *NameNodeSRPCServer) NetINodeMustGet(serviceReq *snettypes.NetQuery) error {
	return p.doNetINodeGet(true, serviceReq)
}

func (p *NameNodeSRPCServer) NetINodeCommitSizeInDB(serviceReq *snettypes.NetQuery) error {
	var (
		param           = make([]byte, serviceReq.BodySize)
		req             sdfsprotocol.NetINodeCommitSizeInDBRequest
		uNetINode       sdfsapitypes.NetINodeUintptr
		netINodeID      sdfsapitypes.NetINodeID
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
	uNetINode, err = p.nameNode.netINodeDriver.GetNetINode(netINodeID)
	defer p.nameNode.netINodeDriver.ReleaseNetINode(uNetINode)

	if err != nil {
		sdfsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		err = nil
		goto SERVICE_DONE
	}

	err = p.nameNode.metaStg.NetINodeDriver.NetINodeTruncate(uNetINode, req.Size())
	if err != nil {
		sdfsapi.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		err = nil
		goto SERVICE_DONE
	}

	// response
	sdfsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	serviceReq.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])

SERVICE_DONE:
	return err
}
