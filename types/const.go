package types

import (
	"soloos/common/fsapitypes"
	"soloos/common/sdfsapitypes"
	"syscall"
)

const (
	MaxInt    = 1<<31 - 1
	MaxUint64 = ^uint64(0)
	MinUint64 = 0
	MaxInt64  = int(MaxUint64 >> 1)
	MinInt64  = -MaxInt - 1

	TRUE  = 1
	FALSE = 0
)

const (
	FSINODE_TYPE_FILE = iota
	FSINODE_TYPE_DIR
	FSINODE_TYPE_HARD_LINK
	FSINODE_TYPE_SOFT_LINK
	FSINODE_TYPE_FIFO
	FSINODE_TYPE_UNKOWN
)

const (
	FS_MAX_PATH_LENGTH = sdfsapitypes.FS_MAX_PATH_LENGTH
	FS_MAX_NAME_LENGTH = sdfsapitypes.FS_MAX_NAME_LENGTH
	FS_RDEV            = 0

	FS_EEXIST       = fsapitypes.Status(syscall.EEXIST)
	FS_ENOTEMPTY    = fsapitypes.Status(syscall.ENOTEMPTY)
	FS_ENAMETOOLONG = fsapitypes.Status(syscall.ENAMETOOLONG)

	FS_INODE_LOCK_SH = syscall.LOCK_SH
	FS_INODE_LOCK_EX = syscall.LOCK_EX
	FS_INODE_LOCK_UN = syscall.LOCK_UN

	FS_XATTR_SOFT_LNKMETA_KEY = "sdfs.soft.link.metadata"
)

const (
	FS_PERM_SETUID uint32 = 1 << (12 - 1 - iota)
	FS_PERM_SETGID
	FS_PERM_STICKY
	FS_PERM_USER_READ
	FS_PERM_USER_WRITE
	FS_PERM_USER_EXECUTE
	FS_PERM_GROUP_READ
	FS_PERM_GROUP_WRITE
	FS_PERM_GROUP_EXECUTE
	FS_PERM_OTHER_READ
	FS_PERM_OTHER_WRITE
	FS_PERM_OTHER_EXECUTE
)
