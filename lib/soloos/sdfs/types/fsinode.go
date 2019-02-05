package types

import (
	"soloos/util/offheap"
	"unsafe"

	"github.com/hanwen/go-fuse/fuse"
)

type FsINodeID = uint64
type DirTreeTime = uint64
type DirTreeTimeNsec = uint32

const (
	ZombieFsINodeParentID = FsINodeID(0)
	RootFsINodeParentID   = FsINodeID(0)
	RootFsINodeID         = FsINodeID(1)
	FsINodeStructSize     = unsafe.Sizeof(FsINode{})
	MaxFsINodeID          = MaxUint64
)

type FsINodeUintptr uintptr

func (u FsINodeUintptr) Ptr() *FsINode { return (*FsINode)(unsafe.Pointer(u)) }

type FsINode struct {
	SharedPointer     offheap.SharedPointer
	LastModifyACMTime int64
	LoadInMemAt       int64

	Ino         FsINodeID
	HardLinkIno FsINodeID
	NetINodeID  NetINodeID
	ParentID    FsINodeID
	Name        string
	Type        int
	Atime       DirTreeTime
	Ctime       DirTreeTime
	Mtime       DirTreeTime
	Atimensec   DirTreeTimeNsec
	Ctimensec   DirTreeTimeNsec
	Mtimensec   DirTreeTimeNsec
	Mode        uint32
	Nlink       int32
	Uid         uint32
	Gid         uint32
	Rdev        uint32
	UNetINode   NetINodeUintptr
}

func (p *FsINode) Reset() {
	p.SharedPointer.Reset()
}

func FsModeToFsINodeType(mode uint32) int {
	if mode&fuse.S_IFDIR != 0 {
		return FSINODE_TYPE_DIR
	}
	if mode&fuse.S_IFREG != 0 {
		return FSINODE_TYPE_FILE
	}
	if mode&fuse.S_IFIFO != 0 {
		return FSINODE_TYPE_FIFO
	}
	if mode&fuse.S_IFLNK != 0 {
		return FSINODE_TYPE_SOFT_LINK
	}
	return FSINODE_TYPE_UNKOWN
}

func FsINodeTypeToFuseType(fsINodeType int) int {
	switch fsINodeType {
	case FSINODE_TYPE_DIR:
		return fuse.S_IFDIR
	case FSINODE_TYPE_FILE:
		return fuse.S_IFREG
	case FSINODE_TYPE_FIFO:
		return fuse.S_IFIFO
	case FSINODE_TYPE_SOFT_LINK:
		return fuse.S_IFLNK
	}
	return fuse.S_IFREG
}
