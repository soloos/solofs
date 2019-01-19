package types

import "syscall"

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
	FSINODE_TYPE_IFREG = syscall.S_IFREG
	FSINODE_TYPE_IFDIR = syscall.S_IFDIR
	FSINODE_TYPE_IFLNK = syscall.S_IFLNK
	FSINODE_TYPE_IFIFO = syscall.S_IFIFO
)
