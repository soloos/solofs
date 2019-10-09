package memstg

import (
	"soloos/common/snet"
	"soloos/common/solofsprotocol"
	"soloos/common/solofstypes"
)

func (p *NetBlockDriver) PReadMemBlockFromNet(uNetINode solofstypes.NetINodeUintptr,
	uNetBlock solofstypes.NetBlockUintptr,
	netBlockIndex int32,
	uMemBlock solofstypes.MemBlockUintptr,
	memBlockIndex int32,
	offset uint64, length int,
) (int, error) {
	if uNetBlock.Ptr().IsLocalDataBackendExists {
		return p.helper.PreadMemBlockWithDisk(
			uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, offset, length)
	}

	var peer snet.Peer
	var err error
	peer, err = p.SNetDriver.GetPeer(uNetBlock.Ptr().StorDataBackends.Arr[0])
	if err != nil {
		return 0, err
	}

	switch peer.ServiceProtocol {
	case snet.ProtocolSolofs:
		return p.doPReadMemBlockWithSrpc(peer.ID,
			uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, offset, length)
	}

	return 0, solofstypes.ErrServiceNotExists
}

func (p *NetBlockDriver) doPReadMemBlockWithSrpc(peerID snet.PeerID,
	uNetINode solofstypes.NetINodeUintptr,
	uNetBlock solofstypes.NetBlockUintptr,
	netBlockIndex int32,
	uMemBlock solofstypes.MemBlockUintptr,
	memBlockIndex int32,
	offset uint64, length int,
) (int, error) {
	var (
		snetReq  snet.SNetReq
		snetResp snet.SNetResp
		err      error
	)

	var req solofsprotocol.NetINodePReadReq
	req.NetINodeID = uNetBlock.Ptr().NetINodeID
	req.Offset = offset
	req.Length = int32(length)

	// TODO choose solodn
	err = p.solodnClient.Call(peerID, "/NetINode/PRead", &snetReq, &snetResp, req)
	if err != nil {
		return 0, err
	}

	var (
		resp                        solofsprotocol.NetINodePReadResp
		offsetInMemBlock, readedLen int
	)
	err = p.solodnClient.SimpleReadResponse(peerID, &snetReq, &snetResp, &resp)
	if err != nil {
		return 0, err
	}

	offsetInMemBlock = int(offset - uint64(uMemBlock.Ptr().Bytes.Cap)*uint64(memBlockIndex))
	readedLen = int(snetResp.ConnBytesLeft)
	err = p.solodnClient.ReadResponse(peerID, &snetReq, &snetResp,
		(*uMemBlock.Ptr().BytesSlice())[offsetInMemBlock:readedLen])
	if err != nil {
		return 0, err
	}

	return int(resp.Length), err
}
