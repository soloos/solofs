package api

import (
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

type NameNodeClient struct {
	snetDriver       *snet.SNetDriver
	snetClientDriver *snet.ClientDriver
	nameNodePeer     snettypes.PeerUintptr
}

func (p *NameNodeClient) Init(snetDriver *snet.SNetDriver,
	snetClientDriver *snet.ClientDriver,
	nameNodePeer snettypes.PeerUintptr) error {
	p.snetDriver = snetDriver
	p.snetClientDriver = snetClientDriver
	p.nameNodePeer = nameNodePeer
	return nil
}

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
		"/NetBlock/MustGet", &req, &resp)
	var body = make([]byte, resp.BodySize)[:resp.BodySize]
	p.snetClientDriver.ReadResponse(p.nameNodePeer, &req, &resp, body)
	if err != nil {
		return err
	}

	var (
		pNetBlock    = uNetBlock.Ptr()
		netBlockInfo protocol.INodeNetBlockInfoResponse
		backend      protocol.NetBlockBackend
		peerID       snettypes.PeerID
		uPeer        snettypes.PeerUintptr
		i            int
	)
	netBlockInfo.Init(body, flatbuffers.GetUOffsetT(body))
	copy(pNetBlock.ID[:], netBlockInfo.NetBlockID())
	pNetBlock.IndexInInode = netBlockIndex
	pNetBlock.Len = netBlockInfo.Len()
	pNetBlock.Cap = netBlockInfo.Cap()
	pNetBlock.DataNodes.Reset()
	for i = 0; i < netBlockInfo.BackendsLength(); i++ {
		netBlockInfo.Backends(&backend, i)
		copy(peerID[:], netBlockInfo.NetBlockID())
		uPeer = p.snetDriver.MustGetPeer(&peerID, string(backend.Address()), types.DefaultSDFSRPCProtocol)
		pNetBlock.DataNodes.Append(uPeer)
	}

	pNetBlock.IsMetaDataInited = true

	return nil
}
