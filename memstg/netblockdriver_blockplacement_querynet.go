package memstg

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
)

// TODO make this configurable
func (p *NetBlockDriver) doPrepareNetBlockMetaData(uNetBlock solofsapitypes.NetBlockUintptr,
	uNetINode solofsapitypes.NetINodeUintptr, netBlockIndex int32,
) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		peerID    snettypes.PeerID
		req       solofsprotocol.NetINodeNetBlockInfoReq
		i         int
		err       error
	)

	req.NetINodeID = uNetINode.Ptr().ID
	req.NetBlockIndex = int32(netBlockIndex)
	req.Cap = int32(uNetINode.Ptr().NetBlockCap)

	// TODO choose solonn
	var netBlockInfo solofsprotocol.NetINodeNetBlockInfoResp
	err = p.solonnClient.Dispatch("/NetBlock/PrepareMetaData", &netBlockInfo, req)
	if err != nil {
		return err
	}

	pNetBlock.NetINodeID = uNetINode.Ptr().ID
	pNetBlock.IndexInNetINode = netBlockIndex
	pNetBlock.Len = int(netBlockInfo.Len)
	pNetBlock.Cap = int(netBlockInfo.Cap)

	pNetBlock.StorDataBackends.Reset()
	for i = 0; i < len(netBlockInfo.Backends); i++ {
		copy(peerID[:], netBlockInfo.Backends[i])
		pNetBlock.StorDataBackends.Append(peerID)
	}

	return nil
}
