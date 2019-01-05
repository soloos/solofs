package types

const (
	MaxInt = 1<<31 - 1

	TRUE  = 1
	FALSE = 0

	DefaultNetBlockCap int = 1024 * 1024 * 8
	DefaultMemBlockCap int = 1024 * 1024 * 2
)

const (
	FSINODE_TYPE_FILE = iota
	FSINODE_TYPE_DIR
)
