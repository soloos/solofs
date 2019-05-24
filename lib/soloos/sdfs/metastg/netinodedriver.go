package metastg

import (
	"soloos/common/sdbapi"
	"soloos/common/sdbapitypes"
	"soloos/common/sdfsapitypes"
	"soloos/common/soloosbase"
)

type NetINodeDriverHelper struct {
	ChooseDataNodesForNewNetBlock sdfsapitypes.ChooseDataNodesForNewNetBlock
}

type NetINodeDriver struct {
	*soloosbase.SoloOSEnv
	dbConn *sdbapi.Connection
	helper NetINodeDriverHelper
}

func (p *NetINodeDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	dbConn *sdbapi.Connection,
	chooseOneDataNode sdfsapitypes.ChooseDataNodesForNewNetBlock,
) error {
	p.dbConn = dbConn
	p.SetHelper(chooseOneDataNode)
	return nil
}

func (p *NetINodeDriver) SetHelper(
	chooseOneDataNode sdfsapitypes.ChooseDataNodesForNewNetBlock,
) {
	p.helper.ChooseDataNodesForNewNetBlock = chooseOneDataNode
}

func (p *NetINodeDriver) PrepareNetINodeMetaDataOnlyLoadDB(uNetINode sdfsapitypes.NetINodeUintptr) error {
	var (
		pNetINode = uNetINode.Ptr()
		err       error
	)

	err = p.FetchNetINodeFromDB(pNetINode)
	if err != nil {
		goto PREPARE_DONE
	}

PREPARE_DONE:
	if err == nil {
		pNetINode.IsDBMetaDataInited.Store(sdbapitypes.MetaDataStateInited)
	}
	return err
}

func (p *NetINodeDriver) PrepareNetINodeMetaDataWithStorDB(uNetINode sdfsapitypes.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int) error {
	var (
		pNetINode = uNetINode.Ptr()
		err       error
	)

	err = p.FetchNetINodeFromDB(pNetINode)
	if err != nil {
		if err != sdfsapitypes.ErrObjectNotExists {
			goto PREPARE_DONE
		}

		pNetINode.Size = size
		pNetINode.NetBlockCap = netBlockCap
		pNetINode.MemBlockCap = memBlockCap
		err = p.StoreNetINodeInDB(pNetINode)
		if err != nil {
			goto PREPARE_DONE
		}
	}

PREPARE_DONE:
	if err == nil {
		pNetINode.IsDBMetaDataInited.Store(sdbapitypes.MetaDataStateInited)
	}
	return err
}

func (p *NetINodeDriver) NetINodeTruncate(uNetINode sdfsapitypes.NetINodeUintptr, size uint64) error {
	return p.NetINodeCommitSizeInDB(uNetINode, size)
}
