package metastg

import (
	"soloos/common/sdbapi"
	"soloos/common/soloosbase"
)

// FIXAttrDriver is FsINode XAttr driver
type FIXAttrDriver struct {
	*soloosbase.SoloOSEnv
	dbConn *sdbapi.Connection
}

func (p *FIXAttrDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	dbConn *sdbapi.Connection,
) error {
	p.SoloOSEnv = soloOSEnv
	p.dbConn = dbConn
	return nil
}
