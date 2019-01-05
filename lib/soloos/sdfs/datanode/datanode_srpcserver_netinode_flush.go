package datanode

import (
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeSRPCServer) NetINodeFlush(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	var (
		reqParamData = make([]byte, reqParamSize)
		reqParam     protocol.NetINodePWriteRequest
		err          error
	)

	// request param
	err = conn.ReadAll(reqParamData)
	if err != nil {
		return err
	}
	reqParam.Init(reqParamData[:reqParamSize], flatbuffers.GetUOffsetT(reqParamData[:reqParamSize]))

	// response

	// get uNetINode
	var (
		protocolBuilder flatbuffers.Builder
		netINodeID      types.NetINodeID
		uNetINode       types.NetINodeUintptr
	)
	copy(netINodeID[:], reqParam.NetINodeID())
	uNetINode, err = p.dataNode.netINodeDriver.GetNetINode(netINodeID)
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
		conn.SkipReadRemaining()
		return nil
	}

	if err == nil {
		api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	}

	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	err = conn.SimpleResponse(reqID, respBody)
	if err != nil {
		return err
	}

	return nil

}
