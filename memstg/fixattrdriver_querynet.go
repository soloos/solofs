package memstg

import (
	"soloos/common/snet"
	"soloos/common/solofstypes"
)

func (p *FIXAttrDriver) DeleteFIXAttrInDB(
	nsID solofstypes.NameSpaceID,
	fsINodeIno solofstypes.FsINodeIno) error {
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FIXAttr/DeleteFIXAttrInDB", nil,
		nsID, fsINodeIno)
	if err != nil {
		return err
	}
	return nil
}

func (p *FIXAttrDriver) ReplaceFIXAttrInDB(
	nsID solofstypes.NameSpaceID,
	fsINodeIno solofstypes.FsINodeIno,
	xattr solofstypes.FsINodeXAttr) error {
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FIXAttr/ReplaceFIXAttrInDB", nil,
		nsID, fsINodeIno, xattr)
	if err != nil {
		return err
	}
	return nil
}

func (p *FIXAttrDriver) GetFIXAttrByInoFromDB(
	nsID solofstypes.NameSpaceID,
	fsINodeIno solofstypes.FsINodeIno) (solofstypes.FsINodeXAttr, error) {
	var ret = snet.Response{RespData: solofstypes.FsINodeXAttr{}}
	var err = p.posixFs.MemStg.SolonnClient.Dispatch("/FIXAttr/GetFIXAttrByInoFromDB", &ret,
		nsID, fsINodeIno)
	if err != nil {
		return solofstypes.FsINodeXAttr{}, err
	}
	return ret.RespData.(solofstypes.FsINodeXAttr), nil
}
