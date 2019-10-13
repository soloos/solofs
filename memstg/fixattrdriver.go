package memstg

import (
	"bytes"
	"soloos/common/solofstypes"
	"soloos/common/util"
)

type FIXAttrDriverHelper struct {
	solofstypes.DeleteFIXAttrInDB
	solofstypes.ReplaceFIXAttrInDB
	solofstypes.GetFIXAttrByInoFromDB
}

// FIXAttrDriver is FsINode XAttr driver
type FIXAttrDriver struct {
	posixFs *PosixFs
	helper  FIXAttrDriverHelper

	xattrsRWMutex util.RWMutex
	xattrs        map[solofstypes.FsINodeIno]solofstypes.FsINodeXAttr
}

func (p *FIXAttrDriver) Init(
	posixFs *PosixFs,
	deleteFIXAttrInDB solofstypes.DeleteFIXAttrInDB,
	replaceFIXAttrInDB solofstypes.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB solofstypes.GetFIXAttrByInoFromDB,
) error {
	p.posixFs = posixFs
	p.SetHelper(
		deleteFIXAttrInDB,
		replaceFIXAttrInDB,
		getFIXAttrByInoFromDB,
	)

	p.xattrs = make(map[solofstypes.FsINodeIno]solofstypes.FsINodeXAttr)

	return nil
}

func (p *FIXAttrDriver) SetHelper(
	deleteFIXAttrInDB solofstypes.DeleteFIXAttrInDB,
	replaceFIXAttrInDB solofstypes.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB solofstypes.GetFIXAttrByInoFromDB,
) {
	p.helper = FIXAttrDriverHelper{
		DeleteFIXAttrInDB:     deleteFIXAttrInDB,
		ReplaceFIXAttrInDB:    replaceFIXAttrInDB,
		GetFIXAttrByInoFromDB: getFIXAttrByInoFromDB,
	}
}

func (p *FIXAttrDriver) getXAttrFromCache(fsINodeIno solofstypes.FsINodeIno) (solofstypes.FsINodeXAttr, bool) {
	var (
		xattr  solofstypes.FsINodeXAttr
		exists bool
	)
	p.xattrsRWMutex.RLock()
	xattr, exists = p.xattrs[fsINodeIno]
	p.xattrsRWMutex.RUnlock()
	return xattr, exists
}

func (p *FIXAttrDriver) setXAttrInCache(fsINodeIno solofstypes.FsINodeIno, xattr solofstypes.FsINodeXAttr) {
	p.xattrsRWMutex.Lock()
	p.xattrs[fsINodeIno] = xattr
	p.xattrsRWMutex.Unlock()
}

func (p *FIXAttrDriver) xAttrFetchAttr(xattr solofstypes.FsINodeXAttr, attr string) []byte {
	var ret []byte
	p.xattrsRWMutex.RLock()
	ret = xattr[attr]
	p.xattrsRWMutex.RUnlock()
	return ret
}

func (p *FIXAttrDriver) xAttrSetAttr(xattr solofstypes.FsINodeXAttr, attr string, paramData []byte) {
	var data = append([]byte{}, paramData...)
	p.xattrsRWMutex.Lock()
	xattr[attr] = data
	p.xattrsRWMutex.Unlock()
}

func (p *FIXAttrDriver) xAttrRemoveAttr(xattr solofstypes.FsINodeXAttr, attr string) {
	p.xattrsRWMutex.Lock()
	delete(xattr, attr)
	p.xattrsRWMutex.Unlock()
}

func (p *FIXAttrDriver) getXAttr(fsINodeIno solofstypes.FsINodeIno) (solofstypes.FsINodeXAttr, error) {
	var (
		xattr       solofstypes.FsINodeXAttr
		xattrExists bool
		err         error
	)

	xattr, xattrExists = p.getXAttrFromCache(fsINodeIno)
	if xattrExists {
		return xattr, nil
	}

	xattr, err = p.helper.GetFIXAttrByInoFromDB(p.posixFs.NameSpaceID, fsINodeIno)
	if err != nil && err.Error() != solofstypes.ErrObjectNotExists.Error() {
		return xattr, err
	}

	p.setXAttrInCache(fsINodeIno, xattr)
	return xattr, nil
}

func (p *FIXAttrDriver) GetXAttrSize(fsINodeIno solofstypes.FsINodeIno, attr string) (int, error) {
	var (
		value []byte
		err   error
	)
	value, err = p.GetXAttrData(fsINodeIno, attr)
	return len(value), err
}

func (p *FIXAttrDriver) GetXAttrData(fsINodeIno solofstypes.FsINodeIno, attr string) ([]byte, error) {
	var (
		xattr solofstypes.FsINodeXAttr
		value []byte
		err   error
	)
	xattr, err = p.getXAttr(fsINodeIno)
	if err != nil {
		return nil, err
	}
	value = p.xAttrFetchAttr(xattr, attr)
	return value, nil
}

func (p *FIXAttrDriver) ListXAttr(fsINodeIno solofstypes.FsINodeIno) ([]byte, error) {
	var (
		xattr solofstypes.FsINodeXAttr
		err   error
	)
	xattr, err = p.getXAttr(fsINodeIno)
	if err != nil {
		return nil, err
	}

	var b = bytes.NewBuffer([]byte{})
	for k, _ := range xattr {
		b.WriteString(k)
		b.WriteByte(0)
	}

	return b.Bytes(), nil
}

func (p *FIXAttrDriver) SetXAttr(fsINodeIno solofstypes.FsINodeIno, attr string, data []byte) error {
	var (
		xattr solofstypes.FsINodeXAttr
		err   error
	)

	xattr, err = p.getXAttr(fsINodeIno)
	if err != nil {
		return err
	}

	if xattr == nil {
		xattr = solofstypes.InitFsINodeXAttr()
	}

	p.xAttrSetAttr(xattr, attr, data)

	err = p.helper.ReplaceFIXAttrInDB(p.posixFs.NameSpaceID, fsINodeIno, xattr)
	if err != nil {
		return err
	}

	p.setXAttrInCache(fsINodeIno, xattr)

	return nil
}

func (p *FIXAttrDriver) RemoveXAttr(fsINodeIno solofstypes.FsINodeIno, attr string) error {
	var (
		xattr solofstypes.FsINodeXAttr
		err   error
	)

	xattr, err = p.getXAttr(fsINodeIno)
	if err != nil {
		return err
	}

	p.xAttrRemoveAttr(xattr, attr)
	err = p.helper.ReplaceFIXAttrInDB(p.posixFs.NameSpaceID, fsINodeIno, xattr)
	if err != nil {
		return err
	}

	return nil
}
