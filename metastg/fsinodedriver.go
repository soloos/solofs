package metastg

import (
	"soloos/common/solodbapi"
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"sync/atomic"
)

type FsINodeDriver struct {
	*soloosbase.SoloosEnv
	dbConn *solodbapi.Connection

	allocINodeIDDalta solofstypes.FsINodeID
	lastFsINodeIDInDB solofstypes.FsINodeID
	maxFsINodeID      solofstypes.FsINodeID
}

func (p *FsINodeDriver) Init(soloosEnv *soloosbase.SoloosEnv,
	dbConn *solodbapi.Connection,
) error {
	var err error

	p.SoloosEnv = soloosEnv
	p.dbConn = dbConn

	err = p.prepareINodes()
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) prepareINodes() error {
	var err error

	p.allocINodeIDDalta = 10000 * 10
	for p.lastFsINodeIDInDB <= solofstypes.RootFsINodeID {
		p.lastFsINodeIDInDB, err = FetchAndUpdateMaxID(p.dbConn, "b_fsinode", p.allocINodeIDDalta)
		if err != nil {
			return err
		}
		p.maxFsINodeID = p.lastFsINodeIDInDB
	}

	return nil
}

func (p *FsINodeDriver) AllocFsINodeID() solofstypes.FsINodeID {
	var ret = atomic.AddUint64(&p.maxFsINodeID, 1)
	if p.lastFsINodeIDInDB-ret < p.allocINodeIDDalta/100 {
		util.AssertErrIsNil1(FetchAndUpdateMaxID(p.dbConn, "b_fsinode", p.allocINodeIDDalta))
		p.lastFsINodeIDInDB += p.allocINodeIDDalta
	}
	return ret
}
