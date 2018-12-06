package metastg

import (
	"database/sql"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"

	"github.com/gocraft/dbr"
)

func (p *MetaStg) FetchNetBlock(pNetBlock *types.NetBlock) (exsists bool, err error) {
	var (
		sess    *dbr.Session
		sqlRows *sql.Rows
	)

	sess = p.DBConn.NewSession(nil)
	sqlRows, err = sess.Select("index_in_inode", "netblock_len", "netblock_cap").
		From("b_netblock").
		Where("netblock_id=?", pNetBlock.IDStr()).Rows()
	for sqlRows.Next() {
		sqlRows.Scan(&pNetBlock.IndexInInode, &pNetBlock.Len, &pNetBlock.Cap)
		exsists = true
	}
	err = sqlRows.Close()
	if err != nil {
		return
	}

	sqlRows, err = sess.Select("netblock_id", "peer_id").
		From("r_netblock_store_peer").
		Where("netblock_id=?", pNetBlock.IDStr()).Rows()
	for sqlRows.Next() {
		// TODO load datanodes
		// sqlRows.Scan(&pNetBlock.IndexInInode, &pNetBlock.Size)
		exsists = true
	}
	err = sqlRows.Close()
	if err != nil {
		return
	}

	return
}

func (p *MetaStg) insertNetBlockDataNodes(sess *dbr.Session, pNetBlock *types.NetBlock) error {
	var (
		netBlockIDStr = pNetBlock.IDStr()
		uPeer         snettypes.PeerUintptr
		err           error
	)
	for i := 0; i < pNetBlock.DataNodes.Len; i++ {
		uPeer = pNetBlock.DataNodes.Arr[i]
		_, err = sess.InsertInto("r_netblock_store_peer").
			Columns("netblock_id", "peer_id").
			Values(netBlockIDStr, uPeer.Ptr().PeerIDStr()).
			Exec()
		if err != nil {
			return err
		}
	}

	return err
}

func (p *MetaStg) StoreNetBlock(pINode *types.INode, pNetBlock *types.NetBlock, isRefreshDataNodes bool) error {
	var (
		sess          *dbr.Session
		tx            *dbr.Tx
		netBlockIDStr = pNetBlock.IDStr()
		inodeIDStr    = pINode.IDStr()
		err           error
	)

	sess = p.DBConn.NewSession(nil)
	tx, err = sess.Begin()
	if err != nil {
		return err
	}

	_, err = sess.InsertInto("b_netblock").
		Columns("netblock_id", "inode_id", "index_in_inode", "netblock_len", "netblock_cap").
		Values(netBlockIDStr, inodeIDStr, pNetBlock.IndexInInode, pNetBlock.Len, pNetBlock.Cap).
		Exec()
	if err != nil {
		_, err = sess.Update("b_netblock").
			Set("inode_id", inodeIDStr).
			Set("index_in_inode", pNetBlock.IndexInInode).
			Set("netblock_len", pNetBlock.Len).
			Set("netblock_cap", pNetBlock.Cap).
			Where("netblock_id=?", netBlockIDStr).
			Exec()
	}

	if err == nil && isRefreshDataNodes {
		_, err = sess.DeleteFrom("r_netblock_store_peer").
			Where("netblock_id=?", inodeIDStr).
			Exec()
		if err != nil {
			goto EXEC_DONE
		}

		p.insertNetBlockDataNodes(sess, pNetBlock)
	}

EXEC_DONE:
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
