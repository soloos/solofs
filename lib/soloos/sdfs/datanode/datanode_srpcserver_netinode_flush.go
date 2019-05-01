package datanode

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeSRPCServer) NetINodeFlush(serviceReq snettypes.ServiceRequest) error {
	var (
		reqParamData = make([]byte, serviceReq.ReqParamSize)
		reqParam     protocol.NetINodePWriteRequest
		err          error
	)

	// request param
	err = serviceReq.Conn.ReadAll(reqParamData)
	if err != nil {
		return err
	}
	reqParam.Init(reqParamData[:serviceReq.ReqParamSize],
		flatbuffers.GetUOffsetT(reqParamData[:serviceReq.ReqParamSize]))

	// response

	// get uNetINode
	var (
		protocolBuilder flatbuffers.Builder
		netINodeID      types.NetINodeID
		uNetINode       types.NetINodeUintptr
	)
	copy(netINodeID[:], reqParam.NetINodeID())
	uNetINode, err = p.dataNode.netINodeDriver.GetNetINodeWithReadAcquire(false, netINodeID)
	defer p.dataNode.netINodeDriver.ReleaseNetINodeWithReadRelease(uNetINode)
	if err != nil {
		if err == types.ErrObjectNotExists {
			api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_404)
			goto SERVICE_DONE
		} else {
			api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
			goto SERVICE_DONE
		}
	}

	err = p.dataNode.netINodeDriver.Flush(uNetINode)
	if err != nil {
		api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
		goto SERVICE_DONE
	}

SERVICE_DONE:
	if err != nil {
		serviceReq.Conn.SkipReadRemaining()
		return nil
	}

	if err == nil {
		api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	}

	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	err = serviceReq.Conn.SimpleResponse(serviceReq.ReqID, respBody)
	if err != nil {
		return err
	}

	return nil

}
