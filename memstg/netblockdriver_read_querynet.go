package memstg

import (
	"soloos/common/snet"
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
)

func (p *NetBlockDriver) PReadMemBlockFromNet(uNetINode solofsapitypes.NetINodeUintptr,
	uNetBlock solofsapitypes.NetBlockUintptr,
	netBlockIndex int32,
	uMemBlock solofsapitypes.MemBlockUintptr,
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

	return 0, solofsapitypes.ErrServiceNotExists
}

func (p *NetBlockDriver) doPReadMemBlockWithSrpc(peerID snet.PeerID,
	uNetINode solofsapitypes.NetINodeUintptr,
	uNetBlock solofsapitypes.NetBlockUintptr,
	netBlockIndex int32,
	uMemBlock solofsapitypes.MemBlockUintptr,
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
	snetReq.Param = snet.MustSpecMarshalRequest(req)
	err = p.solodnClient.Call(peerID, "/NetINode/PRead", &snetReq, &snetResp)
	if err != nil {
		return 0, err
	}

	var (
		respParamBs                 = make([]byte, snetResp.ParamSize)
		resp                        solofsprotocol.NetINodePReadResp
		offsetInMemBlock, readedLen int
	)
	err = p.solodnClient.ReadResponse(peerID, &snetReq, &snetResp, respParamBs, &resp)
	if err != nil {
		return 0, err
	}

	offsetInMemBlock = int(offset - uint64(uMemBlock.Ptr().Bytes.Cap)*uint64(memBlockIndex))
	readedLen = int(snetResp.ConnBytesLeft)
	err = p.solodnClient.ReadRawResponse(peerID, &snetReq, &snetResp,
		(*uMemBlock.Ptr().BytesSlice())[offsetInMemBlock:readedLen])
	if err != nil {
		return 0, err
	}

	return int(resp.Length), err
}
