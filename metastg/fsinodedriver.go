package metastg

import (
	"soloos/common/sdbapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"sync/atomic"
)

type FsINodeDriverHelper struct {
	GetNetINode     sdfsapitypes.GetNetINode
	MustGetNetINode sdfsapitypes.MustGetNetINode
	ReleaseNetINode sdfsapitypes.ReleaseNetINode
}

type FsINodeDriver struct {
	*soloosbase.SoloOSEnv
	dbConn *sdbapi.Connection
	helper FsINodeDriverHelper

	allocINodeIDDalta sdfsapitypes.FsINodeID
	lastFsINodeIDInDB sdfsapitypes.FsINodeID
	maxFsINodeID      sdfsapitypes.FsINodeID
}

func (p *FsINodeDriver) Init(soloOSEnv *soloosbase.SoloOSEnv,
	dbConn *sdbapi.Connection,
	getNetINode sdfsapitypes.GetNetINode,
	mustGetNetINode sdfsapitypes.MustGetNetINode,
	releaseNetINode sdfsapitypes.ReleaseNetINode,
) error {
	var err error

	p.SoloOSEnv = soloOSEnv
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
	getNetINode sdfsapitypes.GetNetINode,
	mustGetNetINode sdfsapitypes.MustGetNetINode,
	ReleaseNetINode sdfsapitypes.ReleaseNetINode,
) {
	p.helper.GetNetINode = getNetINode
	p.helper.MustGetNetINode = mustGetNetINode
	p.helper.ReleaseNetINode = ReleaseNetINode
}

func (p *FsINodeDriver) prepareINodes() error {
	var err error

	p.allocINodeIDDalta = 10000 * 10
	for p.lastFsINodeIDInDB <= sdfsapitypes.RootFsINodeID {
		p.lastFsINodeIDInDB, err = FetchAndUpdateMaxID(p.dbConn, "b_fsinode", p.allocINodeIDDalta)
		if err != nil {
			return err
		}
		p.maxFsINodeID = p.lastFsINodeIDInDB
	}

	return nil
}

func (p *FsINodeDriver) AllocFsINodeID() sdfsapitypes.FsINodeID {
	var ret = atomic.AddUint64(&p.maxFsINodeID, 1)
	if p.lastFsINodeIDInDB-ret < p.allocINodeIDDalta/100 {
		util.AssertErrIsNil1(FetchAndUpdateMaxID(p.dbConn, "b_fsinode", p.allocINodeIDDalta))
		p.lastFsINodeIDInDB += p.allocINodeIDDalta
	}
	return ret
}
