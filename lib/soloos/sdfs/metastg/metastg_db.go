package metastg

func commonSchemaSqls() []string {
	var sql []string
	// sql = append(sql, `
	// drop table b_netnetINode;
	// `)
	sql = append(sql, `
	create table if not exists b_netnetINode (
	netnetINode_id char(64),
	netnetINode_size int,
	netblock_cap int,
	memblock_cap int,
	primary key(netnetINode_id)
	);
`)

	// sql = append(sql, `
	// drop table b_netblock;
	// `)
	sql = append(sql, `
	create table if not exists b_netblock (
	netblock_id char(64),
	netnetINode_id char(64),
	index_in_netnetINode int,
	netblock_len int,
	netblock_cap int,
	backend_peer_id_arr text,
	primary key(netblock_id)
	);
`)
	sql = append(sql, `
	create index if not exists b_netnetINode_netblock_netnetINode_id_index_in_netnetINode on b_netblock(netnetINode_id, index_in_netnetINode);
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
