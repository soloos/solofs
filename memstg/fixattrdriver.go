package memstg

import (
	"bytes"
	"soloos/common/solofsapitypes"
	"soloos/common/util"
)

type FIXAttrDriverHelper struct {
	solofsapitypes.DeleteFIXAttrInDB
	solofsapitypes.ReplaceFIXAttrInDB
	solofsapitypes.GetFIXAttrByInoFromDB
}

// FIXAttrDriver is FsINode XAttr driver
type FIXAttrDriver struct {
	posixFs *PosixFs
	helper  FIXAttrDriverHelper

	xattrsRWMutex util.RWMutex
	xattrs        map[solofsapitypes.FsINodeID]solofsapitypes.FsINodeXAttr
}

func (p *FIXAttrDriver) Init(
	posixFs *PosixFs,
	deleteFIXAttrInDB solofsapitypes.DeleteFIXAttrInDB,
	replaceFIXAttrInDB solofsapitypes.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB solofsapitypes.GetFIXAttrByInoFromDB,
) error {
	p.posixFs = posixFs
	p.SetHelper(
		deleteFIXAttrInDB,
		replaceFIXAttrInDB,
		getFIXAttrByInoFromDB,
	)

	p.xattrs = make(map[solofsapitypes.FsINodeID]solofsapitypes.FsINodeXAttr)

	return nil
}

func (p *FIXAttrDriver) SetHelper(
	deleteFIXAttrInDB solofsapitypes.DeleteFIXAttrInDB,
	replaceFIXAttrInDB solofsapitypes.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB solofsapitypes.GetFIXAttrByInoFromDB,
) {
	p.helper = FIXAttrDriverHelper{
		DeleteFIXAttrInDB:     deleteFIXAttrInDB,
		ReplaceFIXAttrInDB:    replaceFIXAttrInDB,
		GetFIXAttrByInoFromDB: getFIXAttrByInoFromDB,
	}
}

func (p *FIXAttrDriver) getXAttrFromCache(fsINodeID solofsapitypes.FsINodeID) (solofsapitypes.FsINodeXAttr, bool) {
	var (
		xattr  solofsapitypes.FsINodeXAttr
		exists bool
	)
	p.xattrsRWMutex.RLock()
	xattr, exists = p.xattrs[fsINodeID]
	p.xattrsRWMutex.RUnlock()
	return xattr, exists
}

func (p *FIXAttrDriver) setXAttrInCache(fsINodeID solofsapitypes.FsINodeID, xattr solofsapitypes.FsINodeXAttr) {
	p.xattrsRWMutex.Lock()
	p.xattrs[fsINodeID] = xattr
	p.xattrsRWMutex.Unlock()
}

func (p *FIXAttrDriver) xAttrFetchAttr(xattr solofsapitypes.FsINodeXAttr, attr string) []byte {
	var ret []byte
	p.xattrsRWMutex.RLock()
	ret = xattr[attr]
	p.xattrsRWMutex.RUnlock()
	return ret
}

func (p *FIXAttrDriver) xAttrSetAttr(xattr solofsapitypes.FsINodeXAttr, attr string, paramData []byte) {
	var data = append([]byte{}, paramData...)
	p.xattrsRWMutex.Lock()
	xattr[attr] = data
	p.xattrsRWMutex.Unlock()
}

func (p *FIXAttrDriver) xAttrRemoveAttr(xattr solofsapitypes.FsINodeXAttr, attr string) {
	p.xattrsRWMutex.Lock()
	delete(xattr, attr)
	p.xattrsRWMutex.Unlock()
}

func (p *FIXAttrDriver) getXAttr(fsINodeID solofsapitypes.FsINodeID) (solofsapitypes.FsINodeXAttr, error) {
	var (
		xattr       solofsapitypes.FsINodeXAttr
		xattrExists bool
		err         error
	)

	xattr, xattrExists = p.getXAttrFromCache(fsINodeID)
	if xattrExists {
		return xattr, nil
	}

	xattr, err = p.helper.GetFIXAttrByInoFromDB(p.posixFs.NameSpaceID, fsINodeID)
	if err != nil && err != solofsapitypes.ErrObjectNotExists {
		return xattr, err
	}

	p.setXAttrInCache(fsINodeID, xattr)
	return xattr, nil
}

func (p *FIXAttrDriver) GetXAttrSize(fsINodeID solofsapitypes.FsINodeID, attr string) (int, error) {
	var (
		value []byte
		err   error
	)
	value, err = p.GetXAttrData(fsINodeID, attr)
	return len(value), err
}

func (p *FIXAttrDriver) GetXAttrData(fsINodeID solofsapitypes.FsINodeID, attr string) ([]byte, error) {
	var (
		xattr solofsapitypes.FsINodeXAttr
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

func (p *FIXAttrDriver) ListXAttr(fsINodeID solofsapitypes.FsINodeID) ([]byte, error) {
	var (
		xattr solofsapitypes.FsINodeXAttr
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

func (p *FIXAttrDriver) SetXAttr(fsINodeID solofsapitypes.FsINodeID, attr string, data []byte) error {
	var (
		xattr solofsapitypes.FsINodeXAttr
		err   error
	)

	xattr, err = p.getXAttr(fsINodeID)
	if err != nil {
		return err
	}

	if xattr == nil {
		xattr = solofsapitypes.InitFsINodeXAttr()
	}

	p.xAttrSetAttr(xattr, attr, data)

	err = p.helper.ReplaceFIXAttrInDB(p.posixFs.NameSpaceID, fsINodeID, xattr)
	if err != nil {
		return err
	}

	p.setXAttrInCache(fsINodeID, xattr)

	return nil
}

func (p *FIXAttrDriver) RemoveXAttr(fsINodeID solofsapitypes.FsINodeID, attr string) error {
	var (
		xattr solofsapitypes.FsINodeXAttr
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
