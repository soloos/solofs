package metastg

import (
	"database/sql"
	"soloos/common/solodbapi"
	"soloos/common/solofstypes"
)

func (p *NetINodeDriver) NetINodeCommitSizeInDB(uNetINode solofstypes.NetINodeUintptr, size uint64) error {
	var (
		sess solodbapi.Session
		err  error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	var query = sess.Update("b_netinode").
		Set("netinode_size", size).
		Where("netinode_id=?", uNetINode.Ptr().IDStr())
	if sess.Dialect != "sqlite3" {
		query = query.Limit(1)
	}

	_, err = query.Exec()
	if err != nil {
		return err
	}

	uNetINode.Ptr().Size = size

	return nil
}

func (p *NetINodeDriver) FetchNetINodeFromDB(pNetINode *solofstypes.NetINode) error {
	var (
		sess    solodbapi.Session
		sqlRows *sql.Rows
		err     error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		goto QUERY_DONE
	}

	sqlRows, err = sess.Select("netinode_size", "netblock_cap", "memblock_cap").
		From("b_netinode").
		Where("netinode_id=?", pNetINode.IDStr()).Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = solofstypes.ErrObjectNotExists
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

func (p *NetINodeDriver) StoreNetINodeInDB(pNetINode *solofstypes.NetINode) error {
	var (
		sess          solodbapi.Session
		netINodeIDStr = pNetINode.IDStr()
		err           error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	err = sess.ReplaceInto("b_netinode").
		PrimaryColumns("netinode_id").PrimaryValues(netINodeIDStr).
		Columns("netinode_size", "netblock_cap", "memblock_cap").
		Values(pNetINode.Size, pNetINode.NetBlockCap, pNetINode.MemBlockCap).
		Exec()
	if err != nil {
		return err
	}

	return nil
}
