package metastg

func fsSchemaSqls() []string {
	var sql []string
	// sql = append(sql, `
	// drop table b_fsinode;
	// `)
	sql = append(sql, `
	create table if not exists b_fsinode (
	fsinode_id int,
	parent_fsinode_id int,
	fsinode_name char(64),
	netinode_id char(64),
	primary key(netinode_id)
	);
`)

	sql = append(sql, `
	create index if not exists i_b_fsindoe_parent_fsinode_id on b_fsinode(parent_fsinode_id);
`)

	return sql
}

func commonSchemaSqls() []string {
	var sql []string

	sql = append(fsSchemaSqls())

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
	backend_peer_id_arr text,
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
