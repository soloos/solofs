package metastg

import (
	"fmt"
	"soloos/sdfs/types"
)

func baseDataSqls() []string {
	var sql []string
	sql = append(sql, fmt.Sprintf(`
	insert into b_fsinode (fsinode_id,parent_fsinode_id,name,netinode_id,fsinode_type) values(%d,%d,"","",%d);
	`, 0, -1, types.FSINODE_TYPE_DIR))
	return sql
}

func fsSchemaSqls() []string {
	var sql []string
	// sql = append(sql, `
	// drop table b_fsinode;
	// `)
	sql = append(sql, `
	create table if not exists b_fsinode (
	fsinode_id int,
	parent_fsinode_id int,
	name char(64),
	flag int default 0,
	permission int default 0,
	netinode_id char(64),
	fsinode_type int,
	primary key(fsinode_id)
	);
`)

	sql = append(sql, `
	create index if not exists i_b_fsinode_parent_fsinode_id_and_name on b_fsinode(parent_fsinode_id, name);
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
