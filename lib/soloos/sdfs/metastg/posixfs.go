package metastg

import (
	"soloos/common/sdbapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/soloosbase"
)

type PosixFS struct {
	nameSpaceID sdfsapitypes.NameSpaceID
	*soloosbase.SoloOSEnv
	dbConn        *sdbapi.Connection
	FsINodeDriver FsINodeDriver
	FIXAttrDriver FIXAttrDriver
}

func (p *PosixFS) Init(soloOSEnv *soloosbase.SoloOSEnv,
	dbConn *sdbapi.Connection,
	getNetINode sdfsapitypes.GetNetINode,
	mustGetNetINode sdfsapitypes.MustGetNetINode,
	releaseNetINode sdfsapitypes.ReleaseNetINode,
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
