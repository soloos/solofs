package metastg

import (
	"soloos/util"
	"sync/atomic"
)

func (p *DirTreeDriver) PrepareINodes() error {
	var err error
	p.allocINodeIDDalta = 10000 * 10
	p.lastFsINodeIDInDB, err = p.helper.FetchAndUpdateMaxID("b_fsinode", p.allocINodeIDDalta)
	if err != nil {
		return err
	}
	p.maxFsINodeID = p.lastFsINodeIDInDB

	return nil
}

func (p *DirTreeDriver) AllocFsINodeID() int64 {
	var ret = atomic.AddInt64(&p.maxFsINodeID, 1)
	if ret-p.lastFsINodeIDInDB > p.allocINodeIDDalta/100 {
		util.AssertErrIsNil1(p.helper.FetchAndUpdateMaxID("b_fsinode", p.allocINodeIDDalta))
	}
	return ret
}
