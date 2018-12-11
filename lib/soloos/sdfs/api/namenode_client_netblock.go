package api

import (
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (p *NameNodeClient) PrepareNetBlockMetadata(uINode types.INodeUintptr,
	netBlockIndex int,
	uNetBlock types.NetBlockUintptr,
) error {
	var (
		req             snettypes.Request
		resp            snettypes.Response
		protocolBuilder flatbuffers.Builder
		err             error
	)

	inodeIDOff := protocolBuilder.CreateString(uINode.Ptr().IDStr())
	protocol.INodeNetBlockInfoRequestStart(&protocolBuilder)
	protocol.INodeNetBlockInfoRequestAddInodeID(&protocolBuilder, inodeIDOff)
	protocol.INodeNetBlockInfoRequestAddNetBlockIndex(&protocolBuilder, int32(netBlockIndex))
	protocol.INodeNetBlockInfoRequestAddCap(&protocolBuilder, int32(uINode.Ptr().NetBlockCap))
	protocolBuilder.Finish(protocol.INodeNetBlockInfoRequestEnd(&protocolBuilder))
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
		netBlockInfo   protocol.INodeNetBlockInfoResponse
		commonResponse protocol.CommonResponse
		backend        protocol.NetBlockBackend
		peerID         snettypes.PeerID
		uPeer          snettypes.PeerUintptr
		i              int
	)
	netBlockInfo.Init(body, flatbuffers.GetUOffsetT(body))
	netBlockInfo.CommonResponse(&commonResponse)
	if commonResponse.Code() != snettypes.CODE_OK {
		return types.ErrRemoteService
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
