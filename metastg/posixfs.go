package metastg

import (
	"soloos/common/solodbapi"
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
)

type PosixFS struct {
	nameSpaceID solofsapitypes.NameSpaceID
	*soloosbase.SoloOSEnv
	dbConn        *solodbapi.Connection
	FsINodeDriver FsINodeDriver
	FIXAttrDriver FIXAttrDriver
}

func (p *PosixFS) Init(soloOSEnv *soloosbase.SoloOSEnv,
	dbConn *solodbapi.Connection,
	getNetINode solofsapitypes.GetNetINode,
	mustGetNetINode solofsapitypes.MustGetNetINode,
	releaseNetINode solofsapitypes.ReleaseNetINode,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
	p.dbConn = dbConn

	err = p.installSchema()
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.Init(p.SoloOSEnv, dbConn,
		getNetINode,
		mustGetNetINode,
		releaseNetINode,
	)
	if err != nil {
		return err
	}

	err = p.FIXAttrDriver.Init(p.SoloOSEnv, dbConn)
	if err != nil {
		return err
	}

	return nil
}

func (p *PosixFS) Close() error {
	var err error
	err = p.dbConn.Close()
	if err != nil {
		return err
	}

	return nil
}
