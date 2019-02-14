package metastg

import "C"
import (
	"soloos/common/sdbapi"
	"soloos/sdfs/api"
	"soloos/sdbone/offheap"
)

type DirTreeStg struct {
	dbConn        *sdbapi.Connection
	FsINodeDriver FsINodeDriver
	FIXAttrDriver FIXAttrDriver
}

func (p *DirTreeStg) Init(offheapDriver *offheap.OffheapDriver,
	dbConn *sdbapi.Connection,
	getNetINodeWithReadAcquire api.GetNetINodeWithReadAcquire,
	mustGetNetINodeWithReadAcquire api.MustGetNetINodeWithReadAcquire,
) error {
	var err error

	p.dbConn = dbConn

	err = p.installSchema()
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.Init(offheapDriver, dbConn,
		getNetINodeWithReadAcquire,
		mustGetNetINodeWithReadAcquire)
	if err != nil {
		return err
	}

	err = p.FIXAttrDriver.Init(offheapDriver, dbConn)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeStg) Close() error {
	var err error
	err = p.dbConn.Close()
	if err != nil {
		return err
	}

	return nil
}
