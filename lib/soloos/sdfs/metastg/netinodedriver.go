package metastg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/util/offheap"

	"github.com/gocraft/dbr"
)

type NetINodeDriverHelper struct {
	DBConn            *dbr.Connection
	ChooseOneDataNode api.ChooseOneDataNode
}

type NetINodeDriver struct {
	helper       NetINodeDriverHelper
	netINodePool types.NetINodePool
}

func (p *NetINodeDriver) Init(offheapDriver *offheap.OffheapDriver,
	dbConn *dbr.Connection,
	chooseOneDataNode api.ChooseOneDataNode,
) error {
	p.SetHelper(dbConn, chooseOneDataNode)
	p.netINodePool.Init(-1, offheapDriver)
	return nil
}

func (p *NetINodeDriver) SetHelper(dbConn *dbr.Connection,
	chooseOneDataNode api.ChooseOneDataNode,
) {
	p.helper.DBConn = dbConn
	p.helper.ChooseOneDataNode = chooseOneDataNode
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
	size int64, netBlockCap int, memBlockCap int) error {
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
