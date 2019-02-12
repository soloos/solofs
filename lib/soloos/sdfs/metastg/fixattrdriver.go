package metastg

import (
	"soloos/sdbapi"
	"soloos/util/offheap"
)

// FIXAttrDriver is FsINode XAttr driver
type FIXAttrDriver struct {
	DBConn        *sdbapi.Connection
	OffheapDriver *offheap.OffheapDriver
}

func (p *FIXAttrDriver) Init(offheapDriver *offheap.OffheapDriver,
	dbConn *sdbapi.Connection,
) error {
	p.OffheapDriver = offheapDriver
	p.DBConn = dbConn
	return nil
}
