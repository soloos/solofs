package metastg

import (
	"soloos/sdfs/api"
	"soloos/util/offheap"

	"github.com/gocraft/dbr"
)

type MetaStg struct {
	offheapDriver *offheap.OffheapDriver
	dbConn        *dbr.Connection
	DataNodeDriver
	NetINodeDriver
	NetBlockDriver
	DirTreeDriver
}

func (p *MetaStg) Init(offheapDriver *offheap.OffheapDriver,
	dbDriver, dsn string,
	mustGetNetINodeForDirTreeDriver api.MustGetNetINode) error {
	var err error

	p.offheapDriver = offheapDriver
	p.dbConn, err = dbr.Open(dbDriver, dsn, nil)
	if err != nil {
		return err
	}

	switch dbDriver {
	case "mysql":
		err = p.InstallMysqlSchema()
	case "sqlite3":
		err = p.InstallSqlite3Schema()
	}

	err = p.DataNodeDriver.Init(p)
	if err != nil {
		return err
	}

	err = p.NetINodeDriver.Init(p.offheapDriver,
		p.dbConn,
		p.DataNodeDriver.ChooseOneDataNode)
	if err != nil {
		return err
	}

	err = p.NetBlockDriver.Init(p.offheapDriver,
		p.dbConn,
		p.DataNodeDriver.GetDataNode,
		p.NetINodeDriver.ChooseDataNodesForNewNetBlock)
	if err != nil {
		return err
	}

	err = p.DirTreeDriver.Init(p.offheapDriver,
		p.dbConn,
		p.FetchAndUpdateMaxID,
		mustGetNetINodeForDirTreeDriver,
	)

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
