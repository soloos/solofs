package metastg

func (p *MetaStg) prepareNetINodesSqls() []string {
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
	netinode_size int,
	netblock_cap int,
	memblock_cap int,
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
	create unique index if not exists i_b_fsinode_netinode_id
	on b_fsinode(netinode_id);
`)

	sql = append(sql, `
	create table if not exists b_fsinode_xattr (
	namespace_id int,
	fsinode_ino bigint,
	xattr text,
	primary key(namespace_id, fsinode_ino)
	);
`)
	// sql = append(sql, `
	// drop table b_netinode;
	// `)
	sql = append(sql, `
	create table if not exists b_netinode (
	netinode_id char(64),
	netinode_size int,
	netblock_cap int,
	memblock_cap int,
	primary key(netinode_id)
	);
`)

	// sql = append(sql, `
	// drop table b_netblock;
	// `)
	sql = append(sql, `
	create table if not exists b_netblock (
	netinode_id char(64),
	index_in_netinode int,
	netblock_len int,
	netblock_cap int,
	backend_peer_id_arr varchar(1024),
	primary key(netinode_id, index_in_netinode)
	);
`)

	// sql = append(sql, `
	// drop table r_netblock_store_peer;
	// `)
	sql = append(sql, `
	create table if not exists r_netblock_store_peer (
	netblock_id char(64),
	peer_id char(64),
	primary key(netblock_id,peer_id)
	);
`)

	return sql
}
