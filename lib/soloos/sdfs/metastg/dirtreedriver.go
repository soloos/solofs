package metastg

import "C"
import (
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/util/offheap"
	"sync"

	"github.com/gocraft/dbr"
)

type DirTreeDriverHelper struct {
	DBConn              *dbr.Connection
	OffheapDriver       *offheap.OffheapDriver
	FetchAndUpdateMaxID api.FetchAndUpdateMaxID
	GetNetINode         api.GetNetINode
	MustGetNetINode     api.MustGetNetINode
}

type DirTreeDriver struct {
	allocINodeIDDalta int64
	lastFsINodeIDInDB int64
	maxFsINodeID      int64
	helper            DirTreeDriverHelper

	fsINodesRWMutex sync.RWMutex
	fsINodes        map[string]types.FsINode
}

func (p *DirTreeDriver) Init(offheapDriver *offheap.OffheapDriver,
	dbConn *dbr.Connection,
	fetchAndUpdateMaxID api.FetchAndUpdateMaxID,
	getNetINode api.GetNetINode,
	mustGetNetINode api.MustGetNetINode,
) error {
	var err error
	p.SetHelper(offheapDriver,
		dbConn,
		fetchAndUpdateMaxID,
		getNetINode,
		mustGetNetINode,
	)

	err = p.PrepareINodes()
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeDriver) SetHelper(offheapDriver *offheap.OffheapDriver,
	dbConn *dbr.Connection,
	fetchAndUpdateMaxID api.FetchAndUpdateMaxID,
	getNetINode api.GetNetINode,
	mustGetNetINode api.MustGetNetINode,
) {
	p.helper.OffheapDriver = offheapDriver
	p.helper.DBConn = dbConn
	p.helper.FetchAndUpdateMaxID = fetchAndUpdateMaxID
	p.helper.GetNetINode = getNetINode
	p.helper.MustGetNetINode = mustGetNetINode
}
