package datanode

import (
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeSRPCServer) NetINodePWrite(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	var (
		reqParamData         = make([]byte, reqParamSize)
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
	firstNetBlockIndex = int(reqParam.Offset() / uint64(uNetINode.Ptr().NetBlockCap))
	lastNetBlockIndex = int((reqParam.Offset() + uint64(reqParam.Length())) / uint64(uNetINode.Ptr().NetBlockCap))
	for netBlockIndex = firstNetBlockIndex; netBlockIndex <= lastNetBlockIndex; netBlockIndex++ {
		uNetBlock, err = p.dataNode.netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
		if err != nil {
			api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
			goto SERVICE_DONE
		}

		if uNetBlock.Ptr().IsSyncDataBackendsInited == false {
			p.dataNode.metaStg.PrepareNetBlockSyncDataBackendsWithLock(uNetBlock, syncDataBackends)
		}
	}

	// request file data
	err = p.dataNode.netINodeDriver.PWriteWithConn(uNetINode, conn,
		int(reqParam.Length()), reqParam.Offset())
	if err != nil {
		return err
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
