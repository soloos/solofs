package metastg

import (
	"soloos/common/solodbapi"
	"soloos/common/soloosbase"
)

// FIXAttrDriver is FsINode XAttr driver
type FIXAttrDriver struct {
	*soloosbase.SoloosEnv
	dbConn *solodbapi.Connection
}

func (p *FIXAttrDriver) Init(soloosEnv *soloosbase.SoloosEnv,
	dbConn *solodbapi.Connection,
) error {
	p.SoloosEnv = soloosEnv
	p.dbConn = dbConn
	return nil
}
