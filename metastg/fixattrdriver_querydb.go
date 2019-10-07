package metastg

import (
	"database/sql"
	"soloos/common/solodbapi"
	"soloos/common/solofstypes"
)

func (p *FIXAttrDriver) DeleteFIXAttrInDB(nsID solofstypes.NameSpaceID,
	fsINodeID solofstypes.FsINodeID) error {
	var (
		sess solodbapi.Session
		err  error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	_, err = sess.DeleteFrom("b_fsinode_xattr").
		Where("namespace_id=? and fsinode_ino=?", nsID, fsINodeID).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *FIXAttrDriver) ReplaceFIXAttrInDB(nsID solofstypes.NameSpaceID,
	fsINodeID solofstypes.FsINodeID,
	xattr solofstypes.FsINodeXAttr) error {
	var (
		sess       solodbapi.Session
		xattrBytes []byte
		err        error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	xattrBytes, err = solofstypes.SerializeFIXAttr(xattr)
	if err != nil {
		return err
	}

	err = sess.ReplaceInto("b_fsinode_xattr").
		PrimaryColumns("namespace_id", "fsinode_ino").PrimaryValues(nsID, fsINodeID).
		Columns("xattr").Values(xattrBytes).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *FIXAttrDriver) GetFIXAttrByInoFromDB(nsID solofstypes.NameSpaceID,
	fsINodeID solofstypes.FsINodeID) (solofstypes.FsINodeXAttr, error) {
	var (
		sess    solodbapi.Session
		sqlRows *sql.Rows
		xattr   = solofstypes.InitFsINodeXAttr()
		bytes   []byte
		err     error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		goto QUERY_DONE
	}

	sqlRows, err = sess.Select("xattr").
		From("b_fsinode_xattr").
		Where("namespace_id=? and fsinode_ino=?", nsID, fsINodeID).
		Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = solofstypes.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(
		&bytes,
	)
	if err != nil {
		goto QUERY_DONE
	}
	solofstypes.DeserializeFIXAttr(bytes, &xattr)

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return xattr, err
}
