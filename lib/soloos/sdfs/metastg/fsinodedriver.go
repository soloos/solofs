package metastg

import (
	"soloos/common/sdbapi"
	soloosbase "soloos/common/soloosapi/base"
	"soloos/common/util"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"sync/atomic"
)

type FsINodeDriverHelper struct {
	GetNetINodeWithReadAcquire     api.GetNetINodeWithReadAcquire
	MustGetNetINodeWithReadAcquire api.MustGetNetINodeWithReadAcquire
}

type FsINodeDriver struct {
	*soloosbase.SoloOSEnv
	dbConn *sdbapi.Connection
	helper FsINodeDriverHelper

	allocINodeIDDalta types.FsINodeID
	lastFsINodeIDInDB types.FsINodeID
	maxFsINodeID      types.FsINodeID
}

func (p *FsINodeDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	dbConn *sdbapi.Connection,
	getNetINodeWithReadAcquire api.GetNetINodeWithReadAcquire,
	mustGetNetINodeWithReadAcquire api.MustGetNetINodeWithReadAcquire,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
	p.dbConn = dbConn
	p.SetHelper(
		getNetINodeWithReadAcquire,
		mustGetNetINodeWithReadAcquire,
	)

	err = p.prepareINodes()
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) SetHelper(
	getNetINodeWithReadAcquire api.GetNetINodeWithReadAcquire,
	mustGetNetINodeWithReadAcquire api.MustGetNetINodeWithReadAcquire,
) {
	p.helper.GetNetINodeWithReadAcquire = getNetINodeWithReadAcquire
	p.helper.MustGetNetINodeWithReadAcquire = mustGetNetINodeWithReadAcquire
}

func (p *FsINodeDriver) prepareINodes() error {
	var err error

	p.allocINodeIDDalta = 10000 * 10
	for p.lastFsINodeIDInDB <= types.RootFsINodeID {
		p.lastFsINodeIDInDB, err = FetchAndUpdateMaxID(p.dbConn, "b_fsinode", p.allocINodeIDDalta)
		if err != nil {
			return err
		}
		p.maxFsINodeID = p.lastFsINodeIDInDB
	}

	return nil
}

func (p *FsINodeDriver) AllocFsINodeID() types.FsINodeID {
	var ret = atomic.AddUint64(&p.maxFsINodeID, 1)
	if p.lastFsINodeIDInDB-ret < p.allocINodeIDDalta/100 {
		util.AssertErrIsNil1(FetchAndUpdateMaxID(p.dbConn, "b_fsinode", p.allocINodeIDDalta))
		p.lastFsINodeIDInDB += p.allocINodeIDDalta
	}
	return ret
}
