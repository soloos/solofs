package memstg

import "soloos/common/solofsapitypes"

func (p *FsINodeDriver) AllocFsINodeID() solofsapitypes.FsINodeID {
	return 0
}

func (p *FsINodeDriver) DeleteFsINodeByIDInDB(nameSpaceID solofsapitypes.NameSpaceID,
	fsINodeID solofsapitypes.FsINodeID) error {
	return nil
}

func (p *FsINodeDriver) ListFsINodeByParentIDFromDB(nameSpaceID solofsapitypes.NameSpaceID,
	parentID solofsapitypes.FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(solofsapitypes.FsINodeMeta) bool,
) error {
	return nil
}

func (p *FsINodeDriver) UpdateFsINodeInDB(nameSpaceID solofsapitypes.NameSpaceID,
	fsINodeMeta solofsapitypes.FsINodeMeta) error {
	return nil
}

func (p *FsINodeDriver) InsertFsINodeInDB(nameSpaceID solofsapitypes.NameSpaceID,
	fsINodeMeta solofsapitypes.FsINodeMeta) error {
	return nil
}

func (p *FsINodeDriver) FetchFsINodeByIDFromDB(nameSpaceID solofsapitypes.NameSpaceID,
	fsINodeID solofsapitypes.FsINodeID) (solofsapitypes.FsINodeMeta, error) {
	var ret solofsapitypes.FsINodeMeta
	return ret, nil
}

func (p *FsINodeDriver) FetchFsINodeByNameFromDB(nameSpaceID solofsapitypes.NameSpaceID,
	parentID solofsapitypes.FsINodeID,
	fsINodeName string) (solofsapitypes.FsINodeMeta, error) {
	var ret solofsapitypes.FsINodeMeta
	return ret, nil
}
