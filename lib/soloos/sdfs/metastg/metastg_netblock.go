package metastg

import (
	"database/sql"
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
	"strings"

	"github.com/gocraft/dbr"
)

func (p *MetaStg) FetchNetBlock(pNetBlock *types.NetBlock) (exsists bool, err error) {
	var (
		sess    *dbr.Session
		sqlRows *sql.Rows
	)

	sess = p.DBConn.NewSession(nil)
	sqlRows, err = sess.Select("index_in_inode", "netblocksize").
		From("b_netblock").
		Where("netblockid=?", pNetBlock.IDStr()).Rows()
	for sqlRows.Next() {
		sqlRows.Scan(&pNetBlock.IndexInInode, &pNetBlock.Size)
		exsists = true
	}
	err = sqlRows.Close()
	if err != nil {
		return
	}

	sqlRows, err = sess.Select("netblockid", "peerid").
		From("r_netblock_store_peer").
		Where("netblockid=?", pNetBlock.IDStr()).Rows()
	for sqlRows.Next() {
		// todo load datanodes
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
			Columns("netblockid", "peerid").
			Values(netBlockIDStr, uPeer.Ptr().IDStr()).
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
		Columns("netblockid", "inodeid", "index_in_inode", "netblocksize").
		Values(netBlockIDStr, inodeIDStr, pNetBlock.IndexInInode, pNetBlock.Size).
		Exec()
	if err != nil {
		if strings.Index(err.Error(), "Duplicate entry") >= 0 {
			_, err = sess.Update("b_netblock").
				Set("inodeid", inodeIDStr).
				Set("index_in_inode", pNetBlock.IndexInInode).
				Set("netblocksize", pNetBlock.Size).
				Where("netblockid=?", netBlockIDStr).
				Exec()
		}
	}

	if err == nil && isRefreshDataNodes {
		_, err = sess.DeleteFrom("r_netblock_store_peer").
			Where("netblockid=?", inodeIDStr).
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
