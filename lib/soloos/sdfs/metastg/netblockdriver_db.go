package metastg

import (
	"bytes"
	"database/sql"
	"soloos/sdfs/types"

	"github.com/gocraft/dbr"
)

func (p *NetBlockDriver) FetchNetBlockFromDB(pNetBlock *types.NetBlock, backendPeerIDArrStr *string) (err error) {
	var (
		sess    *dbr.Session
		sqlRows *sql.Rows
	)

	sess = p.metaStg.DBConn.NewSession(nil)
	sqlRows, err = sess.Select("index_in_inode", "netblock_len", "netblock_cap", "backend_peer_id_arr").
		From("b_netblock").
		Where("netblock_id=?", pNetBlock.IDStr()).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = types.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(&pNetBlock.IndexInInode, &pNetBlock.Len, &pNetBlock.Cap, backendPeerIDArrStr)
	if err != nil {
		goto QUERY_DONE
	}

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return err
}

func (p *NetBlockDriver) StoreNetBlockInDB(pINode *types.INode, pNetBlock *types.NetBlock) error {
	var (
		sess                *dbr.Session
		tx                  *dbr.Tx
		netBlockIDStr       = pNetBlock.IDStr()
		inodeIDStr          = pINode.IDStr()
		backendPeerIDArrStr bytes.Buffer
		i                   int
		err                 error
	)

	sess = p.metaStg.DBConn.NewSession(nil)
	tx, err = sess.Begin()
	if err != nil {
		goto QUERY_DONE
	}

	if pNetBlock.DataNodes.Len > 0 {
		backendPeerIDArrStr.WriteString(pNetBlock.DataNodes.Arr[0].Ptr().PeerIDStr())
		for i = 1; i < pNetBlock.DataNodes.Len; i++ {
			backendPeerIDArrStr.WriteString(",")
			backendPeerIDArrStr.WriteString(pNetBlock.DataNodes.Arr[i].Ptr().PeerIDStr())
		}
	} else {
		backendPeerIDArrStr.WriteString("")
	}

	_, err = sess.InsertInto("b_netblock").
		Columns("netblock_id", "inode_id", "index_in_inode", "netblock_len", "netblock_cap", "backend_peer_id_arr").
		Values(netBlockIDStr, inodeIDStr, pNetBlock.IndexInInode, pNetBlock.Len, pNetBlock.Cap,
			backendPeerIDArrStr.String()).
		Exec()
	if err != nil {
		_, err = sess.Update("b_netblock").
			Set("inode_id", inodeIDStr).
			Set("index_in_inode", pNetBlock.IndexInInode).
			Set("netblock_len", pNetBlock.Len).
			Set("netblock_cap", pNetBlock.Cap).
			Set("backend_peer_id_arr", backendPeerIDArrStr.String()).
			Where("netblock_id=?", netBlockIDStr).
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
