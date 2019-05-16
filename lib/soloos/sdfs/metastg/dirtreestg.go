package metastg

import "C"
import (
	"soloos/common/sdbapi"
	sdfsapitypes "soloos/common/sdfsapi/types"
	soloosbase "soloos/common/soloosapi/base"
)

type DirTreeStg struct {
	*soloosbase.SoloOSEnv
	dbConn        *sdbapi.Connection
	FsINodeDriver FsINodeDriver
	FIXAttrDriver FIXAttrDriver
}

func (p *DirTreeStg) Init(soloOSEnv *soloosbase.SoloOSEnv,
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

func (p *DirTreeStg) Close() error {
	var err error
	err = p.dbConn.Close()
	if err != nil {
		return err
	}

	return nil
}
