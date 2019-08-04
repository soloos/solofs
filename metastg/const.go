package metastg

const (
	// TestMetaStgDBDriver  = "mysql"
	// TestMetaStgDBConnect = "root:hello@/soloos_test"
	TestMetaStgDBDriver  = "sqlite3"
	TestMetaStgDBConnect = "/tmp/soloos_test.db"
)

var (
	schemaDirTreeFsINodeBasicAttr = []string{
		"namespace_id",
		"fsinode_ino",
		"hardlink_ino",
		"netinode_id",
		"parent_fsinode_ino",
		"fsinode_name",
		"fsinode_type",
		"mode",
	}

	schemaDirTreeFsINodeAttr = []string{
		"namespace_id",
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
