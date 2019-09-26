package metastg

import (
	"database/sql"
	"soloos/common/solodbapi"
	"soloos/common/solofsapitypes"
)

func (p *FsINodeDriver) DeleteFsINodeByIDInDB(nameSpaceID solofsapitypes.NameSpaceID,
	fsINodeID solofsapitypes.FsINodeID) error {
	var (
		sess solodbapi.Session
		err  error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	_, err = sess.DeleteFrom("b_fsinode").
		Where("namespace_id=? and fsinode_ino=?", nameSpaceID, fsINodeID).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) ListFsINodeByParentIDFromDB(nameSpaceID solofsapitypes.NameSpaceID,
	parentID solofsapitypes.FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(solofsapitypes.FsINodeMeta) bool,
) error {
	var (
		sess            solodbapi.Session
		sqlRows         *sql.Rows
		ret             solofsapitypes.FsINodeMeta
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
		Where("namespace_id=? and parent_fsinode_ino=?", nameSpaceID, parentID).Rows()
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
			Where("namespace_id=? and parent_fsinode_ino=?", nameSpaceID, parentID).
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
				&ret.NameSpaceID,
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
				&ret.NameSpaceID,
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

func (p *FsINodeDriver) UpdateFsINodeInDB(nameSpaceID solofsapitypes.NameSpaceID,
	fsINodeMeta solofsapitypes.FsINodeMeta) error {
	var (
		sess solodbapi.Session
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
		Where("namespace_id=? and fsinode_ino=?", nameSpaceID, fsINodeMeta.Ino).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) InsertFsINodeInDB(nameSpaceID solofsapitypes.NameSpaceID,
	fsINodeMeta solofsapitypes.FsINodeMeta) error {
	var (
		sess solodbapi.Session
		err  error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	fsINodeMeta.NameSpaceID = nameSpaceID

	_, err = sess.InsertInto("b_fsinode").
		Columns(schemaDirTreeFsINodeAttr...).
		Values(
			fsINodeMeta.NameSpaceID,
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

func (p *FsINodeDriver) FetchFsINodeByIDFromDB(nameSpaceID solofsapitypes.NameSpaceID,
	fsINodeID solofsapitypes.FsINodeID) (solofsapitypes.FsINodeMeta, error) {
	var (
		fsINodeMeta   solofsapitypes.FsINodeMeta
		fsINodeName   string
		sess          solodbapi.Session
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
		Where("namespace_id=? and fsinode_ino=?", nameSpaceID, fsINodeID).
		Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = solofsapitypes.ErrObjectNotExists
		goto QUERY_DONE
	}

	err = sqlRows.Scan(
		&fsINodeMeta.NameSpaceID,
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

func (p *FsINodeDriver) FetchFsINodeByNameFromDB(nameSpaceID solofsapitypes.NameSpaceID,
	parentID solofsapitypes.FsINodeID,
	fsINodeName string) (solofsapitypes.FsINodeMeta, error) {
	var (
		fsINodeMeta   solofsapitypes.FsINodeMeta
		sess          solodbapi.Session
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
		Where("namespace_id=? and parent_fsinode_ino=? and fsinode_name=?",
			nameSpaceID, parentID, fsINodeName,
		).Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = solofsapitypes.ErrObjectNotExists
		goto QUERY_DONE
	}

	{
		var name string
		err = sqlRows.Scan(
			&fsINodeMeta.NameSpaceID,
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
