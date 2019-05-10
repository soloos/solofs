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
	literalFunc func(types.FsINodeMeta) bool,
) error {
	var (
		sess            sdbapi.Session
		sqlRows         *sql.Rows
		ret             types.FsINodeMeta
		fetchRowsLimit  uint64
		fetchRowsOffset uint64
		netINodeIDStr   string
		resultCount     int
		fsINodeName     string
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
				&fsINodeName,
				&ret.Type,
				&ret.Mode,
			)
		} else {
			err = sqlRows.Scan(
				&ret.Ino,
				&ret.HardLinkIno,
				&netINodeIDStr,
				&ret.ParentID,
				&fsINodeName,
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
		ret.SetName(fsINodeName)

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

func (p *FsINodeDriver) UpdateFsINodeInDB(pFsINodeMeta *types.FsINodeMeta) error {
	var (
		sess sdbapi.Session
		err  error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	_, err = sess.Update("b_fsinode").
		Set("fsinode_ino", pFsINodeMeta.Ino).
		Set("hardlink_ino", pFsINodeMeta.HardLinkIno).
		Set("netinode_id", string(pFsINodeMeta.NetINodeID[:])).
		Set("parent_fsinode_ino", pFsINodeMeta.ParentID).
		Set("fsinode_name", pFsINodeMeta.Name()).
		Set("fsinode_type", pFsINodeMeta.Type).
		Set("atime", pFsINodeMeta.Atime).
		Set("ctime", pFsINodeMeta.Ctime).
		Set("mtime", pFsINodeMeta.Mtime).
		Set("atimensec", pFsINodeMeta.Atimensec).
		Set("ctimensec", pFsINodeMeta.Ctimensec).
		Set("mtimensec", pFsINodeMeta.Mtimensec).
		Set("mode", pFsINodeMeta.Mode).
		Set("nlink", pFsINodeMeta.Nlink).
		Set("uid", pFsINodeMeta.Uid).
		Set("gid", pFsINodeMeta.Gid).
		Set("rdev", pFsINodeMeta.Rdev).
		Where("fsinode_ino=?", pFsINodeMeta.Ino).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) InsertFsINodeInDB(pFsINodeMeta *types.FsINodeMeta) error {
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
			pFsINodeMeta.Ino,
			pFsINodeMeta.HardLinkIno,
			string(pFsINodeMeta.NetINodeID[:]),
			pFsINodeMeta.ParentID,
			pFsINodeMeta.Name(),
			pFsINodeMeta.Type,
			pFsINodeMeta.Atime,
			pFsINodeMeta.Ctime,
			pFsINodeMeta.Mtime,
			pFsINodeMeta.Atimensec,
			pFsINodeMeta.Ctimensec,
			pFsINodeMeta.Mtimensec,
			pFsINodeMeta.Mode,
			pFsINodeMeta.Nlink,
			pFsINodeMeta.Uid,
			pFsINodeMeta.Gid,
			pFsINodeMeta.Rdev,
		).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) FetchFsINodeByIDFromDB(pFsINodeMeta *types.FsINodeMeta) error {
	var (
		fsINodeName   string
		sess          sdbapi.Session
		sqlRows       *sql.Rows
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
			pFsINodeMeta.Ino,
		).Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = types.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(
		&pFsINodeMeta.Ino,
		&pFsINodeMeta.HardLinkIno,
		&netINodeIDStr,
		&pFsINodeMeta.ParentID,
		&fsINodeName,
		&pFsINodeMeta.Type,
		&pFsINodeMeta.Atime,
		&pFsINodeMeta.Ctime,
		&pFsINodeMeta.Mtime,
		&pFsINodeMeta.Atimensec,
		&pFsINodeMeta.Ctimensec,
		&pFsINodeMeta.Mtimensec,
		&pFsINodeMeta.Mode,
		&pFsINodeMeta.Nlink,
		&pFsINodeMeta.Uid,
		&pFsINodeMeta.Gid,
		&pFsINodeMeta.Rdev,
	)
	pFsINodeMeta.SetName(fsINodeName)

	if err != nil {
		goto QUERY_DONE
	}
	copy(pFsINodeMeta.NetINodeID[:], []byte(netINodeIDStr))

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return err
}

func (p *FsINodeDriver) FetchFsINodeByNameFromDB(pFsINodeMeta *types.FsINodeMeta) error {
	var (
		fsINodeName   string
		sess          sdbapi.Session
		sqlRows       *sql.Rows
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
			pFsINodeMeta.ParentID, pFsINodeMeta.Name(),
		).Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = types.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(
		&pFsINodeMeta.Ino,
		&pFsINodeMeta.HardLinkIno,
		&netINodeIDStr,
		&pFsINodeMeta.ParentID,
		&fsINodeName,
		&pFsINodeMeta.Type,
		&pFsINodeMeta.Atime,
		&pFsINodeMeta.Ctime,
		&pFsINodeMeta.Mtime,
		&pFsINodeMeta.Atimensec,
		&pFsINodeMeta.Ctimensec,
		&pFsINodeMeta.Mtimensec,
		&pFsINodeMeta.Mode,
		&pFsINodeMeta.Nlink,
		&pFsINodeMeta.Uid,
		&pFsINodeMeta.Gid,
		&pFsINodeMeta.Rdev,
	)
	pFsINodeMeta.SetName(fsINodeName)

	if err != nil {
		goto QUERY_DONE
	}
	copy(pFsINodeMeta.NetINodeID[:], []byte(netINodeIDStr))

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return err
}
