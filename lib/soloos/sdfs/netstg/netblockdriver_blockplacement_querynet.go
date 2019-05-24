package netstg

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
	"soloos/common/sdfsprotocol"
)

// TODO make this configurable
func (p *NetBlockDriver) doPrepareNetBlockMetaData(uNetBlock sdfsapitypes.NetBlockUintptr,
	uNetINode sdfsapitypes.NetINodeUintptr, netblockIndex int32,
) error {
	var (
		pNetBlock    = uNetBlock.Ptr()
		netBlockInfo sdfsprotocol.NetINodeNetBlockInfoResponse
		backend      sdfsprotocol.SNetPeer
		peerID       snettypes.PeerID
		uPeer        snettypes.PeerUintptr
		i            int
		err          error
	)

	err = p.helper.NameNodeClient.PrepareNetBlockMetaData(&netBlockInfo, uNetINode, netblockIndex, uNetBlock)
	if err != nil {
		return err
	}

	pNetBlock.StorDataBackends.Reset()
	for i = 0; i < netBlockInfo.BackendsLength(); i++ {
		netBlockInfo.Backends(&backend, i)
		copy(peerID[:], backend.PeerID())
		uPeer, _ = p.SNetDriver.MustGetPeer(&peerID, string(backend.Address()),
			sdfsapitypes.DefaultSDFSRPCProtocol)
		pNetBlock.StorDataBackends.Append(uPeer)
	}

	pNetBlock.SyncDataBackends = pNetBlock.StorDataBackends
	return nil
}
