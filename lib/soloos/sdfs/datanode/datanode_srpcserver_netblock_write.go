package datanode

import (
	"soloos/log"
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
	"soloos/util"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *DataNodeSRPCServer) NetBlockPWrite(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {
	var (
		reqParamData         = make([]byte, reqParamSize)
		reqParam             protocol.NetBlockPWriteRequest
		syncDataProtoBackend protocol.NetBlockBackend
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

	log.Error(string(reqParam.NetINodeID()))
	log.Error(reqParam.NetBlockIndex())
	log.Error(reqParam.Offset())
	log.Error(reqParam.Length())

	// response

	// get uNetINode
	var (
		protocolBuilder flatbuffers.Builder
		netINodeID      types.NetINodeID
		uNetINode       types.NetINodeUintptr
	)
	copy(netINodeID[:], reqParam.NetINodeID())
	uNetINode, err = p.dataNode.metaStg.GetNetINode(netINodeID)
	if err != nil {
		if err == types.ErrObjectNotExists {
			err = api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_404)
			goto SERVICE_DONE
		} else {
			err = api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
			goto SERVICE_DONE
		}
	}

	// get uNetBlock
	uNetBlock, err = p.dataNode.metaStg.MustGetNetBlock(uNetINode, int(reqParam.NetBlockIndex()))
	if err != nil {
		err = api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_502)
		goto SERVICE_DONE
	}

	if uNetBlock.Ptr().IsSyncDataBackendsInited == false {
		syncDataBackends.Append(p.dataNode.uLocalDiskPeer)
		for i = 0; i < reqParam.TransferBackendsLength(); i++ {
			reqParam.TransferBackends(&syncDataProtoBackend, i)
			copy(peerID[:], syncDataProtoBackend.PeerID())
			uPeer, _ = p.dataNode.snetDriver.MustGetPeer(&peerID,
				string(syncDataProtoBackend.Address()), types.DefaultSDFSRPCProtocol)
			syncDataBackends.Append(uPeer)
		}
		p.dataNode.metaStg.PrepareNetBlockSyncDataBackendsWithLock(uNetBlock, syncDataBackends)
	}

	// TODO pwrite

SERVICE_DONE:
	if err != nil {
		conn.SkipReadRemaining()
		return nil
	}

	// request file data
	// TODO write to disk or transfer
	util.Ignore(uNetINode)
	var tmp = make([]byte, reqBodySize-reqParamSize)
	err = conn.ReadAll(tmp)
	if err != nil {
		return err
	}

	api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	err = conn.SimpleResponse(reqID, respBody)
	if err != nil {
		return err
	}

	return nil
}
