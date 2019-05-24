package metastg

import (
	"soloos/common/sdbapi"
	"soloos/common/soloosbase"
)

type MetaStg struct {
	*soloosbase.SoloOSEnv
	dbConn sdbapi.Connection

	DataNodeDriver
	NetINodeDriver
	NetBlockDriver
}

func (p *MetaStg) Init(soloOSEnv *soloosbase.SoloOSEnv,
	dbDriver, dsn string,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
	err = p.dbConn.Init(dbDriver, dsn)
	if err != nil {
		return err
	}

	err = p.installSchema()
	if err != nil {
		return err
	}

	err = p.DataNodeDriver.Init(p)
	if err != nil {
		return err
	}

	err = p.NetINodeDriver.Init(p.SoloOSEnv,
		&p.dbConn,
		p.DataNodeDriver.ChooseDataNodesForNewNetBlock)
	if err != nil {
		return err
	}

	err = p.NetBlockDriver.Init(p.SoloOSEnv,
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
