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
	xattrs        map[solofstypes.FsINodeID]solofstypes.FsINodeXAttr
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

	p.xattrs = make(map[solofstypes.FsINodeID]solofstypes.FsINodeXAttr)

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

func (p *FIXAttrDriver) getXAttrFromCache(fsINodeID solofstypes.FsINodeID) (solofstypes.FsINodeXAttr, bool) {
	var (
		xattr  solofstypes.FsINodeXAttr
		exists bool
	)
	p.xattrsRWMutex.RLock()
	xattr, exists = p.xattrs[fsINodeID]
	p.xattrsRWMutex.RUnlock()
	return xattr, exists
}

func (p *FIXAttrDriver) setXAttrInCache(fsINodeID solofstypes.FsINodeID, xattr solofstypes.FsINodeXAttr) {
	p.xattrsRWMutex.Lock()
	p.xattrs[fsINodeID] = xattr
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

func (p *FIXAttrDriver) getXAttr(fsINodeID solofstypes.FsINodeID) (solofstypes.FsINodeXAttr, error) {
	var (
		xattr       solofstypes.FsINodeXAttr
		xattrExists bool
		err         error
	)

	xattr, xattrExists = p.getXAttrFromCache(fsINodeID)
	if xattrExists {
		return xattr, nil
	}

	xattr, err = p.helper.GetFIXAttrByInoFromDB(p.posixFs.NameSpaceID, fsINodeID)
	if err != nil && err.Error() != solofstypes.ErrObjectNotExists.Error() {
		return xattr, err
	}

	p.setXAttrInCache(fsINodeID, xattr)
	return xattr, nil
}

func (p *FIXAttrDriver) GetXAttrSize(fsINodeID solofstypes.FsINodeID, attr string) (int, error) {
	var (
		value []byte
		err   error
	)
	value, err = p.GetXAttrData(fsINodeID, attr)
	return len(value), err
}

func (p *FIXAttrDriver) GetXAttrData(fsINodeID solofstypes.FsINodeID, attr string) ([]byte, error) {
	var (
		xattr solofstypes.FsINodeXAttr
		value []byte
		err   error
	)
	xattr, err = p.getXAttr(fsINodeID)
	if err != nil {
		return nil, err
	}
	value = p.xAttrFetchAttr(xattr, attr)
	return value, nil
}

func (p *FIXAttrDriver) ListXAttr(fsINodeID solofstypes.FsINodeID) ([]byte, error) {
	var (
		xattr solofstypes.FsINodeXAttr
		err   error
	)
	xattr, err = p.getXAttr(fsINodeID)
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

func (p *FIXAttrDriver) SetXAttr(fsINodeID solofstypes.FsINodeID, attr string, data []byte) error {
	var (
		xattr solofstypes.FsINodeXAttr
		err   error
	)

	xattr, err = p.getXAttr(fsINodeID)
	if err != nil {
		return err
	}

	if xattr == nil {
		xattr = solofstypes.InitFsINodeXAttr()
	}

	p.xAttrSetAttr(xattr, attr, data)

	err = p.helper.ReplaceFIXAttrInDB(p.posixFs.NameSpaceID, fsINodeID, xattr)
	if err != nil {
		return err
	}

	p.setXAttrInCache(fsINodeID, xattr)

	return nil
}

func (p *FIXAttrDriver) RemoveXAttr(fsINodeID solofstypes.FsINodeID, attr string) error {
	var (
		xattr solofstypes.FsINodeXAttr
		err   error
	)

	xattr, err = p.getXAttr(fsINodeID)
	if err != nil {
		return err
	}

	p.xAttrRemoveAttr(xattr, attr)
	err = p.helper.ReplaceFIXAttrInDB(p.posixFs.NameSpaceID, fsINodeID, xattr)
	if err != nil {
		return err
	}

	return nil
}
