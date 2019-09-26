package memstg

import (
	"soloos/common/solodbapitypes"
	"soloos/common/solofsapitypes"
)

func (p *NetINodeDriver) prepareNetINodeMetaDataCommon(pNetINode *solofsapitypes.NetINode) {
	pNetINode.MemBlockPlacementPolicy.SetType(solofsapitypes.BlockPlacementPolicyDefault)
	pNetINode.IsDBMetaDataInited.Store(solodbapitypes.MetaDataStateInited)
}

func (p *NetINodeDriver) PrepareNetINodeMetaDataOnlyLoadDB(uNetINode solofsapitypes.NetINodeUintptr) error {
	var err error

	err = p.helper.SolonnClient.GetNetINodeMetaData(uNetINode)
	if err != nil {
		return err
	}

	p.prepareNetINodeMetaDataCommon(uNetINode.Ptr())

	return nil
}

func (p *NetINodeDriver) PrepareNetINodeMetaDataWithStorDB(uNetINode solofsapitypes.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int) error {
	var err error

	err = p.helper.SolonnClient.MustGetNetINodeMetaData(uNetINode, size, netBlockCap, memBlockCap)
	if err != nil {
		return err
	}

	p.prepareNetINodeMetaDataCommon(uNetINode.Ptr())

	return nil
}

func (p *NetINodeDriver) NetINodeCommitSizeInDB(uNetINode solofsapitypes.NetINodeUintptr, size uint64) error {
	var err error
	err = p.helper.SolonnClient.NetINodeCommitSizeInDB(uNetINode, size)
	if err != nil {
		return err
	}

	uNetINode.Ptr().Size = size
	return nil
}
