package datanode

import (
	"soloos/common/sdfsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/sdfsprotocol"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *SRPCServer) NetINodeSync(serviceReq *snettypes.NetQuery) error {
	var (
		reqParamData = make([]byte, serviceReq.ParamSize)
		reqParam     sdfsprotocol.NetINodePWriteRequest
		err          error
	)

	// request param
	err = serviceReq.ReadAll(reqParamData)
	if err != nil {
		return err
	}
	reqParam.Init(reqParamData[:serviceReq.ParamSize],
		flatbuffers.GetUOffsetT(reqParamData[:serviceReq.ParamSize]))

	// response

	// get uNetINode
	var (
		protocolBuilder flatbuffers.Builder
		netINodeID      sdfsapitypes.NetINodeID
		uNetINode       sdfsapitypes.NetINodeUintptr
	)
	copy(netINodeID[:], reqParam.NetINodeID())
	uNetINode, err = p.dataNode.netINodeDriver.GetNetINode(netINodeID)
	defer p.dataNode.netINodeDriver.ReleaseNetINode(uNetINode)
	if err != nil {
		if err == sdfsapitypes.ErrObjectNotExists {
			sdfsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_404)
			goto SERVICE_DONE
		} else {
			sdfsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
			goto SERVICE_DONE
		}
	}

	err = p.dataNode.netINodeDriver.Sync(uNetINode)
	if err != nil {
		sdfsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
		goto SERVICE_DONE
	}

SERVICE_DONE:
	if err != nil {
		return nil
	}

	if err == nil {
		sdfsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	}

	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	err = serviceReq.SimpleResponse(serviceReq.ReqID, respBody)
	if err != nil {
		return err
	}

	return nil

}
