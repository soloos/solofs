package metastg

import (
	"database/sql"
	"soloos/common/sdbapi"
	"soloos/sdfs/types"
)

func (p *FIXAttrDriver) DeleteFIXAttrInDB(fsINodeID types.FsINodeID) error {
	var (
		sess sdbapi.Session
		err  error
	)

	err = p.DBConn.InitSession(&sess)
	if err != nil {
		return err
	}

	_, err = sess.DeleteFrom("b_fsinode_xattr").
		Where("fsinode_ino=?", fsINodeID).
		Exec()
	return err
}

func (p *FIXAttrDriver) ReplaceFIXAttrInDB(fsINodeID types.FsINodeID, xattr types.FsINodeXAttr) error {
	var (
		sess       sdbapi.Session
		tx         *sdbapi.Tx
		xattrBytes []byte
		err        error
	)

	err = p.DBConn.InitSession(&sess)
	if err != nil {
		goto QUERY_DONE
	}

	tx, err = sess.Begin()
	if err != nil {
		goto QUERY_DONE
	}

	xattrBytes, err = types.SerializeFIXAttr(xattr)
	if err != nil {
		goto QUERY_DONE
	}

	err = p.DBConn.ReplaceInto(tx, "b_fsinode_xattr", "fsinode_ino", "xattr", fsINodeID, xattrBytes)
	if err != nil {
		goto QUERY_DONE
	}

QUERY_DONE:
	if err != nil {
		tx.RollbackUnlessCommitted()
	} else {
		err = tx.Commit()
	}
	return err
}

func (p *FIXAttrDriver) GetFIXAttrByInoFromDB(fsINodeID types.FsINodeID) (types.FsINodeXAttr, error) {
	var (
		sess    sdbapi.Session
		sqlRows *sql.Rows
		xattr   = types.InitFsINodeXAttr()
		bytes   []byte
		err     error
	)

	err = p.DBConn.InitSession(&sess)
	if err != nil {
		goto QUERY_DONE
	}

	sqlRows, err = sess.Select("xattr").
		From("b_fsinode_xattr").
		Where("fsinode_ino=?",
			fsINodeID,
		).Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = types.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(
		&bytes,
	)
	if err != nil {
		goto QUERY_DONE
	}
	types.DeserializeFIXAttr(bytes, &xattr)

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return xattr, err
}
