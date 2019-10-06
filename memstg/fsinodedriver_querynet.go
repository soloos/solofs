package memstg

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
)

func (p *FsINodeDriver) AllocFsINodeID(
	nsID solofsapitypes.NameSpaceID) (solofsapitypes.FsINodeID, error) {
	var ret = snettypes.Response{RespData: solofsapitypes.FsINodeID(0)}
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

func (p *FsINodeDriver) ListFsINodeByParentIDFromDB(
	nsID solofsapitypes.NameSpaceID,
	parentID solofsapitypes.FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(solofsapitypes.FsINodeMeta) bool,
) error {
	panic("shit")
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
	var ret = snettypes.Response{RespData: solofsapitypes.FsINodeMeta{}}
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
	var ret = snettypes.Response{RespData: solofsapitypes.FsINodeMeta{}}
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FsINode/FetchFsINodeByNameFromDB", &ret,
		nsID, parentID, fsINodeName)
	if err != nil {
		return solofsapitypes.FsINodeMeta{}, err
	}
	return ret.RespData.(solofsapitypes.FsINodeMeta), nil
}
