package metastg

import (
	"soloos/dbcli"
	"soloos/util/offheap"
)

// FIXAttrDriver is FsINode XAttr driver
type FIXAttrDriver struct {
	DBConn        *dbcli.Connection
	OffheapDriver *offheap.OffheapDriver
}

func (p *FIXAttrDriver) Init(offheapDriver *offheap.OffheapDriver,
	dbConn *dbcli.Connection,
) error {
	p.OffheapDriver = offheapDriver
	p.DBConn = dbConn
	return nil
}
