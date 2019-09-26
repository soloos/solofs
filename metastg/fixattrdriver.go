package metastg

import (
	"soloos/common/solodbapi"
	"soloos/common/soloosbase"
)

// FIXAttrDriver is FsINode XAttr driver
type FIXAttrDriver struct {
	*soloosbase.SoloOSEnv
	dbConn *solodbapi.Connection
}

func (p *FIXAttrDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	dbConn *solodbapi.Connection,
) error {
	p.SoloOSEnv = soloOSEnv
	p.dbConn = dbConn
	return nil
}
