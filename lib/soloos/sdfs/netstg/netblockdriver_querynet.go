package netstg

import (
	sdbapitypes "soloos/common/sdbapi/types"
	sdfsapitypes "soloos/common/sdfsapi/types"
	snettypes "soloos/common/snet/types"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
)

// TODO make this configurable
func (p *NetBlockDriver) doPrepareNetBlockMetaData(uNetBlock types.NetBlockUintptr,
	uNetINode types.NetINodeUintptr, netblockIndex int32,
	syncDataPrimaryBackendTransferCount int,
) error {
	var (
		pNetBlock    = uNetBlock.Ptr()
		netBlockInfo protocol.NetINodeNetBlockInfoResponse
		backend      protocol.SNetPeer
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
	pNetBlock.SyncDataPrimaryBackendTransferCount = syncDataPrimaryBackendTransferCount
	return nil
}

func (p *NetBlockDriver) PrepareNetBlockMetaDataWithTransfer(uNetBlock types.NetBlockUintptr,
	uNetINode types.NetINodeUintptr, netblockIndex int32) error {
	var err error
	err = p.doPrepareNetBlockMetaData(uNetBlock, uNetINode, netblockIndex,
		uNetBlock.Ptr().SyncDataBackends.Len-1)
	if err != nil {
		return err
	}
	uNetBlock.Ptr().IsDBMetaDataInited.Store(sdbapitypes.MetaDataStateInited)
	return nil
}

func (p *NetBlockDriver) PrepareNetBlockMetaDataWithFanout(uNetBlock types.NetBlockUintptr,
	uNetINode types.NetINodeUintptr, netblockIndex int32) error {
	var err error
	err = p.doPrepareNetBlockMetaData(uNetBlock, uNetINode, netblockIndex, 0)
	if err != nil {
		return err
	}
	uNetBlock.Ptr().IsDBMetaDataInited.Store(sdbapitypes.MetaDataStateInited)
	return nil
}
