package datanode

import (
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeSRPCServer) NetINodePRead(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	var (
		reqParamData = make([]byte, reqParamSize)
		reqParam     protocol.NetINodePWriteRequest
		uNetBlock    types.NetBlockUintptr
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
		protocolBuilder    flatbuffers.Builder
		netINodeID         types.NetINodeID
		uNetINode          types.NetINodeUintptr
		firstNetBlockIndex int
		lastNetBlockIndex  int
		netBlockIndex      int
		respBody           []byte
		readDataSize       int
	)
	copy(netINodeID[:], reqParam.NetINodeID())
	uNetINode, err = p.dataNode.netINodeDriver.GetNetINode(netINodeID)
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
	if reqParam.Offset()+int64(reqParam.Length()) > uNetINode.Ptr().Size {
		readDataSize = int(uNetINode.Ptr().Size - reqParam.Offset())
	} else {
		readDataSize = int(reqParam.Length())
	}

	// prepare uNetBlock
	firstNetBlockIndex = int(reqParam.Offset() / int64(uNetINode.Ptr().NetBlockCap))
	lastNetBlockIndex = int((reqParam.Offset() + int64(readDataSize)) / int64(uNetINode.Ptr().NetBlockCap))
	for netBlockIndex = firstNetBlockIndex; netBlockIndex <= lastNetBlockIndex; netBlockIndex++ {
		uNetBlock, err = p.dataNode.netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
		if err != nil {
			api.SetNetINodePReadResponseError(&protocolBuilder, snettypes.CODE_502, "")
			goto SERVICE_REQUEST_DONE
		}

		if uNetBlock.Ptr().IsLocalDataBackendInited == false {
			p.dataNode.metaStg.PrepareNetBlockLocalDataBackendWithLock(uNetBlock, p.dataNode.uLocalDiskPeer)
		}
	}

SERVICE_REQUEST_DONE:
	if err != nil {
		conn.SkipReadRemaining()
		return nil
	}

	// request file data
	api.SetNetINodePReadResponse(&protocolBuilder, int32(readDataSize))
	respBody = protocolBuilder.Bytes[protocolBuilder.Head():]
	// TODO set write length
	err = conn.ResponseHeaderParam(reqID, respBody, int(readDataSize))
	if err != nil {
		goto SERVICE_RESPONSE_DONE
	}

	// TODO get readedDataLength
	_, err = p.dataNode.netINodeDriver.PReadWithConn(uNetINode, conn,
		int(readDataSize), reqParam.Offset())
	if err != nil {
		goto SERVICE_RESPONSE_DONE
	}

SERVICE_RESPONSE_DONE:
	return nil
}
