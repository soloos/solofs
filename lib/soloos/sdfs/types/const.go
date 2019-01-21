package types

import (
	"syscall"

	"github.com/hanwen/go-fuse/fuse"
)

const (
	MaxInt    = 1<<31 - 1
	MaxUint64 = ^uint64(0)
	MinUint64 = 0
	MaxInt64  = int(MaxUint64 >> 1)
	MinInt64  = -MaxInt - 1

	TRUE  = 1
	FALSE = 0

	DefaultNetBlockCap int = 1024 * 1024 * 8
	DefaultMemBlockCap int = 1024 * 1024 * 2
)

const (
	FSINODE_TYPE_FILE = iota
	FSINODE_TYPE_DIR
	FSINODE_TYPE_HARD_LINK
	FSINODE_TYPE_SOFT_LINK
	FSINODE_TYPE_FIFO
)

const (
	FS_EEXIST    = fuse.Status(syscall.EEXIST)
	FS_ENOTEMPTY = fuse.Status(syscall.ENOTEMPTY)

	FS_INODE_LOCK_SH = syscall.LOCK_SH
	FS_INODE_LOCK_EX = syscall.LOCK_EX
	FS_INODE_LOCK_UN = syscall.LOCK_UN

	FS_XATTR_SOFT_LNKMETA_KEY = "sdfs.soft.link.metadata"
)
