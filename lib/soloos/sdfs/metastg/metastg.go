package metastg

import (
	"soloos/util/offheap"

	"github.com/gocraft/dbr"
)

type MetaStg struct {
	offheapDriver *offheap.OffheapDriver
	DBConn        *dbr.Connection
	INodeDriver
	NetBlockDriver
}

func (p *MetaStg) Init(offheapDriver *offheap.OffheapDriver, dbDriver, dsn string) error {
	var err error

	p.offheapDriver = offheapDriver
	p.DBConn, err = dbr.Open(dbDriver, dsn, nil)
	if err != nil {
		return err
	}

	switch dbDriver {
	case "mysql":
		err = p.InstallMysqlSchema()
	case "sqlite3":
		err = p.InstallSqlite3Schema()
	}

	err = p.INodeDriver.Init(p)
	if err != nil {
		return err
	}

	err = p.NetBlockDriver.Init(p)
	if err != nil {
		return err
	}

	return nil
}

func (p *MetaStg) Close() error {
	var err error

	err = p.DBConn.Close()
	if err != nil {
		return err
	}

	return nil
}
