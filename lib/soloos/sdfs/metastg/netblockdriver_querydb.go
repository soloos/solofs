package metastg

import (
	"bytes"
	"database/sql"
	"soloos/common/sdbapi"
	"soloos/common/sdfsapitypes"
)

func (p *NetBlockDriver) FetchNetBlockFromDB(pNetINode *sdfsapitypes.NetINode,
	netBlockIndex int32, pNetBlock *sdfsapitypes.NetBlock,
	backendPeerIDArrStr *string) (err error) {
	var (
		sess    sdbapi.Session
		sqlRows *sql.Rows
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		goto QUERY_DONE
	}

	sqlRows, err = sess.Select("netblock_len", "netblock_cap", "backend_peer_id_arr").
		From("b_netblock").
		Where("netinode_id=? and index_in_netinode=?",
			pNetINode.IDStr(), netBlockIndex,
		).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = sdfsapitypes.ErrObjectNotExists
		goto QUERY_DONE
	}

	pNetBlock.IndexInNetINode = netBlockIndex
	pNetBlock.NetINodeID = pNetINode.ID
	err = sqlRows.Scan(&pNetBlock.Len, &pNetBlock.Cap, backendPeerIDArrStr)
	if err != nil {
		goto QUERY_DONE
	}

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return err
}

func (p *NetBlockDriver) StoreNetBlockInDB(pNetINode *sdfsapitypes.NetINode, pNetBlock *sdfsapitypes.NetBlock) error {
	var (
		sess                sdbapi.Session
		netINodeIDStr       = pNetINode.IDStr()
		backendPeerIDArrStr bytes.Buffer
		i                   int
		err                 error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	if pNetBlock.StorDataBackends.Len > 0 {
		backendPeerIDArrStr.WriteString(pNetBlock.StorDataBackends.Arr[0].Ptr().PeerIDStr())
		for i = 1; i < pNetBlock.StorDataBackends.Len; i++ {
			backendPeerIDArrStr.WriteString(",")
			backendPeerIDArrStr.WriteString(pNetBlock.StorDataBackends.Arr[i].Ptr().PeerIDStr())
		}
	} else {
		backendPeerIDArrStr.WriteString("")
	}

	err = sess.ReplaceInto("b_netblock").
		PrimaryColumns("netinode_id", "index_in_netinode").PrimaryValues(netINodeIDStr, pNetBlock.IndexInNetINode).
		Columns("netblock_len", "netblock_cap", "backend_peer_id_arr").
		Values(pNetBlock.Len, pNetBlock.Cap, backendPeerIDArrStr.String()).
		Exec()
	if err != nil {
		return err
	}

	return nil
}
