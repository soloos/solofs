package metastg

import (
	"database/sql"
	"soloos/common/solodbapi"
	"soloos/common/solofstypes"
)

func (p *FsINodeDriver) DeleteFsINodeByIDInDB(
	nsID solofstypes.NameSpaceID,
	fsINodeID solofstypes.FsINodeID) error {
	var (
		sess solodbapi.Session
		err  error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	_, err = sess.DeleteFrom("b_fsinode").
		Where("namespace_id=? and fsinode_ino=?", nsID, fsINodeID).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) UpdateFsINodeInDB(
	nsID solofstypes.NameSpaceID,
	fsINodeMeta solofstypes.FsINodeMeta) error {
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
		Where("namespace_id=? and fsinode_ino=?", nsID, fsINodeMeta.Ino).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) InsertFsINodeInDB(
	nsID solofstypes.NameSpaceID,
	fsINodeMeta solofstypes.FsINodeMeta) error {
	var (
		sess solodbapi.Session
		err  error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		return err
	}

	fsINodeMeta.NameSpaceID = nsID

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

func (p *FsINodeDriver) FetchFsINodeByIDFromDB(
	nsID solofstypes.NameSpaceID,
	fsINodeID solofstypes.FsINodeID) (solofstypes.FsINodeMeta, error) {
	var (
		fsINodeMeta   solofstypes.FsINodeMeta
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
		Where("namespace_id=? and fsinode_ino=?", nsID, fsINodeID).
		Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = solofstypes.ErrObjectNotExists
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

func (p *FsINodeDriver) FetchFsINodeByNameFromDB(
	nsID solofstypes.NameSpaceID,
	parentID solofstypes.FsINodeID,
	fsINodeName string) (solofstypes.FsINodeMeta, error) {
	var (
		fsINodeMeta   solofstypes.FsINodeMeta
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
			nsID, parentID, fsINodeName,
		).Limit(1).Rows()
	if err != nil {
		goto QUERY_DONE
	}

	if sqlRows.Next() == false {
		err = solofstypes.ErrObjectNotExists
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

func (p *FsINodeDriver) ListFsINodeByParentIDSelectCountFromDB(
	nsID solofstypes.NameSpaceID,
	parentID solofstypes.FsINodeID,
) (int64, error) {
	var (
		sess        solodbapi.Session
		sqlRows     *sql.Rows
		resultCount int64
		err         error
	)

	err = p.dbConn.InitSession(&sess)
	if err != nil {
		goto QUERY_DONE
	}

	sqlRows, err = sess.Select("count(fsinode_ino) as result").
		From("b_fsinode").
		Where("namespace_id=? and parent_fsinode_ino=?", nsID, parentID).Rows()
	if err != nil {
		goto QUERY_DONE
	}
	if sqlRows.Next() {
		err = sqlRows.Scan(&resultCount)
		if err != nil {
			goto QUERY_DONE
		}
	}

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return resultCount, err
}

func (p *FsINodeDriver) ListFsINodeByParentIDSelectDataFromDB(
	nsID solofstypes.NameSpaceID,
	parentID solofstypes.FsINodeID,
	fetchRowsLimit uint64,
	fetchRowsOffset uint64,
	isFetchAllCols bool,
) ([]solofstypes.FsINodeMeta, error) {
	var (
		sess          solodbapi.Session
		sqlRows       *sql.Rows
		ret           []solofstypes.FsINodeMeta
		retRow        solofstypes.FsINodeMeta
		netINodeIDStr string
		fsINodeName   string
		err           error
	)

	if fetchRowsLimit == 0 {
		goto QUERY_DONE
	}

	err = p.dbConn.InitSession(&sess)
	if err != nil {
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
			Where("namespace_id=? and parent_fsinode_ino=?", nsID, parentID).
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
				&retRow.NameSpaceID,
				&retRow.Ino,
				&retRow.HardLinkIno,
				&netINodeIDStr,
				&retRow.ParentID,
				&fsINodeName,
				&retRow.Type,
				&retRow.Mode,
			)
		} else {
			err = sqlRows.Scan(
				&retRow.NameSpaceID,
				&retRow.Ino,
				&retRow.HardLinkIno,
				&netINodeIDStr,
				&retRow.ParentID,
				&fsINodeName,
				&retRow.Type,
				&retRow.Atime,
				&retRow.Ctime,
				&retRow.Mtime,
				&retRow.Atimensec,
				&retRow.Ctimensec,
				&retRow.Mtimensec,
				&retRow.Mode,
				&retRow.Nlink,
				&retRow.Uid,
				&retRow.Gid,
				&retRow.Rdev,
			)
		}
		retRow.SetName(fsINodeName)

		if err != nil {
			goto QUERY_DONE
		}
		copy(retRow.NetINodeID[:], []byte(netINodeIDStr))
		ret = append(ret, retRow)
	}

QUERY_DONE:
	if sqlRows != nil {
		sqlRows.Close()
	}
	return ret, err
}
