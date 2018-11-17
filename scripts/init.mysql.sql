CREATE DATABASE `sdfs_test` /*!40100 DEFAULT CHARACTER SET utf8 */;
use sdfs_test;

drop table b_inode;
create table b_inode (
        inodeid char(64),
        inodesize int,
        netblocksize int,
        memblocksize int,
        primary key(inodeid)
);

drop table b_netblock;
create table b_netblock (
        netblockid char(64),
        inodeid char(64),
        index_in_inode int,
        netblocksize int,
        primary key(netblockid),
        key(inodeid, index_in_inode)
);

drop table r_netblock_store_peer;
create table r_netblock_store_peer (
        netblockid char(64),
        peerid char(64),
        primary key(netblockid,peerid)
);
