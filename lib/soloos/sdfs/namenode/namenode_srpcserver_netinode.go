package namenode

import (
	"soloos/common/log"
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeSRPCServer) doNetINodeGet(isMustGet bool, serviceReq snettypes.ServiceRequest) error {
	var (
		param           = make([]byte, serviceReq.ReqBodySize)
		req             protocol.NetINodeInfoRequest
		uNetINode       types.NetINodeUintptr
		netINodeID      types.NetINodeID
		protocolBuilder flatbuffers.Builder
		err             error
	)

	err = serviceReq.Conn.ReadAll(param)
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
		if err == types.ErrObjectNotExists {
			api.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_404, err.Error())
			serviceReq.Conn.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
			err = nil
			goto SERVICE_DONE
		}

		log.Info("get netinode from db error:", err, string(netINodeID[:]))
		api.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		serviceReq.Conn.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	// response
	api.SetNetINodeInfoResponse(&protocolBuilder,
		uNetINode.Ptr().Size, int32(uNetINode.Ptr().NetBlockCap), int32(uNetINode.Ptr().MemBlockCap))
	serviceReq.Conn.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
	err = nil

SERVICE_DONE:
	return err
}

func (p *NameNodeSRPCServer) NetINodeGet(serviceReq snettypes.ServiceRequest) error {
	return p.doNetINodeGet(false, serviceReq)
}

func (p *NameNodeSRPCServer) NetINodeMustGet(serviceReq snettypes.ServiceRequest) error {
	return p.doNetINodeGet(true, serviceReq)
}

func (p *NameNodeSRPCServer) NetINodeCommitSizeInDB(serviceReq snettypes.ServiceRequest) error {
	var (
		param           = make([]byte, serviceReq.ReqBodySize)
		req             protocol.NetINodeCommitSizeInDBRequest
		uNetINode       types.NetINodeUintptr
		netINodeID      types.NetINodeID
		protocolBuilder flatbuffers.Builder
		err             error
	)

	err = serviceReq.Conn.ReadAll(param)
	if err != nil {
		return err
	}

	// request
	req.Init(param, flatbuffers.GetUOffsetT(param))

	copy(netINodeID[:], req.NetINodeID())
	uNetINode, err = p.nameNode.netINodeDriver.GetNetINode(netINodeID)
	defer p.nameNode.netINodeDriver.ReleaseNetINode(uNetINode)

	if err != nil {
		api.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		serviceReq.Conn.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		err = nil
		goto SERVICE_DONE
	}

	err = p.nameNode.metaStg.NetINodeDriver.NetINodeTruncate(uNetINode, req.Size())
	if err != nil {
		api.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		serviceReq.Conn.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		err = nil
		goto SERVICE_DONE
	}

	// response
	api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	serviceReq.Conn.SimpleResponse(serviceReq.ReqID, protocolBuilder.Bytes[protocolBuilder.Head():])

SERVICE_DONE:
	return err
}
