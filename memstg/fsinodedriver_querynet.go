package memstg

import (
	"soloos/common/snet"
	"soloos/common/solofsapitypes"
)

func (p *FsINodeDriver) AllocFsINodeID(
	nsID solofsapitypes.NameSpaceID) (solofsapitypes.FsINodeID, error) {
	var ret = snet.Response{RespData: solofsapitypes.FsINodeID(0)}
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/AllocFsINodeID", &ret, nsID)
	if err != nil {
		return 0, err
	}
	return ret.RespData.(solofsapitypes.FsINodeID), err
}

func (p *FsINodeDriver) DeleteFsINodeByIDInDB(
	nsID solofsapitypes.NameSpaceID,
	fsINodeID solofsapitypes.FsINodeID) error {
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/DeleteFsINodeByIDInDB", nil,
		nsID, fsINodeID)
	if err != nil {
		return err
	}
	return nil
}

func (p *FsINodeDriver) UpdateFsINodeInDB(
	nsID solofsapitypes.NameSpaceID,
	fsINodeMeta solofsapitypes.FsINodeMeta) error {
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/UpdateFsINodeInDB", nil,
		nsID, fsINodeMeta)
	if err != nil {
		return err
	}
	return nil
}

func (p *FsINodeDriver) InsertFsINodeInDB(
	nsID solofsapitypes.NameSpaceID,
	fsINodeMeta solofsapitypes.FsINodeMeta) error {
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/InsertFsINodeInDB", nil,
		nsID, fsINodeMeta)
	if err != nil {
		return err
	}
	return nil
}

func (p *FsINodeDriver) FetchFsINodeByIDFromDB(
	nsID solofsapitypes.NameSpaceID,
	fsINodeID solofsapitypes.FsINodeID) (solofsapitypes.FsINodeMeta, error) {
	var ret = snet.Response{RespData: solofsapitypes.FsINodeMeta{}}
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/FetchFsINodeByIDFromDB", &ret,
		nsID, fsINodeID)
	if err != nil {
		return solofsapitypes.FsINodeMeta{}, err
	}
	return ret.RespData.(solofsapitypes.FsINodeMeta), nil
}

func (p *FsINodeDriver) FetchFsINodeByNameFromDB(
	nsID solofsapitypes.NameSpaceID,
	parentID solofsapitypes.FsINodeID,
	fsINodeName string) (solofsapitypes.FsINodeMeta, error) {
	var ret = snet.Response{RespData: solofsapitypes.FsINodeMeta{}}
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/FetchFsINodeByNameFromDB", &ret,
		nsID, parentID, fsINodeName)
	if err != nil {
		return solofsapitypes.FsINodeMeta{}, err
	}
	return ret.RespData.(solofsapitypes.FsINodeMeta), nil
}

func (p *FsINodeDriver) ListFsINodeByParentIDFromDB(nsID solofsapitypes.NameSpaceID,
	parentID solofsapitypes.FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int64) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(solofsapitypes.FsINodeMeta) bool,
) error {
	var (
		fetchRowsLimit  uint64
		fetchRowsOffset uint64
		err             error
	)

	var retRowsCount = snet.Response{RespData: int64(0)}
	err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/ListFsINodeByParentIDSelectCountFromDB",
		&retRowsCount, nsID, parentID)
	if err != nil {
		return err
	}

	fetchRowsLimit, fetchRowsOffset = beforeLiteralFunc(retRowsCount.RespData.(int64))
	if fetchRowsLimit == 0 {
		return nil
	}

	var iretRows = snet.Response{RespData: []solofsapitypes.FsINodeMeta{}}
	err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/ListFsINodeByParentIDSelectDataFromDB",
		&iretRows, nsID, parentID, fetchRowsLimit, fetchRowsOffset, isFetchAllCols)
	if err != nil {
		return err
	}
	var retRows = iretRows.RespData.([]solofsapitypes.FsINodeMeta)
	for i, _ := range retRows {
		if literalFunc(retRows[i]) == false {
			break
		}
	}

	return nil
}
