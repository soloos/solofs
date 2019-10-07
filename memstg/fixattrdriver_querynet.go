package memstg

import (
	"soloos/common/snet"
	"soloos/common/solofstypes"
)

func (p *FIXAttrDriver) DeleteFIXAttrInDB(
	nsID solofstypes.NameSpaceID,
	fsINodeID solofstypes.FsINodeID) error {
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FIXAttr/DeleteFIXAttrInDB", nil,
		nsID, fsINodeID)
	if err != nil {
		return err
	}
	return nil
}

func (p *FIXAttrDriver) ReplaceFIXAttrInDB(
	nsID solofstypes.NameSpaceID,
	fsINodeID solofstypes.FsINodeID,
	xattr solofstypes.FsINodeXAttr) error {
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FIXAttr/ReplaceFIXAttrInDB", nil,
		nsID, fsINodeID, xattr)
	if err != nil {
		return err
	}
	return nil
}

func (p *FIXAttrDriver) GetFIXAttrByInoFromDB(
	nsID solofstypes.NameSpaceID,
	fsINodeID solofstypes.FsINodeID) (solofstypes.FsINodeXAttr, error) {
	var ret = snet.Response{RespData: solofstypes.FsINodeXAttr{}}
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FIXAttr/GetFIXAttrByInoFromDB", &ret,
		nsID, fsINodeID)
	if err != nil {
		return solofstypes.FsINodeXAttr{}, err
	}
	return ret.RespData.(solofstypes.FsINodeXAttr), nil
}
