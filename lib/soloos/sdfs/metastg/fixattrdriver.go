package metastg

import (
	"soloos/common/sdbapi"
	"soloos/common/util/offheap"
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
