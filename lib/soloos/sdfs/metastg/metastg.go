package metastg

import (
	"soloos/dbcli"
	"soloos/util/offheap"
)

type MetaStg struct {
	offheapDriver *offheap.OffheapDriver
	dbConn        dbcli.Connection
	DataNodeDriver
	NetINodeDriver
	NetBlockDriver
}

func (p *MetaStg) Init(offheapDriver *offheap.OffheapDriver,
	dbDriver, dsn string,
) error {
	var err error

	p.offheapDriver = offheapDriver
	err = p.dbConn.Init(dbDriver, dsn)
	if err != nil {
		return err
	}

	switch dbDriver {
	case "mysql":
		err = InstallMysqlSchema(&p.dbConn)
	case "sqlite3":
		err = InstallSqlite3Schema(&p.dbConn)
	}

	err = p.DataNodeDriver.Init(p)
	if err != nil {
		return err
	}

	err = p.NetINodeDriver.Init(p.offheapDriver,
		&p.dbConn,
		p.DataNodeDriver.ChooseDataNodesForNewNetBlock)
	if err != nil {
		return err
	}

	err = p.NetBlockDriver.Init(p.offheapDriver,
		&p.dbConn,
		p.DataNodeDriver.GetDataNode,
		p.NetINodeDriver.ChooseDataNodesForNewNetBlock)
	if err != nil {
		return err
	}

	return nil
}

func (p *MetaStg) Close() error {
	var err error

	err = p.dbConn.Close()
	if err != nil {
		return err
	}

	return nil
}
