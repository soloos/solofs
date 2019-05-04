package metastg

import (
	"soloos/common/sdbapi"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
)

type NetINodeDriverHelper struct {
	ChooseDataNodesForNewNetBlock api.ChooseDataNodesForNewNetBlock
}

type NetINodeDriver struct {
	*soloosbase.SoloOSEnv
	dbConn *sdbapi.Connection
	helper NetINodeDriverHelper
}

func (p *NetINodeDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	dbConn *sdbapi.Connection,
	chooseOneDataNode api.ChooseDataNodesForNewNetBlock,
) error {
	p.dbConn = dbConn
	p.SetHelper(chooseOneDataNode)
	return nil
}

func (p *NetINodeDriver) SetHelper(
	chooseOneDataNode api.ChooseDataNodesForNewNetBlock,
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
		pNetINode.IsDBMetaDataInited.Store(types.MetaDataStateInited)
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
		pNetINode.IsDBMetaDataInited.Store(types.MetaDataStateInited)
	}
	return err
}

func (p *NetINodeDriver) NetINodeTruncate(uNetINode types.NetINodeUintptr, size uint64) error {
	return p.NetINodeCommitSizeInDB(uNetINode, size)
}
