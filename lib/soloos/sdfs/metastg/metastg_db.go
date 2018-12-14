package metastg

func commonSchemaSqls() []string {
	var sql []string
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
