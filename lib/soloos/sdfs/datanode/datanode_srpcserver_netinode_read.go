package datanode

import (
	sdbapitypes "soloos/common/sdbapi/types"
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeSRPCServer) NetINodePRead(serviceReq snettypes.ServiceRequest) error {
	var (
		reqParamData = make([]byte, serviceReq.ReqParamSize)
		reqParam     protocol.NetINodePWriteRequest
		uNetBlock    types.NetBlockUintptr
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
		protocolBuilder    flatbuffers.Builder
		netINodeID         types.NetINodeID
		uNetINode          types.NetINodeUintptr
		firstNetBlockIndex int32
		lastNetBlockIndex  int32
		netBlockIndex      int32
		respBody           []byte
		readDataSize       int
	)
	copy(netINodeID[:], reqParam.NetINodeID())
	uNetINode, err = p.dataNode.netINodeDriver.GetNetINode(netINodeID)
	defer p.dataNode.netINodeDriver.ReleaseNetINode(uNetINode)
	if err != nil {
		if err == types.ErrObjectNotExists {
			api.SetNetINodePReadResponseError(&protocolBuilder, snettypes.CODE_404, "")
			goto SERVICE_REQUEST_DONE
		} else {
			api.SetNetINodePReadResponseError(&protocolBuilder, snettypes.CODE_502, "")
			goto SERVICE_REQUEST_DONE
		}
	}

	// calculate file data size
	if reqParam.Offset()+uint64(reqParam.Length()) > uNetINode.Ptr().Size {
		readDataSize = int(uNetINode.Ptr().Size - reqParam.Offset())
	} else {
		readDataSize = int(reqParam.Length())
	}

	// prepare uNetBlock
	firstNetBlockIndex = int32(reqParam.Offset() / uint64(uNetINode.Ptr().NetBlockCap))
	lastNetBlockIndex = int32((reqParam.Offset() + uint64(readDataSize)) / uint64(uNetINode.Ptr().NetBlockCap))
	for netBlockIndex = firstNetBlockIndex; netBlockIndex <= lastNetBlockIndex; netBlockIndex++ {
		uNetBlock, err = p.dataNode.netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
		defer p.dataNode.netBlockDriver.ReleaseNetBlock(uNetBlock)
		if err != nil {
			api.SetNetINodePReadResponseError(&protocolBuilder, snettypes.CODE_502, "")
			goto SERVICE_REQUEST_DONE
		}

		if uNetBlock.Ptr().IsLocalDataBackendInited.Load() == sdbapitypes.MetaDataStateUninited {
			p.dataNode.metaStg.PrepareNetBlockLocalDataBackendWithLock(uNetBlock, p.dataNode.uLocalDiskPeer)
		}
	}

SERVICE_REQUEST_DONE:
	if err != nil {
		serviceReq.Conn.SkipReadRemaining()
		return nil
	}

	// request file data
	api.SetNetINodePReadResponse(&protocolBuilder, int32(readDataSize))
	respBody = protocolBuilder.Bytes[protocolBuilder.Head():]
	// TODO set write length
	err = serviceReq.Conn.ResponseHeaderParam(serviceReq.ReqID, respBody, int(readDataSize))
	if err != nil {
		goto SERVICE_RESPONSE_DONE
	}

	// TODO get readedDataLength
	_, err = p.dataNode.netINodeDriver.PReadWithConn(uNetINode, serviceReq.Conn,
		int(readDataSize), reqParam.Offset())
	if err != nil {
		goto SERVICE_RESPONSE_DONE
	}

SERVICE_RESPONSE_DONE:
	return nil
}
