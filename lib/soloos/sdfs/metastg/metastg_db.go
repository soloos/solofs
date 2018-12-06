package metastg

func commonSchemaSqls() []string {
	var sql []string
	// sql = append(sql, `
	// drop table b_inode;
	// `)
	sql = append(sql, `
	create table if not exists b_inode (
	inode_id char(64),
	inode_size int,
	netblock_cap int,
	memblock_cap int,
	primary key(inode_id)
	);
`)

	// sql = append(sql, `
	// drop table b_netblock;
	// `)
	sql = append(sql, `
	create table if not exists b_netblock (
	netblock_id char(64),
	inode_id char(64),
	index_in_inode int,
	netblock_len int,
	netblock_cap int,
	primary key(netblock_id)
	);
`)
	sql = append(sql, `
	create index if not exists b_inode_netblock_inode_id_index_in_inode on b_netblock(inode_id, index_in_inode);
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
