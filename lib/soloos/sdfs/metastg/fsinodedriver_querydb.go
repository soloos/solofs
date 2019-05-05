package metastg

import (
	"database/sql"
	"soloos/common/sdbapi"
	"soloos/sdfs/types"
)

func (p *FsINodeDriver) DeleteFsINodeByIDInDB(fsINodeID types.FsINodeID) error {
	var (
		sess sdbapi.Session
		err  error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	_, err = sess.DeleteFrom("b_fsinode").
		Where("fsinode_ino=?", fsINodeID).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) ListFsINodeByParentIDFromDB(parentID types.FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(types.FsINode) bool,
) error {
	var (
		sess            sdbapi.Session
		sqlRows         *sql.Rows
		ret             types.FsINode
		fetchRowsLimit  uint64
		fetchRowsOffset uint64
		netINodeIDStr   string
		resultCount     int
		err             error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		goto QUERY_DONE
	}

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

	fetchRowsLimit, fetchRowsOffset = beforeLiteralFunc(resultCount)
	if fetchRowsLimit == 0 {
		goto QUERY_DONE
	}

	{
		var schemaAttr []string
		if isFetchAllCols == false {
			schemaAttr = schemaDirTreeFsINodeBasicAttr

		} else {
			schemaAttr = schemaDirTreeFsINodeAttr
		}
		sqlRows, err = sess.Select(schemaAttr...).
			From("b_fsinode").
			Where("parent_fsinode_ino=?", parentID).
			OrderDesc("fsinode_ino").
			Limit(fetchRowsLimit).
			Offset(fetchRowsOffset).
			Rows()
	}

	if err != nil {
		goto QUERY_DONE
	}

	for sqlRows.Next() {
		if isFetchAllCols == false {
			err = sqlRows.Scan(
				&ret.Ino,
				&ret.HardLinkIno,
				&netINodeIDStr,
				&ret.ParentID,
				&ret.Name,
				&ret.Type,
				&ret.Mode,
			)
		} else {
			err = sqlRows.Scan(
				&ret.Ino,
				&ret.HardLinkIno,
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
		}
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

func (p *FsINodeDriver) UpdateFsINodeInDB(fsINode types.FsINode) error {
	var (
		sess sdbapi.Session
		err  error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	_, err = sess.Update("b_fsinode").
		Set("fsinode_ino", fsINode.Ino).
		Set("hardlink_ino", fsINode.HardLinkIno).
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
		return err
	}

	return nil
}

func (p *FsINodeDriver) InsertFsINodeInDB(fsINode types.FsINode) error {
	var (
		sess sdbapi.Session
		err  error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	_, err = sess.InsertInto("b_fsinode").
		Columns(schemaDirTreeFsINodeAttr...).
		Values(
			fsINode.Ino,
			fsINode.HardLinkIno,
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
		return err
	}

	return nil
}

func (p *FsINodeDriver) GetFsINodeByIDFromDB(fsINodeID types.FsINodeID) (types.FsINode, error) {
	var (
		sess          sdbapi.Session
		sqlRows       *sql.Rows
		ret           types.FsINode
		netINodeIDStr string
		err           error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		goto QUERY_DONE
	}

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
		&ret.HardLinkIno,
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
	if sqlRows != nil {
		sqlRows.Close()
	}
	return ret, err
}

func (p *FsINodeDriver) GetFsINodeByNameFromDB(parentID types.FsINodeID, fsINodeName string) (types.FsINode, error) {
	var (
		sess          sdbapi.Session
		sqlRows       *sql.Rows
		ret           types.FsINode
		netINodeIDStr string
		err           error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		goto QUERY_DONE
	}

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
		&ret.HardLinkIno,
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
	if sqlRows != nil {
		sqlRows.Close()
	}
	return ret, err
}
