package metastg

import (
	"database/sql"
	"soloos/sdfs/types"

	"github.com/gocraft/dbr"
)

func (p *DirTreeDriver) DeleteFsINodeByIDInDB(fsINodeID types.FsINodeID) error {
	var (
		sess *dbr.Session
		err  error
	)

	sess = p.helper.DBConn.NewSession(nil)
	_, err = sess.DeleteFrom("b_fsinode").
		Where("fsinode_id=?", fsINodeID).
		Exec()
	return err
}

func (p *DirTreeDriver) ListFsINodeByParentIDFromDB(parentID types.FsINodeID, literalFunc func(types.FsINode) bool) error {
	var (
		sess          *dbr.Session
		sqlRows       *sql.Rows
		ret           types.FsINode
		netINodeIDStr string
		err           error
	)

	sess = p.helper.DBConn.NewSession(nil)
	sqlRows, err = sess.Select("fsinode_id", "parent_fsinode_id", "name", "flag", "permission", "netinode_id", "fsinode_type").
		From("b_fsinode").
		Where("parent_fsinode_id=?",
			parentID,
		).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	for sqlRows.Next() {
		err = sqlRows.Scan(&ret.ID, &ret.ParentID, &ret.Name, &ret.Flag, &ret.Permission, &netINodeIDStr, &ret.Type)
		if err != nil {
			goto QUERY_DONE
		}
		copy(ret.NetINodeID[:], []byte(netINodeIDStr))
		if literalFunc(ret) == false {
			break
		}
	}

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return err
}

func (p *DirTreeDriver) UpdateFsINodeInDB(fsINode types.FsINode) error {
	var (
		sess *dbr.Session
		tx   *dbr.Tx
		err  error
	)

	sess = p.helper.DBConn.NewSession(nil)
	tx, err = sess.Begin()
	if err != nil {
		goto QUERY_DONE
	}

	_, err = sess.Update("b_fsinode").
		Set("fsinode_id", fsINode.ID).
		Set("parent_fsinode_id", fsINode.ParentID).
		Set("name", fsINode.Name).
		Set("flag", fsINode.Flag).
		Set("permission", fsINode.Permission).
		Set("netinode_id", string(fsINode.NetINodeID[:])).
		Set("fsinode_type", fsINode.Type).
		Where("fsinode_id=?", fsINode.ID).
		Limit(1).
		Exec()
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

func (p *DirTreeDriver) InsertFsINodeInDB(fsINode types.FsINode) error {
	var (
		sess *dbr.Session
		tx   *dbr.Tx
		err  error
	)

	sess = p.helper.DBConn.NewSession(nil)
	tx, err = sess.Begin()
	if err != nil {
		goto QUERY_DONE
	}

	_, err = sess.InsertInto("b_fsinode").
		Columns("fsinode_id", "parent_fsinode_id", "name", "flag", "permission", "netinode_id", "fsinode_type").
		Values(fsINode.ID, fsINode.ParentID, fsINode.Name, fsINode.Flag, fsINode.Permission, string(fsINode.NetINodeID[:]), fsINode.Type).
		Exec()
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

func (p *DirTreeDriver) GetFsINodeByIDFromDB(parentID types.FsINodeID, fsINodeName string) (types.FsINode, error) {
	var (
		sess          *dbr.Session
		sqlRows       *sql.Rows
		key           = p.MakeFsINodeKey(parentID, fsINodeName)
		ret           types.FsINode
		netINodeIDStr string
		exists        bool
		err           error
	)

	p.fsINodesRWMutex.RLock()
	ret, exists = p.fsINodes[key]
	p.fsINodesRWMutex.RUnlock()
	if exists {
		return ret, nil
	}

	p.fsINodesRWMutex.Lock()
	sess = p.helper.DBConn.NewSession(nil)
	sqlRows, err = sess.Select("fsinode_id", "parent_fsinode_id", "name", "flag", "permission", "netinode_id", "fsinode_type").
		From("b_fsinode").
		Where("parent_fsinode_id=? and name=?",
			parentID, fsINodeName,
		).Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = types.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(&ret.ID, &ret.ParentID, &ret.Name, &ret.Flag, &ret.Permission, &netINodeIDStr, &ret.Type)
	if err != nil {
		goto QUERY_DONE
	}
	copy(ret.NetINodeID[:], []byte(netINodeIDStr))

QUERY_DONE:
	p.fsINodesRWMutex.Unlock()
	if sqlRows != nil {
		sqlRows.Close()
	}
	return ret, err
}
