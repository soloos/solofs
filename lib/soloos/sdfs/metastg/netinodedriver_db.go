package metastg

import (
	"database/sql"
	"soloos/sdfs/types"

	"github.com/gocraft/dbr"
)

func (p *NetINodeDriver) FetchNetINodeFromDB(pNetINode *types.NetINode) error {
	var (
		sess    *dbr.Session
		sqlRows *sql.Rows
		err     error
	)

	sess = p.metaStg.DBConn.NewSession(nil)
	sqlRows, err = sess.Select("netnetINode_size", "netblock_cap", "memblock_cap").
		From("b_netnetINode").
		Where("netnetINode_id=?", pNetINode.IDStr()).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = types.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(&pNetINode.Size, &pNetINode.NetBlockCap, &pNetINode.MemBlockCap)
	if err != nil {
		goto QUERY_DONE
	}

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return err
}

func (p *NetINodeDriver) StoreNetINodeInDB(pNetINode *types.NetINode) error {
	var (
		sess       *dbr.Session
		tx         *dbr.Tx
		netINodeIDStr = pNetINode.IDStr()
		err        error
	)

	sess = p.metaStg.DBConn.NewSession(nil)
	tx, err = sess.Begin()
	if err != nil {
		goto QUERY_DONE
	}

	_, err = sess.InsertInto("b_netnetINode").
		Columns("netnetINode_id", "netnetINode_size", "netblock_cap", "memblock_cap").
		Values(netINodeIDStr, pNetINode.Size, pNetINode.NetBlockCap, pNetINode.MemBlockCap).
		Exec()
	if err != nil {
		_, err = sess.Update("b_netnetINode").
			Set("netnetINode_size", pNetINode.Size).
			Set("netblock_cap", pNetINode.NetBlockCap).
			Set("memblock_cap", pNetINode.MemBlockCap).
			Where("netnetINode_id=?", netINodeIDStr).
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
