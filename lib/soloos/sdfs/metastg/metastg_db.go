package metastg

func commonSchemaSql() string {
	var sql = `
--drop table b_inode;
create table if not exists b_inode (
        inodeid char(64),
        inodesize int,
        netblocksize int,
        memblocksize int,
        primary key(inodeid)
);

--drop table b_netblock;
create table if not exists b_netblock (
        netblockid char(64),
        inodeid char(64),
        index_in_inode int,
        netblocksize int,
        primary key(netblockid)
);
create index if not exists b_inode_netblock_inodeid_index_in_inode on b_netblock(inodeid, index_in_inode);

--drop table r_netblock_store_peer;
create table if not exists r_netblock_store_peer (
        netblockid char(64),
        peerid char(64),
        primary key(netblockid,peerid)
);
`
	return sql
}
