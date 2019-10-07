package memstg

import (
	"soloos/common/solodbtypes"
	"soloos/common/solofstypes"
)

func (p *NetBlockDriver) prepareNetBlockMetaDataWithTransfer(uNetBlock solofstypes.NetBlockUintptr,
	uNetINode solofstypes.NetINodeUintptr, netblockIndex int32) error {
	var (
		pNetBlock = uNetBlock.Ptr()
	)

	pNetBlock.SyncDataBackends.Append(pNetBlock.StorDataBackends.Arr[0], pNetBlock.StorDataBackends.Len-1)
	for i := 1; i < pNetBlock.StorDataBackends.Len; i++ {
		pNetBlock.SyncDataBackends.Append(pNetBlock.StorDataBackends.Arr[i], 0)
	}

	return nil
}

func (p *NetBlockDriver) prepareNetBlockMetaDataWithFanout(uNetBlock solofstypes.NetBlockUintptr,
	uNetINode solofstypes.NetINodeUintptr, netblockIndex int32) error {
	var (
		pNetBlock = uNetBlock.Ptr()
	)

	for i := 0; i < pNetBlock.StorDataBackends.Len; i++ {
		pNetBlock.SyncDataBackends.Append(pNetBlock.StorDataBackends.Arr[i], 0)
	}

	return nil
}

func (p *NetBlockDriver) PrepareNetBlockMetaData(uNetBlock solofstypes.NetBlockUintptr,
	uNetINode solofstypes.NetINodeUintptr, netblockIndex int32) error {
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
	case solofstypes.BlockPlacementPolicyDefault:
		err = p.prepareNetBlockMetaDataWithFanout(uNetBlock, uNetINode, netblockIndex)

	case solofstypes.BlockPlacementPolicySolomq:
		err = p.solomqClient.PrepareNetBlockMetaData(uNetBlock, uNetINode, netblockIndex)
	}

	if err != nil {
		return err
	}

	pNetBlock.IsDBMetaDataInited.Store(solodbtypes.MetaDataStateInited)
	return nil
}
