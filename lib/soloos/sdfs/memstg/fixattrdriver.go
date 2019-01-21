package memstg

import (
	"bytes"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"sync"
)

type FIXAttrDriverHelper struct {
	api.DeleteFIXAttrInDB
	api.ReplaceFIXAttrInDB
	api.GetFIXAttrByInoFromDB
}

// FIXAttrDriver is FsINode XAttr driver
type FIXAttrDriver struct {
	helper FIXAttrDriverHelper

	xattrsRWMutex sync.RWMutex
	xattrs        map[types.FsINodeID]types.FsINodeXAttr
}

func (p *FIXAttrDriver) Init(
	deleteFIXAttrInDB api.DeleteFIXAttrInDB,
	replaceFIXAttrInDB api.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB api.GetFIXAttrByInoFromDB,
) error {
	p.SetHelper(
		deleteFIXAttrInDB,
		replaceFIXAttrInDB,
		getFIXAttrByInoFromDB,
	)

	p.xattrs = make(map[types.FsINodeID]types.FsINodeXAttr)

	return nil
}

func (p *FIXAttrDriver) SetHelper(
	deleteFIXAttrInDB api.DeleteFIXAttrInDB,
	replaceFIXAttrInDB api.ReplaceFIXAttrInDB,
	getFIXAttrByInoFromDB api.GetFIXAttrByInoFromDB,
) {
	p.helper = FIXAttrDriverHelper{
		DeleteFIXAttrInDB:     deleteFIXAttrInDB,
		ReplaceFIXAttrInDB:    replaceFIXAttrInDB,
		GetFIXAttrByInoFromDB: getFIXAttrByInoFromDB,
	}
}

func (p *FIXAttrDriver) getXAttrFromCache(fsINodeID types.FsINodeID) (types.FsINodeXAttr, bool) {
	var (
		xattr  types.FsINodeXAttr
		exists bool
	)
	p.xattrsRWMutex.RLock()
	xattr, exists = p.xattrs[fsINodeID]
	p.xattrsRWMutex.RUnlock()
	return xattr, exists
}

func (p *FIXAttrDriver) setXAttrInCache(fsINodeID types.FsINodeID, xattr types.FsINodeXAttr) {
	p.xattrsRWMutex.Lock()
	p.xattrs[fsINodeID] = xattr
	p.xattrsRWMutex.Unlock()
}

func (p *FIXAttrDriver) xAttrFetchAttr(xattr types.FsINodeXAttr, attr string) []byte {
	var ret []byte
	p.xattrsRWMutex.RLock()
	ret = xattr[attr]
	p.xattrsRWMutex.RUnlock()
	return ret
}

func (p *FIXAttrDriver) xAttrSetAttr(xattr types.FsINodeXAttr, attr string, paramData []byte) {
	var data = append([]byte{}, paramData...)
	p.xattrsRWMutex.Lock()
	xattr[attr] = data
	p.xattrsRWMutex.Unlock()
}

func (p *FIXAttrDriver) xAttrRemoveAttr(xattr types.FsINodeXAttr, attr string) {
	p.xattrsRWMutex.Lock()
	delete(xattr, attr)
	p.xattrsRWMutex.Unlock()
}

func (p *FIXAttrDriver) getXAttr(fsINodeID types.FsINodeID) (types.FsINodeXAttr, error) {
	var (
		xattr       types.FsINodeXAttr
		xattrExists bool
		err         error
	)

	xattr, xattrExists = p.getXAttrFromCache(fsINodeID)
	if xattrExists {
		return xattr, nil
	}

	xattr, err = p.helper.GetFIXAttrByInoFromDB(fsINodeID)
	if err != nil && err != types.ErrObjectNotExists {
		return xattr, err
	}

	p.setXAttrInCache(fsINodeID, xattr)
	return xattr, nil
}

func (p *FIXAttrDriver) GetXAttrSize(fsINodeID types.FsINodeID, attr string) (int, error) {
	var (
		value []byte
		err   error
	)
	value, err = p.GetXAttrData(fsINodeID, attr)
	return len(value), err
}

func (p *FIXAttrDriver) GetXAttrData(fsINodeID types.FsINodeID, attr string) ([]byte, error) {
	var (
		xattr types.FsINodeXAttr
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

func (p *FIXAttrDriver) ListXAttr(fsINodeID types.FsINodeID) ([]byte, error) {
	var (
		xattr types.FsINodeXAttr
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

func (p *FIXAttrDriver) SetXAttr(fsINodeID types.FsINodeID, attr string, data []byte) error {
	var (
		xattr types.FsINodeXAttr
		err   error
	)

	xattr, err = p.getXAttr(fsINodeID)
	if err != nil {
		return err
	}

	if xattr == nil {
		xattr = types.InitFsINodeXAttr()
	}

	p.xAttrSetAttr(xattr, attr, data)

	err = p.helper.ReplaceFIXAttrInDB(fsINodeID, xattr)
	if err != nil {
		return err
	}

	p.setXAttrInCache(fsINodeID, xattr)

	return nil
}

func (p *FIXAttrDriver) RemoveXAttr(fsINodeID types.FsINodeID, attr string) error {
	var (
		xattr types.FsINodeXAttr
		err   error
	)

	xattr, err = p.getXAttr(fsINodeID)
	if err != nil {
		return err
	}

	p.xAttrRemoveAttr(xattr, attr)
	err = p.helper.ReplaceFIXAttrInDB(fsINodeID, xattr)
	if err != nil {
		return err
	}

	return nil
}
