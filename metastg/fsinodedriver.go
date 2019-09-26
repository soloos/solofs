package metastg

import (
	"soloos/common/solodbapi"
	"soloos/common/solofsapitypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"sync/atomic"
)

type FsINodeDriverHelper struct {
	GetNetINode     solofsapitypes.GetNetINode
	MustGetNetINode solofsapitypes.MustGetNetINode
	ReleaseNetINode solofsapitypes.ReleaseNetINode
}

type FsINodeDriver struct {
	*soloosbase.SoloosEnv
	dbConn *solodbapi.Connection
	helper FsINodeDriverHelper

	allocINodeIDDalta solofsapitypes.FsINodeID
	lastFsINodeIDInDB solofsapitypes.FsINodeID
	maxFsINodeID      solofsapitypes.FsINodeID
}

func (p *FsINodeDriver) Init(soloosEnv *soloosbase.SoloosEnv,
	dbConn *solodbapi.Connection,
	getNetINode solofsapitypes.GetNetINode,
	mustGetNetINode solofsapitypes.MustGetNetINode,
	releaseNetINode solofsapitypes.ReleaseNetINode,
) error {
	var err error

	p.SoloosEnv = soloosEnv
	p.dbConn = dbConn
	p.SetHelper(
		getNetINode,
		mustGetNetINode,
		releaseNetINode,
	)

	err = p.prepareINodes()
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) SetHelper(
	getNetINode solofsapitypes.GetNetINode,
	mustGetNetINode solofsapitypes.MustGetNetINode,
	ReleaseNetINode solofsapitypes.ReleaseNetINode,
) {
	p.helper.GetNetINode = getNetINode
	p.helper.MustGetNetINode = mustGetNetINode
	p.helper.ReleaseNetINode = ReleaseNetINode
}

func (p *FsINodeDriver) prepareINodes() error {
	var err error

	p.allocINodeIDDalta = 10000 * 10
	for p.lastFsINodeIDInDB <= solofsapitypes.RootFsINodeID {
		p.lastFsINodeIDInDB, err = FetchAndUpdateMaxID(p.dbConn, "b_fsinode", p.allocINodeIDDalta)
		if err != nil {
			return err
		}
		p.maxFsINodeID = p.lastFsINodeIDInDB
	}

	return nil
}

func (p *FsINodeDriver) AllocFsINodeID() solofsapitypes.FsINodeID {
	var ret = atomic.AddUint64(&p.maxFsINodeID, 1)
	if p.lastFsINodeIDInDB-ret < p.allocINodeIDDalta/100 {
		util.AssertErrIsNil1(FetchAndUpdateMaxID(p.dbConn, "b_fsinode", p.allocINodeIDDalta))
		p.lastFsINodeIDInDB += p.allocINodeIDDalta
	}
	return ret
}
