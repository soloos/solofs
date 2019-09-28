package metastg

import (
	"soloos/common/solodbapi"
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
)

type PosixFs struct {
	nameSpaceID solofsapitypes.NameSpaceID
	*soloosbase.SoloosEnv
	dbConn        *solodbapi.Connection
	FsINodeDriver FsINodeDriver
	FIXAttrDriver FIXAttrDriver
}

func (p *PosixFs) Init(soloosEnv *soloosbase.SoloosEnv,
	dbConn *solodbapi.Connection,
	getNetINode solofsapitypes.GetNetINode,
	mustGetNetINode solofsapitypes.MustGetNetINode,
	releaseNetINode solofsapitypes.ReleaseNetINode,
) error {
	var err error

	p.SoloosEnv = soloosEnv
	p.dbConn = dbConn

	err = p.installSchema()
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.Init(p.SoloosEnv, dbConn,
		getNetINode,
		mustGetNetINode,
		releaseNetINode,
	)
	if err != nil {
		return err
	}

	err = p.FIXAttrDriver.Init(p.SoloosEnv, dbConn)
	if err != nil {
		return err
	}

	return nil
}

func (p *PosixFs) Close() error {
	var err error
	err = p.dbConn.Close()
	if err != nil {
		return err
	}

	return nil
}
