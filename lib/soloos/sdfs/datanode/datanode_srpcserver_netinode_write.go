package datanode

import (
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeSRPCServer) NetINodePWrite(serviceReq snettypes.ServiceRequest) error {
	var (
		reqParamData         = make([]byte, serviceReq.ReqParamSize)
		reqParam             protocol.NetINodePWriteRequest
		syncDataProtoBackend protocol.SNetPeer
		syncDataBackends     snettypes.PeerUintptrArray8
		uPeer                snettypes.PeerUintptr
		peerID               snettypes.PeerID
		uNetBlock            types.NetBlockUintptr
		i                    int
		err                  error
	)

	// request param
	err = serviceReq.Conn.ReadAll(reqParamData)
	if err != nil {
		return err
	}
	reqParam.Init(reqParamData[:serviceReq.ReqParamSize], flatbuffers.GetUOffsetT(reqParamData[:serviceReq.ReqParamSize]))

	// response

	// get uNetINode
	var (
		protocolBuilder    flatbuffers.Builder
		netINodeID         types.NetINodeID
		uNetINode          types.NetINodeUintptr
		firstNetBlockIndex int32
		lastNetBlockIndex  int32
		netBlockIndex      int32
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

	syncDataBackends.Reset()
	syncDataBackends.Append(p.dataNode.uLocalDiskPeer)
	for i = 0; i < reqParam.TransferBackendsLength(); i++ {
		reqParam.TransferBackends(&syncDataProtoBackend, i)
		copy(peerID[:], syncDataProtoBackend.PeerID())
		uPeer, _ = p.dataNode.snetDriver.MustGetPeer(&peerID,
			string(syncDataProtoBackend.Address()), types.DefaultSDFSRPCProtocol)
		syncDataBackends.Append(uPeer)
	}

	// prepare uNetBlock
	firstNetBlockIndex = int32(reqParam.Offset() / uint64(uNetINode.Ptr().NetBlockCap))
	lastNetBlockIndex = int32((reqParam.Offset() + uint64(reqParam.Length())) / uint64(uNetINode.Ptr().NetBlockCap))
	for netBlockIndex = firstNetBlockIndex; netBlockIndex <= lastNetBlockIndex; netBlockIndex++ {
		uNetBlock, err = p.dataNode.netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
		if err != nil {
			api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
			goto SERVICE_DONE
		}

		if uNetBlock.Ptr().IsSyncDataBackendsInited.Load() == types.MetaDataStateUninited {
			p.dataNode.metaStg.PrepareNetBlockSyncDataBackendsWithLock(uNetBlock, syncDataBackends)
		}
	}

	// request file data
	err = p.dataNode.netINodeDriver.PWriteWithConn(uNetINode, serviceReq.Conn,
		int(reqParam.Length()), reqParam.Offset())
	if err != nil {
		return err
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
