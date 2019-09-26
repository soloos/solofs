package metastg

import "soloos/common/log"

// func prepareDirTreeDataSqls() []string {
// var sql []string
// sql = append(sql, fmt.Sprintf(`
// insert into b_fsinode (fsinode_ino,parent_fsinode_ino,fsinode_name,netinode_id,fsinode_type) values(%d,%d,"","",%d);
// `, solofsapitypes.RootFsINodeID, solofsapitypes.RootFsINodeParentID, solofstypes.FSINODE_TYPE_DIR))
// return sql
// }

func (p *PosixFS) installSchema() error {
	var (
		sqls []string
		err  error
	)

	sqls = p.prepareDirTreeSqls()
	for _, sql := range sqls {
		_, err = p.dbConn.Exec(sql)
		if err != nil {
			log.Error(err, sql)
		}
	}

	return nil
}

func (p *PosixFS) prepareDirTreeSqls() []string {
	var sql []string

	// sql = append(sql, `
	// drop table b_maxid;
	// `)
	sql = append(sql, `
	create table if not exists b_maxid (
	mkey char(64),
	maxid int,
	primary key(mkey)
	);
`)

	// sql = append(sql, `
	// drop table b_fsinode;
	// `)
	sql = append(sql, `
	create table if not exists b_fsinode (
	namespace_id int,
	fsinode_ino bigint,
	hardlink_ino bigint,
	netinode_id char(64),
	parent_fsinode_ino bigint,
	fsinode_name char(128),
	fsinode_type int,
	atime bigint default 0,
	ctime bigint default 0,
	mtime bigint default 0,
	atimensec bigint default 0,
	ctimensec bigint default 0,
	mtimensec bigint default 0,
	mode int default 0,
	nlink int default 0,
	uid int default 0,
	gid int default 0,
	rdev int default 0,
	primary key(namespace_id, fsinode_ino)
	);
`)

	sql = append(sql, `
	create unique index if not exists i_b_fsinode_parent_fsinode_ino_and_fsinode_name 
	on b_fsinode(namespace_id, parent_fsinode_ino, fsinode_name);
`)

	sql = append(sql, `
	create table if not exists b_fsinode_xattr (
	namespace_id int,
	fsinode_ino bigint,
	xattr text,
	primary key(namespace_id, fsinode_ino)
	);
`)

	// sql = append(sql, prepareDirTreeDataSqls()...)

	return sql
}
