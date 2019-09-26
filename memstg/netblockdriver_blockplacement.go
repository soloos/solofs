package memstg

import (
	"soloos/common/solodbapitypes"
	"soloos/common/solofsapitypes"
)

func (p *NetBlockDriver) prepareNetBlockMetaDataWithTransfer(uNetBlock solofsapitypes.NetBlockUintptr,
	uNetINode solofsapitypes.NetINodeUintptr, netblockIndex int32) error {
	var (
		pNetBlock = uNetBlock.Ptr()
	)

	pNetBlock.SyncDataBackends.Append(pNetBlock.StorDataBackends.Arr[0], pNetBlock.StorDataBackends.Len-1)
	for i := 1; i < pNetBlock.StorDataBackends.Len; i++ {
		pNetBlock.SyncDataBackends.Append(pNetBlock.StorDataBackends.Arr[i], 0)
	}

	return nil
}

func (p *NetBlockDriver) prepareNetBlockMetaDataWithFanout(uNetBlock solofsapitypes.NetBlockUintptr,
	uNetINode solofsapitypes.NetINodeUintptr, netblockIndex int32) error {
	var (
		pNetBlock = uNetBlock.Ptr()
	)

	for i := 0; i < pNetBlock.StorDataBackends.Len; i++ {
		pNetBlock.SyncDataBackends.Append(pNetBlock.StorDataBackends.Arr[i], 0)
	}

	return nil
}

func (p *NetBlockDriver) PrepareNetBlockMetaData(uNetBlock solofsapitypes.NetBlockUintptr,
	uNetINode solofsapitypes.NetINodeUintptr, netblockIndex int32) error {
	var (
		pNetBlock = uNetBlock.Ptr()
		pNetINode = uNetINode.Ptr()
		err       error
	)

	err = p.doPrepareNetBlockMetaData(uNetBlock, uNetINode, netblockIndex)
	if err != nil {
		return err
	}

	switch pNetINode.MemBlockPlacementPolicy.GetType() {
	case solofsapitypes.BlockPlacementPolicyDefault:
		err = p.prepareNetBlockMetaDataWithFanout(uNetBlock, uNetINode, netblockIndex)

	case solofsapitypes.BlockPlacementPolicySOLOMQ:
		err = p.helper.SOLOMQClient.PrepareNetBlockMetaData(uNetBlock, uNetINode, netblockIndex)
	}

	if err != nil {
		return err
	}

	pNetBlock.IsDBMetaDataInited.Store(solodbapitypes.MetaDataStateInited)
	return nil
}
