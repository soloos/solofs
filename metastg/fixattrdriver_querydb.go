package metastg

import (
	"database/sql"
	"soloos/common/sdbapi"
	"soloos/common/sdfsapitypes"
)

func (p *FIXAttrDriver) DeleteFIXAttrInDB(nameSpaceID sdfsapitypes.NameSpaceID,
	fsINodeID sdfsapitypes.FsINodeID) error {
	var (
		sess sdbapi.Session
		err  error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	_, err = sess.DeleteFrom("b_fsinode_xattr").
		Where("namespace_id=? and fsinode_ino=?", nameSpaceID, fsINodeID).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *FIXAttrDriver) ReplaceFIXAttrInDB(nameSpaceID sdfsapitypes.NameSpaceID,
	fsINodeID sdfsapitypes.FsINodeID,
	xattr sdfsapitypes.FsINodeXAttr) error {
	var (
		sess       sdbapi.Session
		xattrBytes []byte
		err        error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	xattrBytes, err = sdfsapitypes.SerializeFIXAttr(xattr)
	if err != nil {
		return err
	}

	err = sess.ReplaceInto("b_fsinode_xattr").
		PrimaryColumns("namespace_id", "fsinode_ino").PrimaryValues(nameSpaceID, fsINodeID).
		Columns("xattr").Values(xattrBytes).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *FIXAttrDriver) GetFIXAttrByInoFromDB(nameSpaceID sdfsapitypes.NameSpaceID,
	fsINodeID sdfsapitypes.FsINodeID) (sdfsapitypes.FsINodeXAttr, error) {
	var (
		sess    sdbapi.Session
		sqlRows *sql.Rows
		xattr   = sdfsapitypes.InitFsINodeXAttr()
		bytes   []byte
		err     error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		goto QUERY_DONE
	}

	sqlRows, err = sess.Select("xattr").
		From("b_fsinode_xattr").
		Where("namespace_id=? and fsinode_ino=?", nameSpaceID, fsINodeID).
		Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = sdfsapitypes.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(
		&bytes,
	)
	if err != nil {
		goto QUERY_DONE
	}
	sdfsapitypes.DeserializeFIXAttr(bytes, &xattr)

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return xattr, err
}
