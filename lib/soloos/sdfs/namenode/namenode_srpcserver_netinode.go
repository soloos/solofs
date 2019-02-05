package namenode

import (
	"soloos/log"
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeSRPCServer) doNetINodeGet(isMustGet bool,
	reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	var (
		param           = make([]byte, reqBodySize)
		req             protocol.NetINodeInfoRequest
		uNetINode       types.NetINodeUintptr
		netINodeID      types.NetINodeID
		protocolBuilder flatbuffers.Builder
		err             error
	)

	err = conn.ReadAll(param)
	if err != nil {
		return err
	}

	// request
	req.Init(param, flatbuffers.GetUOffsetT(param))

	copy(netINodeID[:], req.NetINodeID())
	if isMustGet {
		uNetINode, err = p.nameNode.netINodeDriver.MustGetNetINodeWithReadAcquire(netINodeID,
			req.Size(), int(req.NetBlockCap()), int(req.MemBlockCap()))
	} else {
		uNetINode, err = p.nameNode.netINodeDriver.GetNetINodeWithReadAcquire(false, netINodeID)
	}

	if err != nil {
		if err == types.ErrObjectNotExists {
			api.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_404, err.Error())
			conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
			err = nil
			goto SERVICE_DONE
		}

		log.Info("get netinode from db error:", err, string(netINodeID[:]))
		api.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		goto SERVICE_DONE
	}

	// response
	api.SetNetINodeInfoResponse(&protocolBuilder,
		uNetINode.Ptr().Size, int32(uNetINode.Ptr().NetBlockCap), int32(uNetINode.Ptr().MemBlockCap))
	conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
	err = nil

SERVICE_DONE:
	return err
}

func (p *NameNodeSRPCServer) NetINodeGet(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	return p.doNetINodeGet(false, reqID, reqBodySize, reqParamSize, conn)
}

func (p *NameNodeSRPCServer) NetINodeMustGet(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	return p.doNetINodeGet(true, reqID, reqBodySize, reqParamSize, conn)
}

func (p *NameNodeSRPCServer) NetINodeCommitSizeInDB(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	var (
		param           = make([]byte, reqBodySize)
		req             protocol.NetINodeCommitSizeInDBRequest
		uNetINode       types.NetINodeUintptr
		netINodeID      types.NetINodeID
		protocolBuilder flatbuffers.Builder
		err             error
	)

	err = conn.ReadAll(param)
	if err != nil {
		return err
	}

	// request
	req.Init(param, flatbuffers.GetUOffsetT(param))

	copy(netINodeID[:], req.NetINodeID())
	uNetINode, err = p.nameNode.netINodeDriver.GetNetINodeWithReadAcquire(false, netINodeID)
	if err != nil {
		api.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		err = nil
		goto SERVICE_DONE
	}

	err = p.nameNode.metaStg.NetINodeDriver.NetINodeTruncate(uNetINode, req.Size())
	if err != nil {
		api.SetNetINodeNetBlockInfoResponseError(&protocolBuilder, snettypes.CODE_502, err.Error())
		conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])
		err = nil
		goto SERVICE_DONE
	}

	// response
	api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])

SERVICE_DONE:
	return err
}
