package types

import (
	"unsafe"
)

type FsINodeID = uint64
type DirTreeTime = uint64
type DirTreeTimeNsec = uint32

const (
	RootFsINodeParentID = FsINodeID(0)
	RootFsINodeID       = FsINodeID(1)
)

type FsINodeUintptr uintptr

func (u FsINodeUintptr) Ptr() *FsINode { return (*FsINode)(unsafe.Pointer(u)) }

type FsINode struct {
	Ino        FsINodeID
	NetINodeID NetINodeID
	ParentID   FsINodeID
	Name       string
	Type       int
	Atime      DirTreeTime
	Ctime      DirTreeTime
	Mtime      DirTreeTime
	Atimensec  DirTreeTimeNsec
	Ctimensec  DirTreeTimeNsec
	Mtimensec  DirTreeTimeNsec
	Mode       uint32
	Nlink      uint32
	Uid        uint32
	Gid        uint32
	Rdev       uint32
	UNetINode  NetINodeUintptr
}

type FsINodeFileHandler struct {
	FsINodeID      FsINodeID
	AppendPosition int64
	ReadPosition   int64
}

func (p *FsINodeFileHandler) Reset() {
	p.AppendPosition = 0
	p.ReadPosition = 0
}
