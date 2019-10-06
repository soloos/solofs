package memstg

import (
	"soloos/common/snettypes"
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

	var peer snettypes.Peer
	var err error
	peer, err = p.SNetDriver.GetPeer(uNetBlock.Ptr().StorDataBackends.Arr[0])
	if err != nil {
		return 0, err
	}

	switch peer.ServiceProtocol {
	case snettypes.ProtocolSolofs:
		return p.doPReadMemBlockWithSrpc(peer.ID,
			uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, offset, length)
	}

	return 0, solofsapitypes.ErrServiceNotExists
}

func (p *NetBlockDriver) doPReadMemBlockWithSrpc(peerID snettypes.PeerID,
	uNetINode solofsapitypes.NetINodeUintptr,
	uNetBlock solofsapitypes.NetBlockUintptr,
	netBlockIndex int32,
	uMemBlock solofsapitypes.MemBlockUintptr,
	memBlockIndex int32,
	offset uint64, length int,
) (int, error) {
	var (
		snetReq  snettypes.SNetReq
		snetResp snettypes.SNetResp
		err      error
	)

	var req solofsprotocol.NetINodePReadReq
	req.NetINodeID = uNetBlock.Ptr().NetINodeID
	req.Offset = offset
	req.Length = int32(length)

	// TODO choose solodn
	err = p.SNetClientDriver.Call(peerID, "/NetINode/PRead", &snetReq, &snetResp, req)
	if err != nil {
		return 0, err
	}

	var (
		respParamBs                 = make([]byte, snetResp.ParamSize)
		resp                        solofsprotocol.NetINodePReadResp
		offsetInMemBlock, readedLen int
	)
	err = p.SNetClientDriver.ReadResponse(peerID, &snetReq, &snetResp, respParamBs, &resp)
	if err != nil {
		return 0, err
	}

	offsetInMemBlock = int(offset - uint64(uMemBlock.Ptr().Bytes.Cap)*uint64(memBlockIndex))
	readedLen = int(snetResp.ConnBytesLeft)
	err = p.SNetClientDriver.ReadRawResponse(peerID, &snetReq, &snetResp,
		(*uMemBlock.Ptr().BytesSlice())[offsetInMemBlock:readedLen])
	if err != nil {
		return 0, err
	}

	return int(resp.Length), err
}
