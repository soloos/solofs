package metastg

import "C"
import (
	"soloos/sdfs/api"
	"soloos/util/offheap"

	"github.com/gocraft/dbr"
)

type DirTreeDriverHelper struct {
	DBConn              *dbr.Connection
	OffheapDriver       *offheap.OffheapDriver
	FetchAndUpdateMaxID api.FetchAndUpdateMaxID
	MustGetNetINode     api.MustGetNetINode
}

type DirTreeDriver struct {
	allocINodeIDDalta int64
	lastFsINodeIDInDB int64
	maxFsINodeID      int64
	helper            DirTreeDriverHelper
}

func (p *DirTreeDriver) Init(offheapDriver *offheap.OffheapDriver,
	dbConn *dbr.Connection,
	fetchAndUpdateMaxID api.FetchAndUpdateMaxID,
	mustGetNetINode api.MustGetNetINode,
) error {
	var err error
	p.SetHelper(offheapDriver,
		dbConn,
		fetchAndUpdateMaxID,
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
	mustGetNetINode api.MustGetNetINode,
) {
	p.helper.OffheapDriver = offheapDriver
	p.helper.DBConn = dbConn
	p.helper.FetchAndUpdateMaxID = fetchAndUpdateMaxID
	p.helper.MustGetNetINode = mustGetNetINode
}
