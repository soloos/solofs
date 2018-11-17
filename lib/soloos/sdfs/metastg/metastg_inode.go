package metastg

import (
	"database/sql"
	"soloos/sdfs/types"
	"strings"

	"github.com/gocraft/dbr"
)

func (p *MetaStg) FetchINode(pINode *types.INode) (exsists bool, err error) {
	var (
		sess    *dbr.Session
		sqlRows *sql.Rows
	)

	sess = p.DBConn.NewSession(nil)
	sqlRows, err = sess.Select("inodesize", "netblocksize", "memblocksize").
		From("b_inode").
		Where("inodeid=?", pINode.IDStr()).Rows()
	for sqlRows.Next() {
		sqlRows.Scan(&pINode.Size, &pINode.NetBlockSize, &pINode.MemBlockSize)
		exsists = true
	}
	err = sqlRows.Close()
	if err != nil {
		return
	}

	return
}

func (p *MetaStg) StoreINode(pINode *types.INode) error {
	var (
		sess *dbr.Session
		tx   *dbr.Tx
		err  error
	)

	sess = p.DBConn.NewSession(nil)
	tx, err = sess.Begin()
	if err != nil {
		return err
	}

	_, err = sess.InsertInto("b_inode").
		Columns("inodeid", "inodesize", "netblocksize", "memblocksize").
		Values(pINode.IDStr(), pINode.Size, pINode.NetBlockSize, pINode.MemBlockSize).
		Exec()
	if err != nil {
		if strings.Index(err.Error(), "Duplicate entry") >= 0 {
			_, err = sess.Update("b_inode").
				Set("inodesize", pINode.Size).
				Set("netblocksize", pINode.NetBlockSize).
				Set("memblocksize", pINode.MemBlockSize).
				Where("inodeid=?", pINode.IDStr()).
				Exec()
		}
	}

	if err != nil {
		tx.RollbackUnlessCommitted()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
