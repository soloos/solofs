package metastg

import (
	"soloos/sdbapi"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/util/offheap"
)

type NetINodeDriverHelper struct {
	DBConn                        *sdbapi.Connection
	ChooseDataNodesForNewNetBlock api.ChooseDataNodesForNewNetBlock
}

type NetINodeDriver struct {
	helper NetINodeDriverHelper
}

func (p *NetINodeDriver) Init(offheapDriver *offheap.OffheapDriver,
	dbConn *sdbapi.Connection,
	chooseOneDataNode api.ChooseDataNodesForNewNetBlock,
) error {
	p.SetHelper(dbConn, chooseOneDataNode)
	return nil
}

func (p *NetINodeDriver) SetHelper(dbConn *sdbapi.Connection,
	chooseOneDataNode api.ChooseDataNodesForNewNetBlock,
) {
	p.helper.DBConn = dbConn
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
		pNetINode.IsDBMetaDataInited = true
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
		pNetINode.IsDBMetaDataInited = true
	}
	return err
}

func (p *NetINodeDriver) NetINodeTruncate(uNetINode types.NetINodeUintptr, size uint64) error {
	return p.NetINodeCommitSizeInDB(uNetINode, size)
}
