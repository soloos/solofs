package memstg

import (
	"soloos/common/snettypes"
	"soloos/common/solofsapitypes"
)

func (p *FIXAttrDriver) DeleteFIXAttrInDB(
	nsID solofsapitypes.NameSpaceID,
	fsINodeID solofsapitypes.FsINodeID) error {
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FIXAttr/DeleteFIXAttrInDB", nil,
		nsID, fsINodeID)
	if err != nil {
		return err
	}
	return nil
}

func (p *FIXAttrDriver) ReplaceFIXAttrInDB(
	nsID solofsapitypes.NameSpaceID,
	fsINodeID solofsapitypes.FsINodeID,
	xattr solofsapitypes.FsINodeXAttr) error {
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FIXAttr/ReplaceFIXAttrInDB", nil,
		nsID, fsINodeID, xattr)
	if err != nil {
		return err
	}
	return nil
}

func (p *FIXAttrDriver) GetFIXAttrByInoFromDB(
	nsID solofsapitypes.NameSpaceID,
	fsINodeID solofsapitypes.FsINodeID) (solofsapitypes.FsINodeXAttr, error) {
	var ret = snettypes.Response{RespData: solofsapitypes.FsINodeXAttr{}}
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FIXAttr/GetFIXAttrByInoFromDB", &ret,
		nsID, fsINodeID)
	if err != nil {
		return solofsapitypes.FsINodeXAttr{}, err
	}
	return ret.RespData.(solofsapitypes.FsINodeXAttr), nil
}
