package netstg

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/sdfsprotocol"
	"soloos/common/snettypes"
)

// TODO make this configurable
func (p *NetBlockDriver) doPrepareNetBlockMetaData(uNetBlock sdfsapitypes.NetBlockUintptr,
	uNetINode sdfsapitypes.NetINodeUintptr, netblockIndex int32,
) error {
	var (
		pNetBlock    = uNetBlock.Ptr()
		netBlockInfo sdfsprotocol.NetINodeNetBlockInfoResponse
		peerID       snettypes.PeerID
		i            int
		err          error
	)

	err = p.helper.NameNodeClient.PrepareNetBlockMetaData(&netBlockInfo, uNetINode, netblockIndex, uNetBlock)
	if err != nil {
		return err
	}

	pNetBlock.StorDataBackends.Reset()
	for i = 0; i < netBlockInfo.BackendsLength(); i++ {
		copy(peerID[:], netBlockInfo.Backends(i))
		pNetBlock.StorDataBackends.Append(peerID)
	}

	pNetBlock.SyncDataBackends = pNetBlock.StorDataBackends
	return nil
}
