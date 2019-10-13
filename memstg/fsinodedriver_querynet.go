package memstg

import (
	"soloos/common/snet"
	"soloos/common/solofstypes"
)

func (p *FsINodeDriver) AllocFsINodeIno(
	nsID solofstypes.NameSpaceID) (solofstypes.FsINodeIno, error) {
	var ret = snet.Response{RespData: solofstypes.FsINodeIno(0)}
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/AllocFsINodeIno", &ret, nsID)
	if err != nil {
		return 0, err
	}
	return ret.RespData.(solofstypes.FsINodeIno), err
}

func (p *FsINodeDriver) DeleteFsINodeByIDInDB(
	nsID solofstypes.NameSpaceID,
	fsINodeIno solofstypes.FsINodeIno) error {
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/DeleteFsINodeByIDInDB", nil,
		nsID, fsINodeIno)
	if err != nil {
		return err
	}
	return nil
}

func (p *FsINodeDriver) UpdateFsINodeInDB(
	nsID solofstypes.NameSpaceID,
	fsINodeMeta solofstypes.FsINodeMeta) error {
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/UpdateFsINodeInDB", nil,
		nsID, fsINodeMeta)
	if err != nil {
		return err
	}
	return nil
}

func (p *FsINodeDriver) InsertFsINodeInDB(
	nsID solofstypes.NameSpaceID,
	fsINodeMeta solofstypes.FsINodeMeta) error {
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/InsertFsINodeInDB", nil,
		nsID, fsINodeMeta)
	if err != nil {
		return err
	}
	return nil
}

func (p *FsINodeDriver) FetchFsINodeByIDFromDB(
	nsID solofstypes.NameSpaceID,
	fsINodeIno solofstypes.FsINodeIno) (solofstypes.FsINodeMeta, error) {
	var ret = snet.Response{RespData: solofstypes.FsINodeMeta{}}
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/FetchFsINodeByIDFromDB", &ret,
		nsID, fsINodeIno)
	if err != nil {
		return solofstypes.FsINodeMeta{}, err
	}
	return ret.RespData.(solofstypes.FsINodeMeta), nil
}

func (p *FsINodeDriver) FetchFsINodeByNameFromDB(
	nsID solofstypes.NameSpaceID,
	parentID solofstypes.FsINodeIno,
	fsINodeName string) (solofstypes.FsINodeMeta, error) {
	var ret = snet.Response{RespData: solofstypes.FsINodeMeta{}}
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/FetchFsINodeByNameFromDB", &ret,
		nsID, parentID, fsINodeName)
	if err != nil {
		return solofstypes.FsINodeMeta{}, err
	}
	return ret.RespData.(solofstypes.FsINodeMeta), nil
}

func (p *FsINodeDriver) ListFsINodeByParentIDFromDB(nsID solofstypes.NameSpaceID,
	parentID solofstypes.FsINodeIno,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int64) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(solofstypes.FsINodeMeta) bool,
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

	var iretRows = snet.Response{RespData: []solofstypes.FsINodeMeta{}}
	err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/ListFsINodeByParentIDSelectDataFromDB",
		&iretRows, nsID, parentID, fetchRowsLimit, fetchRowsOffset, isFetchAllCols)
	if err != nil {
		return err
	}
	var retRows = iretRows.RespData.([]solofstypes.FsINodeMeta)
	for i, _ := range retRows {
		if literalFunc(retRows[i]) == false {
			break
		}
	}

	return nil
}
