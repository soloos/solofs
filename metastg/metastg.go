package metastg

import (
	"soloos/common/solodbapi"
	"soloos/common/soloosbase"
)

type MetaStg struct {
	*soloosbase.SoloosEnv
	dbConn solodbapi.Connection

	SolodnDriver
	NetINodeDriver
	NetBlockDriver
	FsINodeDriver
	FIXAttrDriver
}

func (p *MetaStg) Init(soloosEnv *soloosbase.SoloosEnv,
	dbDriver, dsn string,
) error {
	var err error

	p.SoloosEnv = soloosEnv
	err = p.dbConn.Init(dbDriver, dsn)
	if err != nil {
		return err
	}

	err = p.installSchema()
	if err != nil {
		return err
	}

	err = p.SolodnDriver.Init(p)
	if err != nil {
		return err
	}

	err = p.NetINodeDriver.Init(p.SoloosEnv,
		&p.dbConn,
		p.SolodnDriver.ChooseSolodnsForNewNetBlock)
	if err != nil {
		return err
	}

	err = p.NetBlockDriver.Init(p.SoloosEnv,
		&p.dbConn,
		p.SolodnDriver.GetSolodn,
		p.NetINodeDriver.ChooseSolodnsForNewNetBlock)
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.Init(p.SoloosEnv, &p.dbConn)
	if err != nil {
		return err
	}

	err = p.FIXAttrDriver.Init(p.SoloosEnv, &p.dbConn)
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
