package solodn

import (
	"soloos/common/snettypes"
	"soloos/common/solodbapitypes"
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
)

func (p *SrpcServer) NetINodePWrite(reqCtx *snettypes.SNetReqContext,
	req solofsprotocol.NetINodePWriteReq,
) error {
	var (
		syncDataBackends snettypes.PeerGroup
		peerID           snettypes.PeerID
		uNetBlock        solofsapitypes.NetBlockUintptr
		err              error
	)

	// get uNetINode
	var (
		netINodeID         solofsapitypes.NetINodeID
		uNetINode          solofsapitypes.NetINodeUintptr
		firstNetBlockIndex int32
		lastNetBlockIndex  int32
		netBlockIndex      int32
	)
	netINodeID = req.NetINodeID
	uNetINode, err = p.solodn.netINodeDriver.GetNetINode(netINodeID)
	defer p.solodn.netINodeDriver.ReleaseNetINode(uNetINode)
	if err != nil {
		return err
	}

	// TODO no need prepare syncDataBackends every pwrite
	syncDataBackends.Reset()
	syncDataBackends.Append(p.solodn.localFsSNetPeer.ID)
	for i, _ := range req.TransferBackends {
		copy(peerID[:], req.TransferBackends[i])
		syncDataBackends.Append(peerID)
	}

	// prepare uNetBlock
	firstNetBlockIndex = int32(req.Offset / uint64(uNetINode.Ptr().NetBlockCap))
	lastNetBlockIndex = int32((req.Offset + uint64(req.Length)) / uint64(uNetINode.Ptr().NetBlockCap))
	for netBlockIndex = firstNetBlockIndex; netBlockIndex <= lastNetBlockIndex; netBlockIndex++ {
		uNetBlock, err = p.solodn.netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
		defer p.solodn.netBlockDriver.ReleaseNetBlock(uNetBlock)
		if err != nil {
			return err
		}

		if uNetBlock.Ptr().IsSyncDataBackendsInited.Load() == solodbapitypes.MetaDataStateUninited {
			p.solodn.PrepareNetBlockSyncDataBackends(uNetBlock, syncDataBackends)
		}
	}

	// request file data
	err = p.solodn.netINodeDriver.PWriteWithNetQuery(uNetINode, &reqCtx.NetQuery,
		int(req.Length), req.Offset)
	if err != nil {
		return err
	}

	return nil
}
