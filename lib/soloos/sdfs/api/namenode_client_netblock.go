package api

import (
	"soloos/log"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeClient) PrepareNetBlockMetadata(uNetINode types.NetINodeUintptr,
	netBlockIndex int,
	uNetBlock types.NetBlockUintptr,
) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		err             error
	)

	netINodeIDOff := protocolBuilder.CreateString(uNetINode.Ptr().IDStr())
	protocol.NetINodeNetBlockInfoRequestStart(&protocolBuilder)
	protocol.NetINodeNetBlockInfoRequestAddInodeID(&protocolBuilder, netINodeIDOff)
	protocol.NetINodeNetBlockInfoRequestAddNetBlockIndex(&protocolBuilder, int32(netBlockIndex))
	protocol.NetINodeNetBlockInfoRequestAddCap(&protocolBuilder, int32(uNetINode.Ptr().NetBlockCap))
	protocolBuilder.Finish(protocol.NetINodeNetBlockInfoRequestEnd(&protocolBuilder))
	req.Param = protocolBuilder.Bytes[protocolBuilder.Head():]

	// TODO choose namenode
	err = p.snetClientDriver.Call(p.nameNodePeer,
		"/NetBlock/PrepareMetadata", &req, &resp)
	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	p.snetClientDriver.ReadResponse(p.nameNodePeer, &req, &resp, body)
	if err != nil {
		return err
	}

	var (
		pNetBlock      = uNetBlock.Ptr()
		netBlockInfo   protocol.NetINodeNetBlockInfoResponse
		commonResponse protocol.CommonResponse
		backend        protocol.NetBlockBackend
		peerID         snettypes.PeerID
		uPeer          snettypes.PeerUintptr
		i              int
	)
	netBlockInfo.Init(body, flatbuffers.GetUOffsetT(body))
	netBlockInfo.CommonResponse(&commonResponse)
	if commonResponse.Code() != snettypes.CODE_OK {
		if commonResponse.Code() == snettypes.CODE_404 {
			return types.ErrObjectNotExists
		} else {
			log.Warn(string(commonResponse.Error()))
			return types.ErrRemoteService
		}
	}

	copy(pNetBlock.ID[:], netBlockInfo.NetBlockID())
	pNetBlock.IndexInInode = netBlockIndex
	pNetBlock.Len = int(netBlockInfo.Len())
	pNetBlock.Cap = int(netBlockInfo.Cap())
	pNetBlock.DataNodes.Reset()
	for i = 0; i < netBlockInfo.BackendsLength(); i++ {
		netBlockInfo.Backends(&backend, i)
		copy(peerID[:], netBlockInfo.NetBlockID())
		uPeer = p.snetDriver.MustGetPeer(&peerID, string(backend.Address()), types.DefaultSDFSRPCProtocol)
		pNetBlock.DataNodes.Append(uPeer)
	}

	return nil
}
