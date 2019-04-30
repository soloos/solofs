package metastg

import (
	"soloos/common/sdbapi"
	"soloos/common/snet"
	"soloos/sdbone/offheap"
)

type MetaStg struct {
	offheapDriver *offheap.OffheapDriver
	dbConn        sdbapi.Connection
	SnetDriver    snet.NetDriver
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

	err = p.installSchema()
	if err != nil {
		return err
	}

	err = p.SnetDriver.Init(p.offheapDriver, "MetaStgNetDriver")
	if err != nil {
		return err
	}

	err = p.DataNodeDriver.Init(p, &p.SnetDriver)
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
