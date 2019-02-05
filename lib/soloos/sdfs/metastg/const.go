package metastg

const (
	// TestMetaStgDBDriver  = "mysql"
	// TestMetaStgDBConnect = "root:hello@/sdfs_test"
	TestMetaStgDBDriver  = "sqlite3"
	TestMetaStgDBConnect = "/tmp/sdfs_test.db"
)

var (
	schemaDirTreeFsINodeBasicAttr = []string{
		"fsinode_ino",
		"hardlink_ino",
		"netinode_id",
		"parent_fsinode_ino",
		"fsinode_name",
		"fsinode_type",
		"mode",
	}

	schemaDirTreeFsINodeAttr = []string{
		"fsinode_ino",
		"hardlink_ino",
		"netinode_id",
		"parent_fsinode_ino",
		"fsinode_name",
		"fsinode_type",
		"atime",
		"ctime",
		"mtime",
		"atimensec",
		"ctimensec",
		"mtimensec",
		"mode",
		"nlink",
		"uid",
		"gid",
		"rdev",
	}
)
