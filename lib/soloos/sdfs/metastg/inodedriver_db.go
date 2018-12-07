package metastg

import (
	"database/sql"
	"soloos/sdfs/types"

	"github.com/gocraft/dbr"
)

func (p *INodeDriver) FetchINodeFromDB(pINode *types.INode) error {
	var (
		sess    *dbr.Session
		sqlRows *sql.Rows
		err     error
	)

	sess = p.metaStg.DBConn.NewSession(nil)
	sqlRows, err = sess.Select("inode_size", "netblock_cap", "memblock_cap").
		From("b_inode").
		Where("inode_id=?", pINode.IDStr()).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = types.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(&pINode.Size, &pINode.NetBlockCap, &pINode.MemBlockCap)
	if err != nil {
		goto QUERY_DONE
	}

QUERY_DONE:
	sqlRows.Close()
	return err
}

func (p *INodeDriver) StoreINodeInDB(pINode *types.INode) error {
	var (
		sess       *dbr.Session
		tx         *dbr.Tx
		inodeIDStr = pINode.IDStr()
		err        error
	)

	sess = p.metaStg.DBConn.NewSession(nil)
	tx, err = sess.Begin()
	if err != nil {
		goto QUERY_DONE
	}

	_, err = sess.InsertInto("b_inode").
		Columns("inode_id", "inode_size", "netblock_cap", "memblock_cap").
		Values(inodeIDStr, pINode.Size, pINode.NetBlockCap, pINode.MemBlockCap).
		Exec()
	if err != nil {
		_, err = sess.Update("b_inode").
			Set("inode_size", pINode.Size).
			Set("netblock_cap", pINode.NetBlockCap).
			Set("memblock_cap", pINode.MemBlockCap).
			Where("inode_id=?", inodeIDStr).
			Exec()
	}

QUERY_DONE:
	if err != nil {
		tx.RollbackUnlessCommitted()
	} else {
		err = tx.Commit()
	}
	return err
}
