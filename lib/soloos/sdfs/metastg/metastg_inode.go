package metastg

import (
	"database/sql"
	"soloos/sdfs/types"

	"github.com/gocraft/dbr"
)

func (p *MetaStg) FetchINode(pINode *types.INode) (exsists bool, err error) {
	var (
		sess    *dbr.Session
		sqlRows *sql.Rows
	)

	sess = p.DBConn.NewSession(nil)
	sqlRows, err = sess.Select("inode_size", "netblock_size", "memblock_size").
		From("b_inode").
		Where("inode_id=?", pINode.IDStr()).Rows()
	if sqlRows == nil {
		return
	}
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
		sess       *dbr.Session
		tx         *dbr.Tx
		inodeIDStr = pINode.IDStr()
		err        error
	)

	sess = p.DBConn.NewSession(nil)
	tx, err = sess.Begin()
	if err != nil {
		return err
	}

	_, err = sess.InsertInto("b_inode").
		Columns("inode_id", "inode_size", "netblock_size", "memblock_size").
		Values(inodeIDStr, pINode.Size, pINode.NetBlockSize, pINode.MemBlockSize).
		Exec()
	if err != nil {
		_, err = sess.Update("b_inode").
			Set("inode_size", pINode.Size).
			Set("netblock_size", pINode.NetBlockSize).
			Set("memblock_size", pINode.MemBlockSize).
			Where("inode_id=?", inodeIDStr).
			Exec()
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
