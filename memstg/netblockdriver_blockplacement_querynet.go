package memstg

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
)

// TODO make this configurable
func (p *NetBlockDriver) doPrepareNetBlockMetaData(uNetBlock solofsapitypes.NetBlockUintptr,
	uNetINode solofsapitypes.NetINodeUintptr, netblockIndex int32,
) error {
	var (
		pNetBlock    = uNetBlock.Ptr()
		netBlockInfo solofsprotocol.NetINodeNetBlockInfoResponse
		peerID       snettypes.PeerID
		i            int
		err          error
	)

	err = p.helper.SolonnClient.PrepareNetBlockMetaData(&netBlockInfo, uNetINode, netblockIndex, uNetBlock)
	if err != nil {
		return err
	}

	pNetBlock.StorDataBackends.Reset()
	for i = 0; i < netBlockInfo.BackendsLength(); i++ {
		copy(peerID[:], netBlockInfo.Backends(i))
		pNetBlock.StorDataBackends.Append(peerID)
	}

	return nil
}
