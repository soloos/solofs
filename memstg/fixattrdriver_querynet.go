package memstg

import "soloos/common/solofsapitypes"

func (p *FIXAttrDriver) DeleteFIXAttrInDB(nameSpaceID solofsapitypes.NameSpaceID,
	fsINodeID solofsapitypes.FsINodeID) error {
	return nil
}

func (p *FIXAttrDriver) ReplaceFIXAttrInDB(nameSpaceID solofsapitypes.NameSpaceID,
	fsINodeID solofsapitypes.FsINodeID,
	xattr solofsapitypes.FsINodeXAttr) error {
	return nil
}

func (p *FIXAttrDriver) GetFIXAttrByInoFromDB(nameSpaceID solofsapitypes.NameSpaceID,
	fsINodeID solofsapitypes.FsINodeID) (solofsapitypes.FsINodeXAttr, error) {
	var ret solofsapitypes.FsINodeXAttr
	return ret, nil
}
