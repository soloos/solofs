package metastg

import (
	"database/sql"
	"soloos/sdfs/types"

	"github.com/gocraft/dbr"
)

func (p *FIXAttrDriver) DeleteFIXAttrInDB(fsINodeID types.FsINodeID) error {
	var (
		sess *dbr.Session
		err  error
	)

	sess = p.DBConn.NewSession(nil)
	_, err = sess.DeleteFrom("b_fsinode_xattr").
		Where("fsinode_ino=?", fsINodeID).
		Exec()
	return err
}

func (p *FIXAttrDriver) ReplaceFIXAttrInDB(fsINodeID types.FsINodeID, xattr types.FsINodeXAttr) error {
	var (
		sess       *dbr.Session
		tx         *dbr.Tx
		xattrBytes []byte
		err        error
	)

	sess = p.DBConn.NewSession(nil)
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
		sess    *dbr.Session
		sqlRows *sql.Rows
		xattr   = types.InitFsINodeXAttr()
		bytes   []byte
		err     error
	)

	sess = p.DBConn.NewSession(nil)
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
