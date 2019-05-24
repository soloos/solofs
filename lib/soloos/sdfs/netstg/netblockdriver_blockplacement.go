package netstg

import (
	"soloos/common/sdbapitypes"
	"soloos/common/sdfsapitypes"
)

func (p *NetBlockDriver) prepareNetBlockMetaDataWithTransfer(uNetBlock sdfsapitypes.NetBlockUintptr,
	uNetINode sdfsapitypes.NetINodeUintptr, netblockIndex int32) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)
	err = p.doPrepareNetBlockMetaData(uNetBlock, uNetINode, netblockIndex)
	if err != nil {
		return err
	}
	pNetBlock.SyncDataPrimaryBackendTransferCount = pNetBlock.SyncDataBackends.Len - 1
	return nil
}

func (p *NetBlockDriver) prepareNetBlockMetaDataWithFanout(uNetBlock sdfsapitypes.NetBlockUintptr,
	uNetINode sdfsapitypes.NetINodeUintptr, netblockIndex int32) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		err       error
	)
	err = p.doPrepareNetBlockMetaData(uNetBlock, uNetINode, netblockIndex)
	if err != nil {
		return err
	}
	pNetBlock.SyncDataPrimaryBackendTransferCount = 0
	return nil
}

func (p *NetBlockDriver) PrepareNetBlockMetaData(uNetBlock sdfsapitypes.NetBlockUintptr,
	uNetINode sdfsapitypes.NetINodeUintptr, netblockIndex int32) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		pNetINode = uNetINode.Ptr()
		err       error
	)

	err = p.prepareNetBlockMetaDataWithFanout(uNetBlock, uNetINode, netblockIndex)
	if err != nil {
		return err
	}

	switch pNetINode.MemBlockPlacementPolicy.GetType() {
	case sdfsapitypes.BlockPlacementPolicyDefault:

	case sdfsapitypes.BlockPlacementPolicySWAL:
		err = p.helper.SWALClient.PrepareNetBlockMetaData(uNetBlock, uNetINode, netblockIndex)
	}

	if err != nil {
		return err
	}

	pNetBlock.IsDBMetaDataInited.Store(sdbapitypes.MetaDataStateInited)
	return nil
}
