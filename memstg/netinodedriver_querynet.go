package memstg

import (
	"soloos/common/sdbapitypes"
	"soloos/common/sdfsapitypes"
)

func (p *NetINodeDriver) prepareNetINodeMetaDataCommon(pNetINode *sdfsapitypes.NetINode) {
	pNetINode.MemBlockPlacementPolicy.SetType(sdfsapitypes.BlockPlacementPolicyDefault)
	pNetINode.IsDBMetaDataInited.Store(sdbapitypes.MetaDataStateInited)
}

func (p *NetINodeDriver) PrepareNetINodeMetaDataOnlyLoadDB(uNetINode sdfsapitypes.NetINodeUintptr) error {
	var err error

	err = p.helper.NameNodeClient.GetNetINodeMetaData(uNetINode)
	if err != nil {
		return err
	}

	p.prepareNetINodeMetaDataCommon(uNetINode.Ptr())

	return nil
}

func (p *NetINodeDriver) PrepareNetINodeMetaDataWithStorDB(uNetINode sdfsapitypes.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int) error {
	var err error

	err = p.helper.NameNodeClient.MustGetNetINodeMetaData(uNetINode, size, netBlockCap, memBlockCap)
	if err != nil {
		return err
	}

	p.prepareNetINodeMetaDataCommon(uNetINode.Ptr())

	return nil
}

func (p *NetINodeDriver) NetINodeCommitSizeInDB(uNetINode sdfsapitypes.NetINodeUintptr, size uint64) error {
	var err error
	err = p.helper.NameNodeClient.NetINodeCommitSizeInDB(uNetINode, size)
	if err != nil {
		return err
	}

	uNetINode.Ptr().Size = size
	return nil
}
