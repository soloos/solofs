package metastg

import "C"
import (
	"soloos/common/sdbapi"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/sdfs/api"
)

type DirTreeStg struct {
	*soloosbase.SoloOSEnv
	dbConn        *sdbapi.Connection
	FsINodeDriver FsINodeDriver
	FIXAttrDriver FIXAttrDriver
}

func (p *DirTreeStg) Init(soloOSEnv *soloosbase.SoloOSEnv,
	dbConn *sdbapi.Connection,
	getNetINodeWithReadAcquire api.GetNetINodeWithReadAcquire,
	mustGetNetINodeWithReadAcquire api.MustGetNetINodeWithReadAcquire,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
	p.dbConn = dbConn

	err = p.installSchema()
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.Init(p.SoloOSEnv, dbConn,
		getNetINodeWithReadAcquire,
		mustGetNetINodeWithReadAcquire)
	if err != nil {
		return err
	}

	err = p.FIXAttrDriver.Init(p.SoloOSEnv, dbConn)
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
