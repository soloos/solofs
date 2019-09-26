package solodn

import (
	"soloos/common/solofsapi"
	"soloos/common/solofsapitypes"
	"soloos/common/snettypes"
	"soloos/common/solofsprotocol"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *SRPCServer) NetINodeSync(serviceReq *snettypes.NetQuery) error {
	var (
		reqParamData = make([]byte, serviceReq.ParamSize)
		reqParam     solofsprotocol.NetINodePWriteRequest
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
		netINodeID      solofsapitypes.NetINodeID
		uNetINode       solofsapitypes.NetINodeUintptr
	)
	copy(netINodeID[:], reqParam.NetINodeID())
	uNetINode, err = p.solodn.netINodeDriver.GetNetINode(netINodeID)
	defer p.solodn.netINodeDriver.ReleaseNetINode(uNetINode)
	if err != nil {
		if err == solofsapitypes.ErrObjectNotExists {
			solofsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_404)
			goto SERVICE_DONE
		} else {
			solofsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
			goto SERVICE_DONE
		}
	}

	err = p.solodn.netINodeDriver.Sync(uNetINode)
	if err != nil {
		solofsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
		goto SERVICE_DONE
	}

SERVICE_DONE:
	if err != nil {
		return nil
	}

	if err == nil {
		solofsapi.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	}

	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	err = serviceReq.SimpleResponse(serviceReq.ReqID, respBody)
	if err != nil {
		return err
	}

	return nil

}
