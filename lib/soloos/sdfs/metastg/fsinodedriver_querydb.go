package metastg

import (
	"database/sql"
	"soloos/common/sdbapi"
	"soloos/common/sdfsapitypes"
)

func (p *FsINodeDriver) DeleteFsINodeByIDInDB(fsINodeID sdfsapitypes.FsINodeID) error {
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

func (p *FsINodeDriver) ListFsINodeByParentIDFromDB(parentID sdfsapitypes.FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(sdfsapitypes.FsINodeMeta) bool,
) error {
	var (
		sess            sdbapi.Session
		sqlRows         *sql.Rows
		ret             sdfsapitypes.FsINodeMeta
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

func (p *FsINodeDriver) UpdateFsINodeInDB(fsINodeMeta sdfsapitypes.FsINodeMeta) error {
	var (
		sess sdbapi.Session
		err  error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	_, err = sess.Update("b_fsinode").
		Set("fsinode_ino", fsINodeMeta.Ino).
		Set("hardlink_ino", fsINodeMeta.HardLinkIno).
		Set("netinode_id", string(fsINodeMeta.NetINodeID[:])).
		Set("parent_fsinode_ino", fsINodeMeta.ParentID).
		Set("fsinode_name", fsINodeMeta.Name()).
		Set("fsinode_type", fsINodeMeta.Type).
		Set("atime", fsINodeMeta.Atime).
		Set("ctime", fsINodeMeta.Ctime).
		Set("mtime", fsINodeMeta.Mtime).
		Set("atimensec", fsINodeMeta.Atimensec).
		Set("ctimensec", fsINodeMeta.Ctimensec).
		Set("mtimensec", fsINodeMeta.Mtimensec).
		Set("mode", fsINodeMeta.Mode).
		Set("nlink", fsINodeMeta.Nlink).
		Set("uid", fsINodeMeta.Uid).
		Set("gid", fsINodeMeta.Gid).
		Set("rdev", fsINodeMeta.Rdev).
		Where("fsinode_ino=?", fsINodeMeta.Ino).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) InsertFsINodeInDB(fsINodeMeta sdfsapitypes.FsINodeMeta) error {
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
			fsINodeMeta.Ino,
			fsINodeMeta.HardLinkIno,
			string(fsINodeMeta.NetINodeID[:]),
			fsINodeMeta.ParentID,
			fsINodeMeta.Name(),
			fsINodeMeta.Type,
			fsINodeMeta.Atime,
			fsINodeMeta.Ctime,
			fsINodeMeta.Mtime,
			fsINodeMeta.Atimensec,
			fsINodeMeta.Ctimensec,
			fsINodeMeta.Mtimensec,
			fsINodeMeta.Mode,
			fsINodeMeta.Nlink,
			fsINodeMeta.Uid,
			fsINodeMeta.Gid,
			fsINodeMeta.Rdev,
		).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) FetchFsINodeByIDFromDB(fsINodeID sdfsapitypes.FsINodeID) (sdfsapitypes.FsINodeMeta, error) {
	var (
		fsINodeMeta   sdfsapitypes.FsINodeMeta
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
		Where("fsinode_ino=?", fsINodeID).
		Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = sdfsapitypes.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(
		&fsINodeMeta.Ino,
		&fsINodeMeta.HardLinkIno,
		&netINodeIDStr,
		&fsINodeMeta.ParentID,
		&fsINodeName,
		&fsINodeMeta.Type,
		&fsINodeMeta.Atime,
		&fsINodeMeta.Ctime,
		&fsINodeMeta.Mtime,
		&fsINodeMeta.Atimensec,
		&fsINodeMeta.Ctimensec,
		&fsINodeMeta.Mtimensec,
		&fsINodeMeta.Mode,
		&fsINodeMeta.Nlink,
		&fsINodeMeta.Uid,
		&fsINodeMeta.Gid,
		&fsINodeMeta.Rdev,
	)
	fsINodeMeta.SetName(fsINodeName)

	if err != nil {
		goto QUERY_DONE
	}
	copy(fsINodeMeta.NetINodeID[:], []byte(netINodeIDStr))

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return fsINodeMeta, err
}

func (p *FsINodeDriver) FetchFsINodeByNameFromDB(parentID sdfsapitypes.FsINodeID,
	fsINodeName string) (sdfsapitypes.FsINodeMeta, error) {
	var (
		fsINodeMeta   sdfsapitypes.FsINodeMeta
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
			parentID, fsINodeName,
		).Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = sdfsapitypes.ErrObjectNotExists
		goto QUERY_DONE
	}

	{
		var name string
		err = sqlRows.Scan(
			&fsINodeMeta.Ino,
			&fsINodeMeta.HardLinkIno,
			&netINodeIDStr,
			&fsINodeMeta.ParentID,
			&name,
			&fsINodeMeta.Type,
			&fsINodeMeta.Atime,
			&fsINodeMeta.Ctime,
			&fsINodeMeta.Mtime,
			&fsINodeMeta.Atimensec,
			&fsINodeMeta.Ctimensec,
			&fsINodeMeta.Mtimensec,
			&fsINodeMeta.Mode,
			&fsINodeMeta.Nlink,
			&fsINodeMeta.Uid,
			&fsINodeMeta.Gid,
			&fsINodeMeta.Rdev,
		)
		fsINodeMeta.SetName(name)
	}

	if err != nil {
		goto QUERY_DONE
	}
	copy(fsINodeMeta.NetINodeID[:], []byte(netINodeIDStr))

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return fsINodeMeta, err
}
