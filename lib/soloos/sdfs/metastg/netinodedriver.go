package metastg

import (
	"soloos/common/sdbapi"
	sdbapitypes "soloos/common/sdbapi/types"
	sdfsapitypes "soloos/common/sdfsapi/types"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/sdfs/types"
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

func (p *NetINodeDriver) PrepareNetINodeMetaDataOnlyLoadDB(uNetINode types.NetINodeUintptr) error {
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

func (p *NetINodeDriver) PrepareNetINodeMetaDataWithStorDB(uNetINode types.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int) error {
	var (
		pNetINode = uNetINode.Ptr()
		err       error
	)

	err = p.FetchNetINodeFromDB(pNetINode)
	if err != nil {
		if err != types.ErrObjectNotExists {
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

func (p *NetINodeDriver) NetINodeTruncate(uNetINode types.NetINodeUintptr, size uint64) error {
	return p.NetINodeCommitSizeInDB(uNetINode, size)
}
