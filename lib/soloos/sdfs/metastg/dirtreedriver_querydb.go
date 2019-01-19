package metastg

import (
	"database/sql"
	"soloos/sdfs/types"
	"time"

	"github.com/gocraft/dbr"
)

func (p *DirTreeDriver) DeleteFsINodeByIDInDB(fsINodeID types.FsINodeID) error {
	var (
		sess *dbr.Session
		err  error
	)

	sess = p.helper.DBConn.NewSession(nil)
	_, err = sess.DeleteFrom("b_fsinode").
		Where("fsinode_ino=?", fsINodeID).
		Exec()
	return err
}

func (p *DirTreeDriver) ListFsINodeByParentIDFromDB(parentID types.FsINodeID,
	beforeLiteralFunc func(resultCount int) bool,
	literalFunc func(types.FsINode) bool,
) error {
	var (
		sess          *dbr.Session
		sqlRows       *sql.Rows
		ret           types.FsINode
		netINodeIDStr string
		resultCount   int
		err           error
	)

	sess = p.helper.DBConn.NewSession(nil)

	sqlRows, err = sess.Select("count(fsinode_ino) as result").
		From("b_fsinode").
		Where("parent_fsinode_ino=?", parentID).Rows()
	if err != nil {
		goto QUERY_DONE
	}
	if sqlRows.Next() {
		err = sqlRows.Scan(&resultCount)
		if err != nil {
			goto QUERY_DONE
		}
	}
	sqlRows.Close()

	if beforeLiteralFunc(resultCount) == false {
		goto QUERY_DONE
	}

	sqlRows, err = sess.Select(schemaDirTreeFsINodeAttr...).
		From("b_fsinode").
		Where("parent_fsinode_ino=?", parentID).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	for sqlRows.Next() {
		err = sqlRows.Scan(
			&ret.Ino,
			&netINodeIDStr,
			&ret.ParentID,
			&ret.Name,
			&ret.Type,
			&ret.Atime,
			&ret.Ctime,
			&ret.Mtime,
			&ret.Atimensec,
			&ret.Ctimensec,
			&ret.Mtimensec,
			&ret.Mode,
			&ret.Nlink,
			&ret.Uid,
			&ret.Gid,
			&ret.Rdev,
		)
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

func (p *DirTreeDriver) UpdateFsINodeInDB(fsINode *types.FsINode) error {
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

	fsINode.Mtime = types.DirTreeTime(time.Now().Unix())
	_, err = sess.Update("b_fsinode").
		Set("fsinode_ino", fsINode.Ino).
		Set("netinode_id", string(fsINode.NetINodeID[:])).
		Set("parent_fsinode_ino", fsINode.ParentID).
		Set("fsinode_name", fsINode.Name).
		Set("fsinode_type", fsINode.Type).
		Set("atime", fsINode.Atime).
		Set("ctime", fsINode.Ctime).
		Set("mtime", fsINode.Mtime).
		Set("atimensec", fsINode.Atimensec).
		Set("ctimensec", fsINode.Ctimensec).
		Set("mtimensec", fsINode.Mtimensec).
		Set("mode", fsINode.Mode).
		Set("nlink", fsINode.Nlink).
		Set("uid", fsINode.Uid).
		Set("gid", fsINode.Gid).
		Set("rdev", fsINode.Rdev).
		Where("fsinode_ino=?", fsINode.Ino).
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
		Columns(schemaDirTreeFsINodeAttr...).
		Values(
			fsINode.Ino,
			string(fsINode.NetINodeID[:]),
			fsINode.ParentID,
			fsINode.Name,
			fsINode.Type,
			fsINode.Atime,
			fsINode.Ctime,
			fsINode.Mtime,
			fsINode.Atimensec,
			fsINode.Ctimensec,
			fsINode.Mtimensec,
			fsINode.Mode,
			fsINode.Nlink,
			fsINode.Uid,
			fsINode.Gid,
			fsINode.Rdev,
		).
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

func (p *DirTreeDriver) GetFsINodeByIDFromDB(fsINodeID types.FsINodeID) (types.FsINode, error) {
	var (
		sess          *dbr.Session
		sqlRows       *sql.Rows
		ret           types.FsINode
		netINodeIDStr string
		exists        bool
		err           error
	)

	p.fsINodesByIDRWMutex.RLock()
	ret, exists = p.fsINodesByID[fsINodeID]
	p.fsINodesByIDRWMutex.RUnlock()
	if exists {
		return ret, nil
	}

	sess = p.helper.DBConn.NewSession(nil)
	sqlRows, err = sess.Select(schemaDirTreeFsINodeAttr...).
		From("b_fsinode").
		Where("fsinode_ino=?",
			fsINodeID,
		).Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = types.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(
		&ret.Ino,
		&netINodeIDStr,
		&ret.ParentID,
		&ret.Name,
		&ret.Type,
		&ret.Atime,
		&ret.Ctime,
		&ret.Mtime,
		&ret.Atimensec,
		&ret.Ctimensec,
		&ret.Mtimensec,
		&ret.Mode,
		&ret.Nlink,
		&ret.Uid,
		&ret.Gid,
		&ret.Rdev,
	)
	if err != nil {
		goto QUERY_DONE
	}
	copy(ret.NetINodeID[:], []byte(netINodeIDStr))

QUERY_DONE:
	if err == nil {
		err = p.prepareAndSetFsINodeCache(&ret)
	}
	if sqlRows != nil {
		sqlRows.Close()
	}
	return ret, err
}

func (p *DirTreeDriver) GetFsINodeByNameFromDB(parentID types.FsINodeID, fsINodeName string) (types.FsINode, error) {
	var (
		sess          *dbr.Session
		sqlRows       *sql.Rows
		ret           types.FsINode
		netINodeIDStr string
		exists        bool
		err           error
	)

	p.fsINodesByPathRWMutex.RLock()
	ret, exists = p.fsINodesByPath[p.MakeFsINodeKey(parentID, fsINodeName)]
	p.fsINodesByPathRWMutex.RUnlock()
	if exists {
		return ret, nil
	}

	sess = p.helper.DBConn.NewSession(nil)
	sqlRows, err = sess.Select(schemaDirTreeFsINodeAttr...).
		From("b_fsinode").
		Where("parent_fsinode_ino=? and fsinode_name=?",
			parentID, fsINodeName,
		).Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = types.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(
		&ret.Ino,
		&netINodeIDStr,
		&ret.ParentID,
		&ret.Name,
		&ret.Type,
		&ret.Atime,
		&ret.Ctime,
		&ret.Mtime,
		&ret.Atimensec,
		&ret.Ctimensec,
		&ret.Mtimensec,
		&ret.Mode,
		&ret.Nlink,
		&ret.Uid,
		&ret.Gid,
		&ret.Rdev,
	)
	if err != nil {
		goto QUERY_DONE
	}
	copy(ret.NetINodeID[:], []byte(netINodeIDStr))

QUERY_DONE:
	if err == nil {
		err = p.prepareAndSetFsINodeCache(&ret)
	}
	if sqlRows != nil {
		sqlRows.Close()
	}
	return ret, err
}
