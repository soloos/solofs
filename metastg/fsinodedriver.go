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

	allocINodeIDDalta solofstypes.FsINodeIno
	lastFsINodeInoInDB solofstypes.FsINodeIno
	maxFsINodeIno      solofstypes.FsINodeIno
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
	for p.lastFsINodeInoInDB <= solofstypes.RootFsINodeIno {
		p.lastFsINodeInoInDB, err = FetchAndUpdateMaxID(p.dbConn, "b_fsinode", p.allocINodeIDDalta)
		if err != nil {
			return err
		}
		p.maxFsINodeIno = p.lastFsINodeInoInDB
	}

	return nil
}

func (p *FsINodeDriver) AllocFsINodeIno() solofstypes.FsINodeIno {
	var ret = atomic.AddUint64(&p.maxFsINodeIno, 1)
	if p.lastFsINodeInoInDB-ret < p.allocINodeIDDalta/100 {
		util.AssertErrIsNil1(FetchAndUpdateMaxID(p.dbConn, "b_fsinode", p.allocINodeIDDalta))
		p.lastFsINodeInoInDB += p.allocINodeIDDalta
	}
	return ret
}
