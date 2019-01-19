package metastg

// func prepareDirTreeDataSqls() []string {
// var sql []string
// sql = append(sql, fmt.Sprintf(`
// insert into b_fsinode (fsinode_ino,parent_fsinode_ino,fsinode_name,netinode_id,fsinode_type) values(%d,%d,"","",%d);
// `, types.RootFsINodeID, types.RootFsINodeParentID, types.FSINODE_TYPE_IFDIR))
// return sql
// }

func prepareDirTreeSqls() []string {
	var sql []string

	// sql = append(sql, `
	// drop table b_maxid;
	// `)
	sql = append(sql, `
	create table if not exists b_maxid (
	key char(64),
	maxid int,
	primary key(key)
	);
`)

	// sql = append(sql, `
	// drop table b_fsinode;
	// `)
	sql = append(sql, `
	create table if not exists b_fsinode (
	fsinode_ino big int,
	netinode_id char(64),
	parent_fsinode_ino big int,
	fsinode_name char(64),
	fsinode_type int,
	atime big int default 0,
	ctime big int default 0,
	mtime big int default 0,
	atimensec big int default 0,
	ctimensec big int default 0,
	mtimensec big int default 0,
	mode int default 0,
	nlink int default 0,
	uid int default 0,
	gid int default 0,
	rdev int default 0,
	primary key(fsinode_ino)
	);
`)

	sql = append(sql, `
	create unique index if not exists i_b_fsinode_parent_fsinode_ino_and_fsinode_name 
	on b_fsinode(parent_fsinode_ino, fsinode_name);
`)

	// sql = append(sql, prepareDirTreeDataSqls()...)

	return sql
}
