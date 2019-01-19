package metastg

import (
	"soloos/dbcli"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/util/offheap"
)

type NetINodeDriverHelper struct {
	DBConn                        *dbcli.Connection
	ChooseDataNodesForNewNetBlock api.ChooseDataNodesForNewNetBlock
}

type NetINodeDriver struct {
	helper NetINodeDriverHelper
}

func (p *NetINodeDriver) Init(offheapDriver *offheap.OffheapDriver,
	dbConn *dbcli.Connection,
	chooseOneDataNode api.ChooseDataNodesForNewNetBlock,
) error {
	p.SetHelper(dbConn, chooseOneDataNode)
	return nil
}

func (p *NetINodeDriver) SetHelper(dbConn *dbcli.Connection,
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
	return nil
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
	return nil
}
